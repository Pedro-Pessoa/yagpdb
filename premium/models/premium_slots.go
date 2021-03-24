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
	"github.com/volatiletech/strmangle"
)

// PremiumSlot is an object representing the database table.
type PremiumSlot struct {
	ID                int64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedAt         time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	AttachedAt        null.Time  `boil:"attached_at" json:"attached_at,omitempty" toml:"attached_at" yaml:"attached_at,omitempty"`
	UserID            int64      `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	GuildID           null.Int64 `boil:"guild_id" json:"guild_id,omitempty" toml:"guild_id" yaml:"guild_id,omitempty"`
	Title             string     `boil:"title" json:"title" toml:"title" yaml:"title"`
	Message           string     `boil:"message" json:"message" toml:"message" yaml:"message"`
	Source            string     `boil:"source" json:"source" toml:"source" yaml:"source"`
	SourceID          int64      `boil:"source_id" json:"source_id" toml:"source_id" yaml:"source_id"`
	FullDuration      int64      `boil:"full_duration" json:"full_duration" toml:"full_duration" yaml:"full_duration"`
	Permanent         bool       `boil:"permanent" json:"permanent" toml:"permanent" yaml:"permanent"`
	DurationRemaining int64      `boil:"duration_remaining" json:"duration_remaining" toml:"duration_remaining" yaml:"duration_remaining"`

	R *premiumSlotR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L premiumSlotL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PremiumSlotColumns = struct {
	ID                string
	CreatedAt         string
	AttachedAt        string
	UserID            string
	GuildID           string
	Title             string
	Message           string
	Source            string
	SourceID          string
	FullDuration      string
	Permanent         string
	DurationRemaining string
}{
	ID:                "id",
	CreatedAt:         "created_at",
	AttachedAt:        "attached_at",
	UserID:            "user_id",
	GuildID:           "guild_id",
	Title:             "title",
	Message:           "message",
	Source:            "source",
	SourceID:          "source_id",
	FullDuration:      "full_duration",
	Permanent:         "permanent",
	DurationRemaining: "duration_remaining",
}

// Generated where

var PremiumSlotWhere = struct {
	ID                whereHelperint64
	CreatedAt         whereHelpertime_Time
	AttachedAt        whereHelpernull_Time
	UserID            whereHelperint64
	GuildID           whereHelpernull_Int64
	Title             whereHelperstring
	Message           whereHelperstring
	Source            whereHelperstring
	SourceID          whereHelperint64
	FullDuration      whereHelperint64
	Permanent         whereHelperbool
	DurationRemaining whereHelperint64
}{
	ID:                whereHelperint64{field: "\"premium_slots\".\"id\""},
	CreatedAt:         whereHelpertime_Time{field: "\"premium_slots\".\"created_at\""},
	AttachedAt:        whereHelpernull_Time{field: "\"premium_slots\".\"attached_at\""},
	UserID:            whereHelperint64{field: "\"premium_slots\".\"user_id\""},
	GuildID:           whereHelpernull_Int64{field: "\"premium_slots\".\"guild_id\""},
	Title:             whereHelperstring{field: "\"premium_slots\".\"title\""},
	Message:           whereHelperstring{field: "\"premium_slots\".\"message\""},
	Source:            whereHelperstring{field: "\"premium_slots\".\"source\""},
	SourceID:          whereHelperint64{field: "\"premium_slots\".\"source_id\""},
	FullDuration:      whereHelperint64{field: "\"premium_slots\".\"full_duration\""},
	Permanent:         whereHelperbool{field: "\"premium_slots\".\"permanent\""},
	DurationRemaining: whereHelperint64{field: "\"premium_slots\".\"duration_remaining\""},
}

// PremiumSlotRels is where relationship names are stored.
var PremiumSlotRels = struct {
	SlotPremiumCodes string
}{
	SlotPremiumCodes: "SlotPremiumCodes",
}

// premiumSlotR is where relationships are stored.
type premiumSlotR struct {
	SlotPremiumCodes PremiumCodeSlice `boil:"SlotPremiumCodes" json:"SlotPremiumCodes" toml:"SlotPremiumCodes" yaml:"SlotPremiumCodes"`
}

// NewStruct creates a new relationship struct
func (*premiumSlotR) NewStruct() *premiumSlotR {
	return &premiumSlotR{}
}

// premiumSlotL is where Load methods for each relationship are stored.
type premiumSlotL struct{}

var (
	premiumSlotAllColumns            = []string{"id", "created_at", "attached_at", "user_id", "guild_id", "title", "message", "source", "source_id", "full_duration", "permanent", "duration_remaining"}
	premiumSlotColumnsWithoutDefault = []string{"created_at", "attached_at", "user_id", "guild_id", "title", "message", "source", "source_id", "full_duration", "permanent", "duration_remaining"}
	premiumSlotColumnsWithDefault    = []string{"id"}
	premiumSlotPrimaryKeyColumns     = []string{"id"}
)

type (
	// PremiumSlotSlice is an alias for a slice of pointers to PremiumSlot.
	// This should generally be used opposed to []PremiumSlot.
	PremiumSlotSlice []*PremiumSlot

	premiumSlotQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	premiumSlotType                 = reflect.TypeOf(&PremiumSlot{})
	premiumSlotMapping              = queries.MakeStructMapping(premiumSlotType)
	premiumSlotPrimaryKeyMapping, _ = queries.BindMapping(premiumSlotType, premiumSlotMapping, premiumSlotPrimaryKeyColumns)
	premiumSlotInsertCacheMut       sync.RWMutex
	premiumSlotInsertCache          = make(map[string]insertCache)
	premiumSlotUpdateCacheMut       sync.RWMutex
	premiumSlotUpdateCache          = make(map[string]updateCache)
	premiumSlotUpsertCacheMut       sync.RWMutex
	premiumSlotUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single premiumSlot record from the query using the global executor.
func (q premiumSlotQuery) OneG(ctx context.Context) (*PremiumSlot, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single premiumSlot record from the query.
func (q premiumSlotQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PremiumSlot, error) {
	o := &PremiumSlot{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for premium_slots")
	}

	return o, nil
}

// AllG returns all PremiumSlot records from the query using the global executor.
func (q premiumSlotQuery) AllG(ctx context.Context) (PremiumSlotSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all PremiumSlot records from the query.
func (q premiumSlotQuery) All(ctx context.Context, exec boil.ContextExecutor) (PremiumSlotSlice, error) {
	var o []*PremiumSlot

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PremiumSlot slice")
	}

	return o, nil
}

// CountG returns the count of all PremiumSlot records in the query, and panics on error.
func (q premiumSlotQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all PremiumSlot records in the query.
func (q premiumSlotQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count premium_slots rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q premiumSlotQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q premiumSlotQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if premium_slots exists")
	}

	return count > 0, nil
}

// SlotPremiumCodes retrieves all the premium_code's PremiumCodes with an executor via slot_id column.
func (o *PremiumSlot) SlotPremiumCodes(mods ...qm.QueryMod) premiumCodeQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"premium_codes\".\"slot_id\"=?", o.ID),
	)

	query := PremiumCodes(queryMods...)
	queries.SetFrom(query.Query, "\"premium_codes\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"premium_codes\".*"})
	}

	return query
}

// LoadSlotPremiumCodes allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (premiumSlotL) LoadSlotPremiumCodes(ctx context.Context, e boil.ContextExecutor, singular bool, maybePremiumSlot interface{}, mods queries.Applicator) error {
	var slice []*PremiumSlot
	var object *PremiumSlot

	if singular {
		object = maybePremiumSlot.(*PremiumSlot)
	} else {
		slice = *maybePremiumSlot.(*[]*PremiumSlot)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &premiumSlotR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &premiumSlotR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`premium_codes`),
		qm.WhereIn(`premium_codes.slot_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load premium_codes")
	}

	var resultSlice []*PremiumCode
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice premium_codes")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on premium_codes")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for premium_codes")
	}

	if singular {
		object.R.SlotPremiumCodes = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &premiumCodeR{}
			}
			foreign.R.Slot = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.SlotID) {
				local.R.SlotPremiumCodes = append(local.R.SlotPremiumCodes, foreign)
				if foreign.R == nil {
					foreign.R = &premiumCodeR{}
				}
				foreign.R.Slot = local
				break
			}
		}
	}

	return nil
}

// AddSlotPremiumCodesG adds the given related objects to the existing relationships
// of the premium_slot, optionally inserting them as new records.
// Appends related to o.R.SlotPremiumCodes.
// Sets related.R.Slot appropriately.
// Uses the global database handle.
func (o *PremiumSlot) AddSlotPremiumCodesG(ctx context.Context, insert bool, related ...*PremiumCode) error {
	return o.AddSlotPremiumCodes(ctx, boil.GetContextDB(), insert, related...)
}

// AddSlotPremiumCodes adds the given related objects to the existing relationships
// of the premium_slot, optionally inserting them as new records.
// Appends related to o.R.SlotPremiumCodes.
// Sets related.R.Slot appropriately.
func (o *PremiumSlot) AddSlotPremiumCodes(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*PremiumCode) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.SlotID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"premium_codes\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"slot_id"}),
				strmangle.WhereClause("\"", "\"", 2, premiumCodePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.SlotID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &premiumSlotR{
			SlotPremiumCodes: related,
		}
	} else {
		o.R.SlotPremiumCodes = append(o.R.SlotPremiumCodes, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &premiumCodeR{
				Slot: o,
			}
		} else {
			rel.R.Slot = o
		}
	}
	return nil
}

// SetSlotPremiumCodesG removes all previously related items of the
// premium_slot replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Slot's SlotPremiumCodes accordingly.
// Replaces o.R.SlotPremiumCodes with related.
// Sets related.R.Slot's SlotPremiumCodes accordingly.
// Uses the global database handle.
func (o *PremiumSlot) SetSlotPremiumCodesG(ctx context.Context, insert bool, related ...*PremiumCode) error {
	return o.SetSlotPremiumCodes(ctx, boil.GetContextDB(), insert, related...)
}

// SetSlotPremiumCodes removes all previously related items of the
// premium_slot replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Slot's SlotPremiumCodes accordingly.
// Replaces o.R.SlotPremiumCodes with related.
// Sets related.R.Slot's SlotPremiumCodes accordingly.
func (o *PremiumSlot) SetSlotPremiumCodes(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*PremiumCode) error {
	query := "update \"premium_codes\" set \"slot_id\" = null where \"slot_id\" = $1"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.SlotPremiumCodes {
			queries.SetScanner(&rel.SlotID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Slot = nil
		}

		o.R.SlotPremiumCodes = nil
	}
	return o.AddSlotPremiumCodes(ctx, exec, insert, related...)
}

// RemoveSlotPremiumCodesG relationships from objects passed in.
// Removes related items from R.SlotPremiumCodes (uses pointer comparison, removal does not keep order)
// Sets related.R.Slot.
// Uses the global database handle.
func (o *PremiumSlot) RemoveSlotPremiumCodesG(ctx context.Context, related ...*PremiumCode) error {
	return o.RemoveSlotPremiumCodes(ctx, boil.GetContextDB(), related...)
}

// RemoveSlotPremiumCodes relationships from objects passed in.
// Removes related items from R.SlotPremiumCodes (uses pointer comparison, removal does not keep order)
// Sets related.R.Slot.
func (o *PremiumSlot) RemoveSlotPremiumCodes(ctx context.Context, exec boil.ContextExecutor, related ...*PremiumCode) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.SlotID, nil)
		if rel.R != nil {
			rel.R.Slot = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("slot_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.SlotPremiumCodes {
			if rel != ri {
				continue
			}

			ln := len(o.R.SlotPremiumCodes)
			if ln > 1 && i < ln-1 {
				o.R.SlotPremiumCodes[i] = o.R.SlotPremiumCodes[ln-1]
			}
			o.R.SlotPremiumCodes = o.R.SlotPremiumCodes[:ln-1]
			break
		}
	}

	return nil
}

// PremiumSlots retrieves all the records using an executor.
func PremiumSlots(mods ...qm.QueryMod) premiumSlotQuery {
	mods = append(mods, qm.From("\"premium_slots\""))
	return premiumSlotQuery{NewQuery(mods...)}
}

// FindPremiumSlotG retrieves a single record by ID.
func FindPremiumSlotG(ctx context.Context, iD int64, selectCols ...string) (*PremiumSlot, error) {
	return FindPremiumSlot(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindPremiumSlot retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPremiumSlot(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*PremiumSlot, error) {
	premiumSlotObj := &PremiumSlot{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"premium_slots\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, premiumSlotObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from premium_slots")
	}

	return premiumSlotObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *PremiumSlot) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PremiumSlot) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no premium_slots provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(premiumSlotColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	premiumSlotInsertCacheMut.RLock()
	cache, cached := premiumSlotInsertCache[key]
	premiumSlotInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			premiumSlotAllColumns,
			premiumSlotColumnsWithDefault,
			premiumSlotColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(premiumSlotType, premiumSlotMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(premiumSlotType, premiumSlotMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"premium_slots\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"premium_slots\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into premium_slots")
	}

	if !cached {
		premiumSlotInsertCacheMut.Lock()
		premiumSlotInsertCache[key] = cache
		premiumSlotInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single PremiumSlot record using the global executor.
// See Update for more documentation.
func (o *PremiumSlot) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the PremiumSlot.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PremiumSlot) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	premiumSlotUpdateCacheMut.RLock()
	cache, cached := premiumSlotUpdateCache[key]
	premiumSlotUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			premiumSlotAllColumns,
			premiumSlotPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update premium_slots, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"premium_slots\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, premiumSlotPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(premiumSlotType, premiumSlotMapping, append(wl, premiumSlotPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update premium_slots row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for premium_slots")
	}

	if !cached {
		premiumSlotUpdateCacheMut.Lock()
		premiumSlotUpdateCache[key] = cache
		premiumSlotUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q premiumSlotQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q premiumSlotQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for premium_slots")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for premium_slots")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PremiumSlotSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PremiumSlotSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumSlotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"premium_slots\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, premiumSlotPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in premiumSlot slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all premiumSlot")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *PremiumSlot) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PremiumSlot) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no premium_slots provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(premiumSlotColumnsWithDefault, o)

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

	premiumSlotUpsertCacheMut.RLock()
	cache, cached := premiumSlotUpsertCache[key]
	premiumSlotUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			premiumSlotAllColumns,
			premiumSlotColumnsWithDefault,
			premiumSlotColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			premiumSlotAllColumns,
			premiumSlotPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert premium_slots, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(premiumSlotPrimaryKeyColumns))
			copy(conflict, premiumSlotPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"premium_slots\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(premiumSlotType, premiumSlotMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(premiumSlotType, premiumSlotMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert premium_slots")
	}

	if !cached {
		premiumSlotUpsertCacheMut.Lock()
		premiumSlotUpsertCache[key] = cache
		premiumSlotUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single PremiumSlot record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *PremiumSlot) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single PremiumSlot record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PremiumSlot) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PremiumSlot provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), premiumSlotPrimaryKeyMapping)
	sql := "DELETE FROM \"premium_slots\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from premium_slots")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for premium_slots")
	}

	return rowsAff, nil
}

func (q premiumSlotQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q premiumSlotQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no premiumSlotQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from premium_slots")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for premium_slots")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o PremiumSlotSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PremiumSlotSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumSlotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"premium_slots\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, premiumSlotPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from premiumSlot slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for premium_slots")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *PremiumSlot) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no PremiumSlot provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PremiumSlot) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPremiumSlot(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PremiumSlotSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty PremiumSlotSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PremiumSlotSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PremiumSlotSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumSlotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"premium_slots\".* FROM \"premium_slots\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, premiumSlotPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PremiumSlotSlice")
	}

	*o = slice

	return nil
}

// PremiumSlotExistsG checks if the PremiumSlot row exists.
func PremiumSlotExistsG(ctx context.Context, iD int64) (bool, error) {
	return PremiumSlotExists(ctx, boil.GetContextDB(), iD)
}

// PremiumSlotExists checks if the PremiumSlot row exists.
func PremiumSlotExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"premium_slots\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if premium_slots exists")
	}

	return exists, nil
}
