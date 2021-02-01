// Code generated by SQLBoiler 4.4.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// VerificationConfig is an object representing the database table.
type VerificationConfig struct {
	GuildID             int64  `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	Enabled             bool   `boil:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	VerifiedRole        int64  `boil:"verified_role" json:"verified_role" toml:"verified_role" yaml:"verified_role"`
	PageContent         string `boil:"page_content" json:"page_content" toml:"page_content" yaml:"page_content"`
	KickUnverifiedAfter int    `boil:"kick_unverified_after" json:"kick_unverified_after" toml:"kick_unverified_after" yaml:"kick_unverified_after"`
	WarnUnverifiedAfter int    `boil:"warn_unverified_after" json:"warn_unverified_after" toml:"warn_unverified_after" yaml:"warn_unverified_after"`
	WarnMessage         string `boil:"warn_message" json:"warn_message" toml:"warn_message" yaml:"warn_message"`
	LogChannel          int64  `boil:"log_channel" json:"log_channel" toml:"log_channel" yaml:"log_channel"`
	DMMessage           string `boil:"dm_message" json:"dm_message" toml:"dm_message" yaml:"dm_message"`

	R *verificationConfigR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L verificationConfigL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VerificationConfigColumns = struct {
	GuildID             string
	Enabled             string
	VerifiedRole        string
	PageContent         string
	KickUnverifiedAfter string
	WarnUnverifiedAfter string
	WarnMessage         string
	LogChannel          string
	DMMessage           string
}{
	GuildID:             "guild_id",
	Enabled:             "enabled",
	VerifiedRole:        "verified_role",
	PageContent:         "page_content",
	KickUnverifiedAfter: "kick_unverified_after",
	WarnUnverifiedAfter: "warn_unverified_after",
	WarnMessage:         "warn_message",
	LogChannel:          "log_channel",
	DMMessage:           "dm_message",
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

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var VerificationConfigWhere = struct {
	GuildID             whereHelperint64
	Enabled             whereHelperbool
	VerifiedRole        whereHelperint64
	PageContent         whereHelperstring
	KickUnverifiedAfter whereHelperint
	WarnUnverifiedAfter whereHelperint
	WarnMessage         whereHelperstring
	LogChannel          whereHelperint64
	DMMessage           whereHelperstring
}{
	GuildID:             whereHelperint64{field: "\"verification_configs\".\"guild_id\""},
	Enabled:             whereHelperbool{field: "\"verification_configs\".\"enabled\""},
	VerifiedRole:        whereHelperint64{field: "\"verification_configs\".\"verified_role\""},
	PageContent:         whereHelperstring{field: "\"verification_configs\".\"page_content\""},
	KickUnverifiedAfter: whereHelperint{field: "\"verification_configs\".\"kick_unverified_after\""},
	WarnUnverifiedAfter: whereHelperint{field: "\"verification_configs\".\"warn_unverified_after\""},
	WarnMessage:         whereHelperstring{field: "\"verification_configs\".\"warn_message\""},
	LogChannel:          whereHelperint64{field: "\"verification_configs\".\"log_channel\""},
	DMMessage:           whereHelperstring{field: "\"verification_configs\".\"dm_message\""},
}

// VerificationConfigRels is where relationship names are stored.
var VerificationConfigRels = struct {
}{}

// verificationConfigR is where relationships are stored.
type verificationConfigR struct {
}

// NewStruct creates a new relationship struct
func (*verificationConfigR) NewStruct() *verificationConfigR {
	return &verificationConfigR{}
}

// verificationConfigL is where Load methods for each relationship are stored.
type verificationConfigL struct{}

var (
	verificationConfigAllColumns            = []string{"guild_id", "enabled", "verified_role", "page_content", "kick_unverified_after", "warn_unverified_after", "warn_message", "log_channel", "dm_message"}
	verificationConfigColumnsWithoutDefault = []string{"guild_id", "enabled", "verified_role", "page_content", "kick_unverified_after", "warn_unverified_after", "warn_message", "log_channel"}
	verificationConfigColumnsWithDefault    = []string{"dm_message"}
	verificationConfigPrimaryKeyColumns     = []string{"guild_id"}
)

type (
	// VerificationConfigSlice is an alias for a slice of pointers to VerificationConfig.
	// This should generally be used opposed to []VerificationConfig.
	VerificationConfigSlice []*VerificationConfig

	verificationConfigQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	verificationConfigType                 = reflect.TypeOf(&VerificationConfig{})
	verificationConfigMapping              = queries.MakeStructMapping(verificationConfigType)
	verificationConfigPrimaryKeyMapping, _ = queries.BindMapping(verificationConfigType, verificationConfigMapping, verificationConfigPrimaryKeyColumns)
	verificationConfigInsertCacheMut       sync.RWMutex
	verificationConfigInsertCache          = make(map[string]insertCache)
	verificationConfigUpdateCacheMut       sync.RWMutex
	verificationConfigUpdateCache          = make(map[string]updateCache)
	verificationConfigUpsertCacheMut       sync.RWMutex
	verificationConfigUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single verificationConfig record from the query using the global executor.
func (q verificationConfigQuery) OneG(ctx context.Context) (*VerificationConfig, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single verificationConfig record from the query.
func (q verificationConfigQuery) One(ctx context.Context, exec boil.ContextExecutor) (*VerificationConfig, error) {
	o := &VerificationConfig{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for verification_configs")
	}

	return o, nil
}

// AllG returns all VerificationConfig records from the query using the global executor.
func (q verificationConfigQuery) AllG(ctx context.Context) (VerificationConfigSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all VerificationConfig records from the query.
func (q verificationConfigQuery) All(ctx context.Context, exec boil.ContextExecutor) (VerificationConfigSlice, error) {
	var o []*VerificationConfig

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to VerificationConfig slice")
	}

	return o, nil
}

// CountG returns the count of all VerificationConfig records in the query, and panics on error.
func (q verificationConfigQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all VerificationConfig records in the query.
func (q verificationConfigQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count verification_configs rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q verificationConfigQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q verificationConfigQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if verification_configs exists")
	}

	return count > 0, nil
}

// VerificationConfigs retrieves all the records using an executor.
func VerificationConfigs(mods ...qm.QueryMod) verificationConfigQuery {
	mods = append(mods, qm.From("\"verification_configs\""))
	return verificationConfigQuery{NewQuery(mods...)}
}

// FindVerificationConfigG retrieves a single record by ID.
func FindVerificationConfigG(ctx context.Context, guildID int64, selectCols ...string) (*VerificationConfig, error) {
	return FindVerificationConfig(ctx, boil.GetContextDB(), guildID, selectCols...)
}

// FindVerificationConfig retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVerificationConfig(ctx context.Context, exec boil.ContextExecutor, guildID int64, selectCols ...string) (*VerificationConfig, error) {
	verificationConfigObj := &VerificationConfig{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"verification_configs\" where \"guild_id\"=$1", sel,
	)

	q := queries.Raw(query, guildID)

	err := q.Bind(ctx, exec, verificationConfigObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from verification_configs")
	}

	return verificationConfigObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *VerificationConfig) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VerificationConfig) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no verification_configs provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(verificationConfigColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	verificationConfigInsertCacheMut.RLock()
	cache, cached := verificationConfigInsertCache[key]
	verificationConfigInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			verificationConfigAllColumns,
			verificationConfigColumnsWithDefault,
			verificationConfigColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(verificationConfigType, verificationConfigMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(verificationConfigType, verificationConfigMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"verification_configs\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"verification_configs\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into verification_configs")
	}

	if !cached {
		verificationConfigInsertCacheMut.Lock()
		verificationConfigInsertCache[key] = cache
		verificationConfigInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single VerificationConfig record using the global executor.
// See Update for more documentation.
func (o *VerificationConfig) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the VerificationConfig.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VerificationConfig) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	verificationConfigUpdateCacheMut.RLock()
	cache, cached := verificationConfigUpdateCache[key]
	verificationConfigUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			verificationConfigAllColumns,
			verificationConfigPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update verification_configs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"verification_configs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, verificationConfigPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(verificationConfigType, verificationConfigMapping, append(wl, verificationConfigPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update verification_configs row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for verification_configs")
	}

	if !cached {
		verificationConfigUpdateCacheMut.Lock()
		verificationConfigUpdateCache[key] = cache
		verificationConfigUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q verificationConfigQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q verificationConfigQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for verification_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for verification_configs")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o VerificationConfigSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VerificationConfigSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), verificationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"verification_configs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, verificationConfigPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in verificationConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all verificationConfig")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *VerificationConfig) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VerificationConfig) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no verification_configs provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(verificationConfigColumnsWithDefault, o)

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

	verificationConfigUpsertCacheMut.RLock()
	cache, cached := verificationConfigUpsertCache[key]
	verificationConfigUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			verificationConfigAllColumns,
			verificationConfigColumnsWithDefault,
			verificationConfigColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			verificationConfigAllColumns,
			verificationConfigPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert verification_configs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(verificationConfigPrimaryKeyColumns))
			copy(conflict, verificationConfigPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"verification_configs\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(verificationConfigType, verificationConfigMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(verificationConfigType, verificationConfigMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert verification_configs")
	}

	if !cached {
		verificationConfigUpsertCacheMut.Lock()
		verificationConfigUpsertCache[key] = cache
		verificationConfigUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single VerificationConfig record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *VerificationConfig) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single VerificationConfig record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VerificationConfig) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no VerificationConfig provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), verificationConfigPrimaryKeyMapping)
	sql := "DELETE FROM \"verification_configs\" WHERE \"guild_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from verification_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for verification_configs")
	}

	return rowsAff, nil
}

func (q verificationConfigQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q verificationConfigQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no verificationConfigQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from verification_configs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for verification_configs")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o VerificationConfigSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VerificationConfigSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), verificationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"verification_configs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, verificationConfigPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from verificationConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for verification_configs")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *VerificationConfig) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no VerificationConfig provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *VerificationConfig) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVerificationConfig(ctx, exec, o.GuildID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VerificationConfigSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty VerificationConfigSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VerificationConfigSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VerificationConfigSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), verificationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"verification_configs\".* FROM \"verification_configs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, verificationConfigPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in VerificationConfigSlice")
	}

	*o = slice

	return nil
}

// VerificationConfigExistsG checks if the VerificationConfig row exists.
func VerificationConfigExistsG(ctx context.Context, guildID int64) (bool, error) {
	return VerificationConfigExists(ctx, boil.GetContextDB(), guildID)
}

// VerificationConfigExists checks if the VerificationConfig row exists.
func VerificationConfigExists(ctx context.Context, exec boil.ContextExecutor, guildID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"verification_configs\" where \"guild_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, guildID)
	}
	row := exec.QueryRowContext(ctx, sql, guildID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if verification_configs exists")
	}

	return exists, nil
}
