package premium

//go:generate sqlboiler psql

import (
	"context"
	"database/sql"
	"encoding/base32"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/premium/models"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var (
	ErrCodeExpired  = errors.New("Code expired")
	ErrCodeNotFound = errors.New("Code not found")
)

func init() {
	RegisterPremiumSource(&CodePremiumSource{})
}

type CodePremiumSource struct{}

func (ps *CodePremiumSource) Init() {}

func (ps *CodePremiumSource) Names() (human string, idname string) {
	return "Redeemed code", "code"
}

func RedeemCode(ctx context.Context, code string, userID int64) error {
	tx, err := common.PQ.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithMessage(err, "BeginTX")
	}

	// Query for the code model
	c, err := models.PremiumCodes(qm.Where("code = ? AND user_id IS NULL", code), qm.For("UPDATE")).One(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return errors.WithMessage(err, "models.PremiumCodes")
	}

	// model found, with no user attached, create the slot for it
	slot, err := CreatePremiumSlot(ctx, tx, userID, "code", "Redeemed code", c.Message, c.ID, time.Duration(c.Duration), PremiumTierPremium)
	if err != nil {
		_ = tx.Rollback()
		return errors.WithMessage(err, "CreatePremiumSlot")
	}

	// Update the code fields
	c.UserID = null.Int64From(userID)
	c.UsedAt = null.TimeFrom(time.Now())
	c.SlotID = null.Int64From(slot.ID)

	_, err = c.Update(ctx, tx, boil.Infer())
	if err != nil {
		_ = tx.Rollback()
		return errors.WithMessage(err, "Update")
	}

	err = tx.Commit()
	return errors.WithMessage(err, "Commit")
}

func LookupCode(ctx context.Context, code string) (*models.PremiumCode, error) {
	c, err := models.PremiumCodes(qm.Where("code = ?", code)).OneG(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCodeNotFound
		}

		return nil, errors.WithMessage(err, "models.PremiumCodes")
	}

	return c, nil
}

func LookUpSlot(ctx context.Context, sourceID int64) (*models.PremiumSlot, error) {
	p, err := models.PremiumSlots(qm.Where("source_id = ?", sourceID)).OneG(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCodeNotFound
		}

		return nil, errors.WithMessage(err, "models.PremiumSlots")
	}

	return p, nil
}

var (
	ErrCodeCollision = errors.New("Code collision")
)

// TryRetryGenerateCode attempts to generate codes, if it enocunters a key collision it retries, returns on all other cases
func TryRetryGenerateCode(ctx context.Context, message string, duration time.Duration) (*models.PremiumCode, error) {
	for {
		code, err := GenerateCode(ctx, message, duration)
		if err != nil && err == ErrCodeCollision {
			logger.WithError(err).Error("Code collision!")
			continue
		}

		return code, err
	}
}

// GenerateCode generates a redeemable premium code with the specified duration (-1 for permanent) and message
func GenerateCode(ctx context.Context, message string, duration time.Duration) (*models.PremiumCode, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.WithMessage(err, "GenerateCode")
	}

	encoded := encodeKey(key)

	model := &models.PremiumCode{
		Code:      encoded,
		Message:   message,
		Permanent: duration == -1,
		Duration:  int64(duration),
	}

	err = model.InsertG(ctx, boil.Infer())
	if err != nil {
		if cast, ok := errors.Cause(err).(*pq.Error); ok {
			if cast.Code == "23505" {
				return nil, ErrCodeCollision
			}
		}
	}

	return model, err
}

var keyEncoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func encodeKey(rawKey []byte) string {
	str := keyEncoder.EncodeToString(rawKey)
	var output strings.Builder
	for i, r := range str {
		if i%6 == 0 && i != 0 {
			output.WriteString("-")
		}
		output.WriteString(string(r))
	}

	return output.String()
}

var cmdGenerateCode = &commands.TIDCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "generatepremiumcode",
	Aliases:              []string{"gpc"},
	Description:          "Generates premium codes",
	HideFromHelp:         true,
	RequiredArgs:         3,
	RunInDM:              true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Duration", Type: &commands.DurationArg{}},
		{Name: "NumCodes", Type: dcmd.Int},
		{Name: "Message", Type: dcmd.String},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		numKeys := data.Args[1].Int()
		duration := data.Args[0].Value.(time.Duration)
		codes := make([]string, 0, numKeys)

		if duration <= 0 {
			duration = -1
		}

		for i := 0; i < numKeys; i++ {
			code, err := TryRetryGenerateCode(data.Context(), data.Args[2].Str(), duration)
			if err != nil {
				return nil, err
			}

			codes = append(codes, code.Code)
		}

		var dm strings.Builder

		dm.WriteString(fmt.Sprintf("Duration: `%s`, Permanent: `%t`, Message: `%s`\n```\n", duration.String(), duration == -1, data.Args[2].Str()))

		for _, v := range codes {
			dm.WriteString(v + "\n")
		}

		dm.WriteString("```")

		if _, err := bot.SendDM(data.Msg.Author.ID, dm.String()); err != nil {
			return fmt.Sprintf("I wasn't able to send you a DM.\nError:%v", err), err
		}

		return "Check your dms", nil
	}),
}

var cmdDeleteCode = &commands.TIDCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "deletepremiumcode",
	Aliases:              []string{"dpc"},
	Description:          "Deletes a premium codes",
	HideFromHelp:         true,
	RequiredArgs:         1,
	RunInDM:              true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Code", Type: dcmd.String},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		code := data.Args[0].Str()
		amount, err := DeletePremiumCode(data.Context(), code)
		if err != nil {
			if err == ErrCodeNotFound {
				return "Código não encontrado", nil
			}

			return fmt.Sprintf("Deleted %d with the following error:\nRaw error: `%#v`\nString error: `%s`", amount, err, err), err
		}

		return "Code " + code + " deleted.\nTotal amount: " + discordgo.StrID(amount), nil
	}),
}

func DeletePremiumCode(ctx context.Context, code string) (amount int64, err error) {
	dbCode, err := LookupCode(ctx, code)
	if err != nil {
		return 0, err
	}

	amount, err = dbCode.DeleteG(ctx)
	if err != nil {
		return amount, err
	}

	slot, err := LookUpSlot(ctx, dbCode.ID)
	if err != nil && err != ErrCodeNotFound {
		return amount, errors.WithMessage(err, "error 1")
	}

	if err == nil {
		err = DetachSlotFromGuild(ctx, slot.ID, slot.UserID)
		if err != nil {
			return amount, errors.WithMessage(err, "error detaching from guild")
		}

		_, err = slot.DeleteG(ctx)
		if err != nil {
			return amount, errors.WithMessage(err, "error 2")
		}
	}

	return amount, err
}

var cmdPremiumInfo = &commands.TIDCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "premium",
	Aliases:     []string{"prem", "premmy"},
	Description: "Gives you information about tidbot's premium",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		title := "Tid Bot Premium"

		host := common.ConfHost.GetString()

		thumb := &discordgo.MessageEmbedThumbnail{
			URL: common.BotUser.AvatarURL("1024"),
		}

		cor := int(rand.Int63n(16777215))

		descr := fmt.Sprintf("Obrigado por querer ser premium! :)\nSer [premium](https://%s/premium) no tidbot tem vários benefícios. Tais como:\n\n• Suporte Prioritário\n• Aumento do armazenamento de mensagens de 1 para 12 horas. O que significa que você pode buscar mensagens apagadas no servidor nas últimas 12 horas\n• Limite de custom commands aumentado de 100 para 250.\n• Aumento de 10 vezes no limite de databases\n• **Permissão de usar as funções de simultaniedade que são até 10 vezes mais rápidas que as funções padrões de tibia.**\n\n• Aumento de limite de várias funções\n→ **getChar** 40 calls por CC;\n→ **getDeaths** 40 calls por CC; \n→ **getDeath** 40 calls por CC;\n→ **getGuild** 5 calls por CC;\n→ **getGuildMembers** 2 calls por CC;\n→ **checkOnline** 5 calls por CC;\n→ **getNews** 3 calls por CC;\n→ libera as funções **getMultipleChars** e **getMultipleCharsDeath**.\n\nPara mais informações, [clique aqui](https://%s/premium).", host, host)

		embed := &discordgo.MessageEmbed{
			Title:       title,
			Thumbnail:   thumb,
			Color:       cor,
			Description: descr,
		}

		return embed, nil
	},
}
