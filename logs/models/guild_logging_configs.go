// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// GuildLoggingConfig is an object representing the database table.
type GuildLoggingConfig struct {
	GuildID                      int64            `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	CreatedAt                    null.Time        `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt                    null.Time        `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	UsernameLoggingEnabled       null.Bool        `boil:"username_logging_enabled" json:"username_logging_enabled,omitempty" toml:"username_logging_enabled" yaml:"username_logging_enabled,omitempty"`
	NicknameLoggingEnabled       null.Bool        `boil:"nickname_logging_enabled" json:"nickname_logging_enabled,omitempty" toml:"nickname_logging_enabled" yaml:"nickname_logging_enabled,omitempty"`
	BlacklistedChannels          null.String      `boil:"blacklisted_channels" json:"blacklisted_channels,omitempty" toml:"blacklisted_channels" yaml:"blacklisted_channels,omitempty"`
	ManageMessagesCanViewDeleted null.Bool        `boil:"manage_messages_can_view_deleted" json:"manage_messages_can_view_deleted,omitempty" toml:"manage_messages_can_view_deleted" yaml:"manage_messages_can_view_deleted,omitempty"`
	EveryoneCanViewDeleted       null.Bool        `boil:"everyone_can_view_deleted" json:"everyone_can_view_deleted,omitempty" toml:"everyone_can_view_deleted" yaml:"everyone_can_view_deleted,omitempty"`
	MessageLogsAllowedRoles      types.Int64Array `boil:"message_logs_allowed_roles" json:"message_logs_allowed_roles,omitempty" toml:"message_logs_allowed_roles" yaml:"message_logs_allowed_roles,omitempty"`
	AccessMode                   int16            `boil:"access_mode" json:"access_mode" toml:"access_mode" yaml:"access_mode"`

	R *guildLoggingConfigR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L guildLoggingConfigL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GuildLoggingConfigColumns = struct {
	GuildID                      string
	CreatedAt                    string
	UpdatedAt                    string
	UsernameLoggingEnabled       string
	NicknameLoggingEnabled       string
	BlacklistedChannels          string
	ManageMessagesCanViewDeleted string
	EveryoneCanViewDeleted       string
	MessageLogsAllowedRoles      string
	AccessMode                   string
}{
	GuildID:                      "guild_id",
	CreatedAt:                    "created_at",
	UpdatedAt:                    "updated_at",
	UsernameLoggingEnabled:       "username_logging_enabled",
	NicknameLoggingEnabled:       "nickname_logging_enabled",
	BlacklistedChannels:          "blacklisted_channels",
	ManageMessagesCanViewDeleted: "manage_messages_can_view_deleted",
	EveryoneCanViewDeleted:       "everyone_can_view_deleted",
	MessageLogsAllowedRoles:      "message_logs_allowed_roles",
	AccessMode:                   "access_mode",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpernull_Bool struct{ field string }

func (w whereHelpernull_Bool) EQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bool) NEQ(x null.Bool) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bool) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bool) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Bool) LT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bool) LTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bool) GT(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bool) GTE(x null.Bool) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpertypes_Int64Array struct{ field string }

func (w whereHelpertypes_Int64Array) EQ(x types.Int64Array) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_Int64Array) NEQ(x types.Int64Array) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_Int64Array) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_Int64Array) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpertypes_Int64Array) LT(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_Int64Array) LTE(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_Int64Array) GT(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_Int64Array) GTE(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperint16 struct{ field string }

func (w whereHelperint16) EQ(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint16) NEQ(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint16) LT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint16) LTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint16) GT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint16) GTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint16) IN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint16) NIN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var GuildLoggingConfigWhere = struct {
	GuildID                      whereHelperint64
	CreatedAt                    whereHelpernull_Time
	UpdatedAt                    whereHelpernull_Time
	UsernameLoggingEnabled       whereHelpernull_Bool
	NicknameLoggingEnabled       whereHelpernull_Bool
	BlacklistedChannels          whereHelpernull_String
	ManageMessagesCanViewDeleted whereHelpernull_Bool
	EveryoneCanViewDeleted       whereHelpernull_Bool
	MessageLogsAllowedRoles      whereHelpertypes_Int64Array
	AccessMode                   whereHelperint16
}{
	GuildID:                      whereHelperint64{field: "\"guild_logging_configs\".\"guild_id\""},
	CreatedAt:                    whereHelpernull_Time{field: "\"guild_logging_configs\".\"created_at\""},
	UpdatedAt:                    whereHelpernull_Time{field: "\"guild_logging_configs\".\"updated_at\""},
	UsernameLoggingEnabled:       whereHelpernull_Bool{field: "\"guild_logging_configs\".\"username_logging_enabled\""},
	NicknameLoggingEnabled:       whereHelpernull_Bool{field: "\"guild_logging_configs\".\"nickname_logging_enabled\""},
	BlacklistedChannels:          whereHelpernull_String{field: "\"guild_logging_configs\".\"blacklisted_channels\""},
	ManageMessagesCanViewDeleted: whereHelpernull_Bool{field: "\"guild_logging_configs\".\"manage_messages_can_view_deleted\""},
	EveryoneCanViewDeleted:       whereHelpernull_Bool{field: "\"guild_logging_configs\".\"everyone_can_view_deleted\""},
	MessageLogsAllowedRoles:      whereHelpertypes_Int64Array{field: "\"guild_logging_configs\".\"message_logs_allowed_roles\""},
	AccessMode:                   whereHelperint16{field: "\"guild_logging_configs\".\"access_mode\""},
}

// GuildLoggingConfigRels is where relationship names are stored.
var GuildLoggingConfigRels = struct {
}{}

// guildLoggingConfigR is where relationships are stored.
type guildLoggingConfigR struct {
}

// NewStruct creates a new relationship struct
func (*guildLoggingConfigR) NewStruct() *guildLoggingConfigR {
	return &guildLoggingConfigR{}
}

// guildLoggingConfigL is where Load methods for each relationship are stored.
type guildLoggingConfigL struct{}

var (
	guildLoggingConfigAllColumns            = []string{"guild_id", "created_at", "updated_at", "username_logging_enabled", "nickname_logging_enabled", "blacklisted_channels", "manage_messages_can_view_deleted", "everyone_can_view_deleted", "message_logs_allowed_roles", "access_mode"}
	guildLoggingConfigColumnsWithoutDefault = []string{"guild_id", "created_at", "updated_at", "username_logging_enabled", "nickname_logging_enabled", "blacklisted_channels", "manage_messages_can_view_deleted", "everyone_can_view_deleted", "message_logs_allowed_roles"}
	guildLoggingConfigColumnsWithDefault    = []string{"access_mode"}
	guildLoggingConfigPrimaryKeyColumns     = []string{"guild_id"}
)

type (
	// GuildLoggingConfigSlice is an alias for a slice of pointers to GuildLoggingConfig.
	// This should generally be used opposed to []GuildLoggingConfig.
	GuildLoggingConfigSlice []*GuildLoggingConfig

	guildLoggingConfigQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	guildLoggingConfigType                 = reflect.TypeOf(&GuildLoggingConfig{})
	guildLoggingConfigMapping              = queries.MakeStructMapping(guildLoggingConfigType)
	guildLoggingConfigPrimaryKeyMapping, _ = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, guildLoggingConfigPrimaryKeyColumns)
	guildLoggingConfigInsertCacheMut       sync.RWMutex
	guildLoggingConfigInsertCache          = make(map[string]insertCache)
	guildLoggingConfigUpdateCacheMut       sync.RWMutex
	guildLoggingConfigUpdateCache          = make(map[string]updateCache)
	guildLoggingConfigUpsertCacheMut       sync.RWMutex
	guildLoggingConfigUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single guildLoggingConfig record from the query using the global executor.
func (q guildLoggingConfigQuery) OneG(ctx context.Context) (*GuildLoggingConfig, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single guildLoggingConfig record from the query.
func (q guildLoggingConfigQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GuildLoggingConfig, error) {
	o := &GuildLoggingConfig{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for guild_logging_configs")
	}

	return o, nil
}

// AllG returns all GuildLoggingConfig records from the query using the global executor.
func (q guildLoggingConfigQuery) AllG(ctx context.Context) (GuildLoggingConfigSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all GuildLoggingConfig records from the query.
func (q guildLoggingConfigQuery) All(ctx context.Context, exec boil.ContextExecutor) (GuildLoggingConfigSlice, error) {
	var o []*GuildLoggingConfig

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to GuildLoggingConfig slice")
	}

	return o, nil
}

// CountG returns the count of all GuildLoggingConfig records in the query, and panics on error.
func (q guildLoggingConfigQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all GuildLoggingConfig records in the query.
func (q guildLoggingConfigQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count guild_logging_configs rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q guildLoggingConfigQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q guildLoggingConfigQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if guild_logging_configs exists")
	}

	return count > 0, nil
}

// GuildLoggingConfigs retrieves all the records using an executor.
func GuildLoggingConfigs(mods ...qm.QueryMod) guildLoggingConfigQuery {
	mods = append(mods, qm.From("\"guild_logging_configs\""))
	return guildLoggingConfigQuery{NewQuery(mods...)}
}

// FindGuildLoggingConfigG retrieves a single record by ID.
func FindGuildLoggingConfigG(ctx context.Context, guildID int64, selectCols ...string) (*GuildLoggingConfig, error) {
	return FindGuildLoggingConfig(ctx, boil.GetContextDB(), guildID, selectCols...)
}

// FindGuildLoggingConfig retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGuildLoggingConfig(ctx context.Context, exec boil.ContextExecutor, guildID int64, selectCols ...string) (*GuildLoggingConfig, error) {
	guildLoggingConfigObj := &GuildLoggingConfig{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"guild_logging_configs\" where \"guild_id\"=$1", sel,
	)

	q := queries.Raw(query, guildID)

	err := q.Bind(ctx, exec, guildLoggingConfigObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from guild_logging_configs")
	}

	return guildLoggingConfigObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *GuildLoggingConfig) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GuildLoggingConfig) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no guild_logging_configs provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(guildLoggingConfigColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	guildLoggingConfigInsertCacheMut.RLock()
	cache, cached := guildLoggingConfigInsertCache[key]
	guildLoggingConfigInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			guildLoggingConfigAllColumns,
			guildLoggingConfigColumnsWithDefault,
			guildLoggingConfigColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"guild_logging_configs\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"guild_logging_configs\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into guild_logging_configs")
	}

	if !cached {
		guildLoggingConfigInsertCacheMut.Lock()
		guildLoggingConfigInsertCache[key] = cache
		guildLoggingConfigInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single GuildLoggingConfig record using the global executor.
// See Update for more documentation.
func (o *GuildLoggingConfig) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the GuildLoggingConfig.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GuildLoggingConfig) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	key := makeCacheKey(columns, nil)
	guildLoggingConfigUpdateCacheMut.RLock()
	cache, cached := guildLoggingConfigUpdateCache[key]
	guildLoggingConfigUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			guildLoggingConfigAllColumns,
			guildLoggingConfigPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update guild_logging_configs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"guild_logging_configs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, guildLoggingConfigPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, append(wl, guildLoggingConfigPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update guild_logging_configs row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for guild_logging_configs")
	}

	if !cached {
		guildLoggingConfigUpdateCacheMut.Lock()
		guildLoggingConfigUpdateCache[key] = cache
		guildLoggingConfigUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q guildLoggingConfigQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q guildLoggingConfigQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for guild_logging_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for guild_logging_configs")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o GuildLoggingConfigSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GuildLoggingConfigSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildLoggingConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"guild_logging_configs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, guildLoggingConfigPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in guildLoggingConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all guildLoggingConfig")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *GuildLoggingConfig) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GuildLoggingConfig) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no guild_logging_configs provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	nzDefaults := queries.NonZeroDefaultSet(guildLoggingConfigColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	guildLoggingConfigUpsertCacheMut.RLock()
	cache, cached := guildLoggingConfigUpsertCache[key]
	guildLoggingConfigUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			guildLoggingConfigAllColumns,
			guildLoggingConfigColumnsWithDefault,
			guildLoggingConfigColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			guildLoggingConfigAllColumns,
			guildLoggingConfigPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert guild_logging_configs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(guildLoggingConfigPrimaryKeyColumns))
			copy(conflict, guildLoggingConfigPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"guild_logging_configs\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(guildLoggingConfigType, guildLoggingConfigMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert guild_logging_configs")
	}

	if !cached {
		guildLoggingConfigUpsertCacheMut.Lock()
		guildLoggingConfigUpsertCache[key] = cache
		guildLoggingConfigUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single GuildLoggingConfig record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *GuildLoggingConfig) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single GuildLoggingConfig record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GuildLoggingConfig) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no GuildLoggingConfig provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), guildLoggingConfigPrimaryKeyMapping)
	sql := "DELETE FROM \"guild_logging_configs\" WHERE \"guild_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from guild_logging_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for guild_logging_configs")
	}

	return rowsAff, nil
}

func (q guildLoggingConfigQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q guildLoggingConfigQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no guildLoggingConfigQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from guild_logging_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for guild_logging_configs")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o GuildLoggingConfigSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GuildLoggingConfigSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildLoggingConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"guild_logging_configs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, guildLoggingConfigPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from guildLoggingConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for guild_logging_configs")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *GuildLoggingConfig) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no GuildLoggingConfig provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GuildLoggingConfig) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGuildLoggingConfig(ctx, exec, o.GuildID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GuildLoggingConfigSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty GuildLoggingConfigSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GuildLoggingConfigSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GuildLoggingConfigSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildLoggingConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"guild_logging_configs\".* FROM \"guild_logging_configs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, guildLoggingConfigPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GuildLoggingConfigSlice")
	}

	*o = slice

	return nil
}

// GuildLoggingConfigExistsG checks if the GuildLoggingConfig row exists.
func GuildLoggingConfigExistsG(ctx context.Context, guildID int64) (bool, error) {
	return GuildLoggingConfigExists(ctx, boil.GetContextDB(), guildID)
}

// GuildLoggingConfigExists checks if the GuildLoggingConfig row exists.
func GuildLoggingConfigExists(ctx context.Context, exec boil.ContextExecutor, guildID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"guild_logging_configs\" where \"guild_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, guildID)
	}
	row := exec.QueryRowContext(ctx, sql, guildID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if guild_logging_configs exists")
	}

	return exists, nil
}
