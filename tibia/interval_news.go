package tibia

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/common"
	"github.com/mediocregopher/radix/v3"
)

type NewsTable struct {
	common.SmallModel

	RunScan bool
	LastID  int
}

func (nt *NewsTable) TableName() string {
	return "news_table"
}

type InnerNewsStruct struct {
	GuildID      int64 `gorm:"primary_key"`
	IsChannelSet bool
	ChannelID    int64
	DMSent       bool
	RunNews      bool
}

func (ins *InnerNewsStruct) TableName() string {
	return "inner_news_struct"
}

var (
	newsMasterwg sync.WaitGroup // WaitGroup to be used on the newsController
	newswg       sync.WaitGroup // WaitGroup to be used while ranging over guild channels and sending the new
	newsLocker   sync.WaitGroup // WaitGroup to be used to prevent race conditions while managing the redis table "news_guilds"
	newspool     chan struct{}  // Channel to be used as a queue for the range loop to send new messages
	outputChan   chan int64     // Channel to update guilds slice
)

func newsController() {
	table := NewsTable{}
	done := make(chan bool, 1)
	go func() {
		for {
			err := common.GORM.Where(&NewsTable{}).First(&table).Error
			alreadySet := err != gorm.ErrRecordNotFound
			if err != nil && alreadySet {
				logger.Errorf("Err on newsController: %v", err)
				return
			}

			if !alreadySet || !table.RunScan {
				done <- true
			}

			ct := time.Now()

			select {
			case <-done:
				close(done)
				logger.Info("News feed done")
				return
			default:
				newsMasterwg.Add(1)
				logger.Info("Running news tracker")
				scanNews()
			}

			newsMasterwg.Wait()
			logger.Infof("News tracker finished in %v", time.Since(ct).Seconds())

			time.Sleep(30 * time.Second)
		}
	}()
}

func scanNews() {
	defer newsMasterwg.Done()

	news, err := GetTibiaNews()
	if err != nil || news == nil {
		logger.Errorf("Error getting tibia's latest new to send embeds. Error: %g", err)
		return
	}

	var lastID int
	table := NewsTable{}
	err = common.GORM.Where(&table).First(&table).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			lastID = news.ID
			table.LastID = news.ID
			err = common.GORM.Save(&table).Error
			if err != nil {
				logger.Error("Error setting last news id.")
			}
		} else {
			logger.Error("Error getting last news id!")
			return
		}
	} else {
		lastID = table.LastID
	}

	if lastID == news.ID {
		return // No news to send
	}

	table.LastID = news.ID

	embed := &discordgo.MessageEmbed{
		Title:       news.Title,
		Description: fmt.Sprintf("%s\n[Clique para ver mais](%s)", news.ShortDescription, news.URL),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %d\nData: %s", news.ID, news.Date),
		},
	}

	err = NewsLoop(embed)
	if err != nil {
		logger.Errorf("Error on the News Loop: %v", err)
	}

	err = common.GORM.Save(&table).Error
	if err != nil {
		logger.Errorf("Error updating last news. Error: %v", err)
	}
}

func NewsLoop(news *discordgo.MessageEmbed) error {
	newsLocker.Add(1)
	defer newsLocker.Done()

	var table []int64
	err := common.RedisPool.Do(radix.Cmd(&table, "SMEMBERS", "news_guilds"))
	if err != nil {
		return errors.WithMessage(err, "Error fetching table data for news loop.")
	}

	newspool = make(chan struct{}, 100)
	outputChan = make(chan int64, len(table))
	rand.Seed(time.Now().UnixNano())
	for _, g := range table {
		newswg.Add(1)
		go newsSender(g, news)
	}

	newswg.Wait()
	close(newspool)
	close(outputChan)

	goBack := make([]string, len(outputChan)+1)
	var i int
	for v := range outputChan {
		goBack[i+1] = fmt.Sprint(v)
		i++
	}

	goBack[0] = "news_guilds"

	if len(goBack) > 1 {
		err = common.RedisPool.Do(radix.Cmd(nil, "SADD", goBack...))
		if err != nil {
			return errors.WithMessage(err, "Error setting table data for news loop. (trying to add)")
		}
	} else {
		logger.Info("else triggered")
		err = common.RedisPool.Do(radix.Cmd(nil, "DEL", "news_guilds"))
		if err != nil {
			return errors.WithMessage(err, "Error setting table data for news loop. (trying to delete)")
		}
	}

	return nil
}

func newsSender(g int64, news *discordgo.MessageEmbed) {
	defer func() {
		if r := recover(); r != nil {
			logger.Infof("Recovered at: %v", r)
		}

		defer newswg.Done()
	}()

	newspool <- struct{}{}
	defer func() { <-newspool }()

	gs := bot.State.Guild(true, g)
	if gs == nil {
		logger.Errorf("Really weird bug with guild %d. It is not in state.", g)
		go func(g int64) {
			table := InnerNewsStruct{}
			err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
			if err == nil {
				common.GORM.Delete(&table)
			}
		}(g)
		return
	}

	gs.RLock()
	defer gs.RUnlock()

	table := InnerNewsStruct{}
	err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}

		outputChan <- g
		logger.Errorf("Error fething DB to send news on guild: %d -- Error: %v", g, err)
		return
	}

	if table.IsChannelSet && table.RunNews {
		news.Color = int(rand.Int63n(16777215))
		perms, _, err := bot.SendMessageEmbed(g, table.ChannelID, news)
		if !perms {
			if !table.DMSent {
				err = bot.SendDM(gs.Guild.OwnerID, fmt.Sprintf("Looks like I don't have perms to send msgs on the news channel of your server. Please give me permission to Send Message on the channel <#%d> or change the news feed to another channel.", table.ChannelID))
				if err != nil {
					logger.Errorf("Failed sending DM to the owner of Guild: %d -- User: %d", g, gs.Guild.OwnerID)
				}

				table.DMSent = true
			}
		} else {
			if err != nil {
				logger.Errorf("Error sending news DM on guild %d", g)
			} else if table.DMSent {
				table.DMSent = false
				logger.Infof("Guild %d news channel is now working!", g)
			}
		}

		err = common.GORM.Save(&table).Error
		if err != nil {
			logger.Errorf("Failed updating guild %d news table", g)
		}
	}

	outputChan <- g
}

func StartNewsLoop() (string, error) {
	table := NewsTable{}
	err := common.GORM.Where(&table).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if alreadySet && table.RunScan {
		return "O tracking de news já está rolando.", nil
	}

	table.RunScan = true

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	newsController()

	return "Tudo certo! O tracking de news está rolando", nil
}

func StopNewsLoop() (string, error) {
	table := NewsTable{}
	err := common.GORM.Where(&table).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if !alreadySet || !table.RunScan {
		return "O tracking de news não está rolando ainda.", nil
	}

	table.RunScan = false

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	return "Tudo certo! O tracking de news foi pausado.", nil
}

func CreateNewsFeed(g, c int64) (string, error) {
	table := InnerNewsStruct{}
	err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "ERROR 1:")
	}

	if alreadySet && table.RunNews {
		return "O news feed ja está rolando.", nil
	}

	if table.IsChannelSet {
		return "O feed ja foi criado, use o comando \"cnfc\" para mudar o canal.", nil
	}

	channel, err := common.BotSession.Channel(c)
	if err != nil {
		return fmt.Sprintf("O canal especificado não foi encontrado.\nError: %v", err), err
	} else if channel == nil {
		return "O canal especificado não foi encontrado", nil
	}

	newsLocker.Wait()
	newsLocker.Add(1)
	var added int
	err = common.RedisPool.Do(radix.FlatCmd(&added, "SADD", "news_guilds", g))
	newsLocker.Done()
	if err != nil {
		return "", errors.WithMessage(err, "ERROR 2:")
	}

	if added == 0 {
		logger.Errorf("Weird bug with guild %d news feed - setting it up", g)
	}

	table.GuildID = g
	table.IsChannelSet = true
	table.RunNews = true
	table.ChannelID = c

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", errors.WithMessage(err, "ERROR 3:")
	}

	return fmt.Sprintf("Tudo certo! O feed de notícias será feito no canal <#%d>", c), nil
}

func EnableNewsFeed(g int64) (string, error) {
	table := InnerNewsStruct{}
	err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if alreadySet && table.RunNews {
		return "O news feed ja está rolando.", nil
	}

	if !table.IsChannelSet {
		return "Você, primeiro, precisa usar o comando \"CreateNewsFeed\" para configurar o feed.", nil
	}

	newsLocker.Wait()
	newsLocker.Add(1)
	var added int
	err = common.RedisPool.Do(radix.FlatCmd(&added, "SADD", "news_guilds", g))
	newsLocker.Done()
	if err != nil {
		return "", err
	}

	if added == 0 {
		logger.Errorf("Weird bug with guild %d news feed - enabling it up", g)
	}

	table.RunNews = true

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	return "Tudo certo! O feed de notícias está ativado!", nil
}

func DisableNewsFeed(g int64) (string, error) {
	table := InnerNewsStruct{}
	err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if !alreadySet || !table.RunNews {
		return "O news feed já está desabilitado neste servidor.", nil
	}

	newsLocker.Wait()
	newsLocker.Add(1)
	var beingTracked bool
	err = common.RedisPool.Do(radix.FlatCmd(&beingTracked, "SREM", "news_guilds", g))
	newsLocker.Done()
	if err != nil {
		return "", err
	}

	if !beingTracked {
		logger.Errorf("Weird bug with guild %d news feed", g)
	}

	table.RunNews = false

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	return "Tudo certo! O feed de notícias foi desativado!", nil
}

func ChangeNewsFeedChannel(g, c int64) (string, error) {
	table := InnerNewsStruct{}
	err := common.GORM.Where(&InnerNewsStruct{GuildID: g}).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if !alreadySet {
		return "O news feed está desabilitado neste servidor.", nil
	}

	if table.ChannelID == c {
		return "Esse já é o canal usado para o news feed.", nil
	}

	channel, err := common.BotSession.Channel(c)
	if err != nil {
		return fmt.Sprintf("O canal especificado não foi encontrado.\nError: %v", err), err
	} else if channel == nil {
		return "O canal especificado não foi encontrado.", nil
	}

	table.ChannelID = c
	table.IsChannelSet = true

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Tudo certo! O canal do feed de notícias foi alterado para <#%d>!", c), nil
}

func DebugNews(c int64) (string, error) {
	var table []int64
	err := common.RedisPool.Do(radix.Cmd(&table, "SMEMBERS", "news_guilds"))
	if err != nil {
		return "", errors.WithMessage(err, "Error: 1")
	}

	newsTable := NewsTable{}
	err = common.GORM.Where(&newsTable).First(&newsTable).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "Error: 2")
	}

	out := "NewsTable:\n"

	if !alreadySet {
		out += "`NewsTable not found`\n\n\n"
	} else {
		out += fmt.Sprintf("`%#v`\n\n\n", newsTable)
	}

	out += fmt.Sprintf("**Redis**:\nLen: %d\n", len(table))

	for i, v := range table {
		out += fmt.Sprintf("**Index**: `%d` -- **Guild**: `%d`\n", i, v)
	}

	if len(out) <= 2000 {
		return out, nil
	}

	msg := &discordgo.MessageSend{}

	var buf bytes.Buffer
	buf.WriteString(out)

	msg.File = &discordgo.File{
		Name:        "Attachment.txt",
		ContentType: "text/plain",
		Reader:      &buf,
	}

	_, err = common.BotSession.ChannelMessageSendComplex(c, msg)
	if err != nil {
		return "", errors.WithMessage(err, "Error: 3")
	}

	return "", nil
}
