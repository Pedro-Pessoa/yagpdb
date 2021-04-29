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

// PremiumCode is an object representing the database table.
type PremiumCode struct {
	ID        int64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	Code      string     `boil:"code" json:"code" toml:"code" yaml:"code"`
	Message   string     `boil:"message" json:"message" toml:"message" yaml:"message"`
	CreatedAt time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UsedAt    null.Time  `boil:"used_at" json:"used_at,omitempty" toml:"used_at" yaml:"used_at,omitempty"`
	SlotID    null.Int64 `boil:"slot_id" json:"slot_id,omitempty" toml:"slot_id" yaml:"slot_id,omitempty"`
	UserID    null.Int64 `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	GuildID   null.Int64 `boil:"guild_id" json:"guild_id,omitempty" toml:"guild_id" yaml:"guild_id,omitempty"`
	Permanent bool       `boil:"permanent" json:"permanent" toml:"permanent" yaml:"permanent"`
	Duration  int64      `boil:"duration" json:"duration" toml:"duration" yaml:"duration"`

	R *premiumCodeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L premiumCodeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PremiumCodeColumns = struct {
	ID        string
	Code      string
	Message   string
	CreatedAt string
	UsedAt    string
	SlotID    string
	UserID    string
	GuildID   string
	Permanent string
	Duration  string
}{
	ID:        "id",
	Code:      "code",
	Message:   "message",
	CreatedAt: "created_at",
	UsedAt:    "used_at",
	SlotID:    "slot_id",
	UserID:    "user_id",
	GuildID:   "guild_id",
	Permanent: "permanent",
	Duration:  "duration",
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

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
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

type whereHelpernull_Int64 struct{ field string }

func (w whereHelpernull_Int64) EQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int64) NEQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int64) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int64) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Int64) LT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int64) LTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int64) GT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int64) GTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

var PremiumCodeWhere = struct {
	ID        whereHelperint64
	Code      whereHelperstring
	Message   whereHelperstring
	CreatedAt whereHelpertime_Time
	UsedAt    whereHelpernull_Time
	SlotID    whereHelpernull_Int64
	UserID    whereHelpernull_Int64
	GuildID   whereHelpernull_Int64
	Permanent whereHelperbool
	Duration  whereHelperint64
}{
	ID:        whereHelperint64{field: "\"premium_codes\".\"id\""},
	Code:      whereHelperstring{field: "\"premium_codes\".\"code\""},
	Message:   whereHelperstring{field: "\"premium_codes\".\"message\""},
	CreatedAt: whereHelpertime_Time{field: "\"premium_codes\".\"created_at\""},
	UsedAt:    whereHelpernull_Time{field: "\"premium_codes\".\"used_at\""},
	SlotID:    whereHelpernull_Int64{field: "\"premium_codes\".\"slot_id\""},
	UserID:    whereHelpernull_Int64{field: "\"premium_codes\".\"user_id\""},
	GuildID:   whereHelpernull_Int64{field: "\"premium_codes\".\"guild_id\""},
	Permanent: whereHelperbool{field: "\"premium_codes\".\"permanent\""},
	Duration:  whereHelperint64{field: "\"premium_codes\".\"duration\""},
}

// PremiumCodeRels is where relationship names are stored.
var PremiumCodeRels = struct {
	Slot string
}{
	Slot: "Slot",
}

// premiumCodeR is where relationships are stored.
type premiumCodeR struct {
	Slot *PremiumSlot `boil:"Slot" json:"Slot" toml:"Slot" yaml:"Slot"`
}

// NewStruct creates a new relationship struct
func (*premiumCodeR) NewStruct() *premiumCodeR {
	return &premiumCodeR{}
}

// premiumCodeL is where Load methods for each relationship are stored.
type premiumCodeL struct{}

var (
	premiumCodeAllColumns            = []string{"id", "code", "message", "created_at", "used_at", "slot_id", "user_id", "guild_id", "permanent", "duration"}
	premiumCodeColumnsWithoutDefault = []string{"code", "message", "created_at", "used_at", "slot_id", "user_id", "guild_id", "permanent", "duration"}
	premiumCodeColumnsWithDefault    = []string{"id"}
	premiumCodePrimaryKeyColumns     = []string{"id"}
)

type (
	// PremiumCodeSlice is an alias for a slice of pointers to PremiumCode.
	// This should generally be used opposed to []PremiumCode.
	PremiumCodeSlice []*PremiumCode

	premiumCodeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	premiumCodeType                 = reflect.TypeOf(&PremiumCode{})
	premiumCodeMapping              = queries.MakeStructMapping(premiumCodeType)
	premiumCodePrimaryKeyMapping, _ = queries.BindMapping(premiumCodeType, premiumCodeMapping, premiumCodePrimaryKeyColumns)
	premiumCodeInsertCacheMut       sync.RWMutex
	premiumCodeInsertCache          = make(map[string]insertCache)
	premiumCodeUpdateCacheMut       sync.RWMutex
	premiumCodeUpdateCache          = make(map[string]updateCache)
	premiumCodeUpsertCacheMut       sync.RWMutex
	premiumCodeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single premiumCode record from the query using the global executor.
func (q premiumCodeQuery) OneG(ctx context.Context) (*PremiumCode, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single premiumCode record from the query.
func (q premiumCodeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PremiumCode, error) {
	o := &PremiumCode{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for premium_codes")
	}

	return o, nil
}

// AllG returns all PremiumCode records from the query using the global executor.
func (q premiumCodeQuery) AllG(ctx context.Context) (PremiumCodeSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all PremiumCode records from the query.
func (q premiumCodeQuery) All(ctx context.Context, exec boil.ContextExecutor) (PremiumCodeSlice, error) {
	var o []*PremiumCode

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PremiumCode slice")
	}

	return o, nil
}

// CountG returns the count of all PremiumCode records in the query, and panics on error.
func (q premiumCodeQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all PremiumCode records in the query.
func (q premiumCodeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count premium_codes rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q premiumCodeQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q premiumCodeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if premium_codes exists")
	}

	return count > 0, nil
}

// Slot pointed to by the foreign key.
func (o *PremiumCode) Slot(mods ...qm.QueryMod) premiumSlotQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SlotID),
	}

	queryMods = append(queryMods, mods...)

	query := PremiumSlots(queryMods...)
	queries.SetFrom(query.Query, "\"premium_slots\"")

	return query
}

// LoadSlot allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (premiumCodeL) LoadSlot(ctx context.Context, e boil.ContextExecutor, singular bool, maybePremiumCode interface{}, mods queries.Applicator) error {
	var slice []*PremiumCode
	var object *PremiumCode

	if singular {
		object = maybePremiumCode.(*PremiumCode)
	} else {
		slice = *maybePremiumCode.(*[]*PremiumCode)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &premiumCodeR{}
		}
		if !queries.IsNil(object.SlotID) {
			args = append(args, object.SlotID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &premiumCodeR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.SlotID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.SlotID) {
				args = append(args, obj.SlotID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`premium_slots`),
		qm.WhereIn(`premium_slots.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load PremiumSlot")
	}

	var resultSlice []*PremiumSlot
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice PremiumSlot")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for premium_slots")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for premium_slots")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Slot = foreign
		if foreign.R == nil {
			foreign.R = &premiumSlotR{}
		}
		foreign.R.SlotPremiumCodes = append(foreign.R.SlotPremiumCodes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.SlotID, foreign.ID) {
				local.R.Slot = foreign
				if foreign.R == nil {
					foreign.R = &premiumSlotR{}
				}
				foreign.R.SlotPremiumCodes = append(foreign.R.SlotPremiumCodes, local)
				break
			}
		}
	}

	return nil
}

// SetSlotG of the premiumCode to the related item.
// Sets o.R.Slot to related.
// Adds o to related.R.SlotPremiumCodes.
// Uses the global database handle.
func (o *PremiumCode) SetSlotG(ctx context.Context, insert bool, related *PremiumSlot) error {
	return o.SetSlot(ctx, boil.GetContextDB(), insert, related)
}

// SetSlot of the premiumCode to the related item.
// Sets o.R.Slot to related.
// Adds o to related.R.SlotPremiumCodes.
func (o *PremiumCode) SetSlot(ctx context.Context, exec boil.ContextExecutor, insert bool, related *PremiumSlot) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"premium_codes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"slot_id"}),
		strmangle.WhereClause("\"", "\"", 2, premiumCodePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.SlotID, related.ID)
	if o.R == nil {
		o.R = &premiumCodeR{
			Slot: related,
		}
	} else {
		o.R.Slot = related
	}

	if related.R == nil {
		related.R = &premiumSlotR{
			SlotPremiumCodes: PremiumCodeSlice{o},
		}
	} else {
		related.R.SlotPremiumCodes = append(related.R.SlotPremiumCodes, o)
	}

	return nil
}

// RemoveSlotG relationship.
// Sets o.R.Slot to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *PremiumCode) RemoveSlotG(ctx context.Context, related *PremiumSlot) error {
	return o.RemoveSlot(ctx, boil.GetContextDB(), related)
}

// RemoveSlot relationship.
// Sets o.R.Slot to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *PremiumCode) RemoveSlot(ctx context.Context, exec boil.ContextExecutor, related *PremiumSlot) error {
	var err error

	queries.SetScanner(&o.SlotID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("slot_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Slot = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.SlotPremiumCodes {
		if queries.Equal(o.SlotID, ri.SlotID) {
			continue
		}

		ln := len(related.R.SlotPremiumCodes)
		if ln > 1 && i < ln-1 {
			related.R.SlotPremiumCodes[i] = related.R.SlotPremiumCodes[ln-1]
		}
		related.R.SlotPremiumCodes = related.R.SlotPremiumCodes[:ln-1]
		break
	}
	return nil
}

// PremiumCodes retrieves all the records using an executor.
func PremiumCodes(mods ...qm.QueryMod) premiumCodeQuery {
	mods = append(mods, qm.From("\"premium_codes\""))
	return premiumCodeQuery{NewQuery(mods...)}
}

// FindPremiumCodeG retrieves a single record by ID.
func FindPremiumCodeG(ctx context.Context, iD int64, selectCols ...string) (*PremiumCode, error) {
	return FindPremiumCode(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindPremiumCode retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPremiumCode(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*PremiumCode, error) {
	premiumCodeObj := &PremiumCode{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"premium_codes\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, premiumCodeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from premium_codes")
	}

	return premiumCodeObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *PremiumCode) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PremiumCode) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no premium_codes provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(premiumCodeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	premiumCodeInsertCacheMut.RLock()
	cache, cached := premiumCodeInsertCache[key]
	premiumCodeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			premiumCodeAllColumns,
			premiumCodeColumnsWithDefault,
			premiumCodeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(premiumCodeType, premiumCodeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(premiumCodeType, premiumCodeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"premium_codes\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"premium_codes\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into premium_codes")
	}

	if !cached {
		premiumCodeInsertCacheMut.Lock()
		premiumCodeInsertCache[key] = cache
		premiumCodeInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single PremiumCode record using the global executor.
// See Update for more documentation.
func (o *PremiumCode) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the PremiumCode.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PremiumCode) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	premiumCodeUpdateCacheMut.RLock()
	cache, cached := premiumCodeUpdateCache[key]
	premiumCodeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			premiumCodeAllColumns,
			premiumCodePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update premium_codes, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"premium_codes\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, premiumCodePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(premiumCodeType, premiumCodeMapping, append(wl, premiumCodePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update premium_codes row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for premium_codes")
	}

	if !cached {
		premiumCodeUpdateCacheMut.Lock()
		premiumCodeUpdateCache[key] = cache
		premiumCodeUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q premiumCodeQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q premiumCodeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for premium_codes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for premium_codes")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PremiumCodeSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PremiumCodeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumCodePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"premium_codes\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, premiumCodePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in premiumCode slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all premiumCode")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *PremiumCode) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PremiumCode) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no premium_codes provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(premiumCodeColumnsWithDefault, o)

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

	premiumCodeUpsertCacheMut.RLock()
	cache, cached := premiumCodeUpsertCache[key]
	premiumCodeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			premiumCodeAllColumns,
			premiumCodeColumnsWithDefault,
			premiumCodeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			premiumCodeAllColumns,
			premiumCodePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert premium_codes, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(premiumCodePrimaryKeyColumns))
			copy(conflict, premiumCodePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"premium_codes\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(premiumCodeType, premiumCodeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(premiumCodeType, premiumCodeMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert premium_codes")
	}

	if !cached {
		premiumCodeUpsertCacheMut.Lock()
		premiumCodeUpsertCache[key] = cache
		premiumCodeUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single PremiumCode record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *PremiumCode) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single PremiumCode record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PremiumCode) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PremiumCode provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), premiumCodePrimaryKeyMapping)
	sql := "DELETE FROM \"premium_codes\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from premium_codes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for premium_codes")
	}

	return rowsAff, nil
}

func (q premiumCodeQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q premiumCodeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no premiumCodeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from premium_codes")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for premium_codes")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o PremiumCodeSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PremiumCodeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumCodePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"premium_codes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, premiumCodePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from premiumCode slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for premium_codes")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *PremiumCode) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no PremiumCode provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PremiumCode) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPremiumCode(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PremiumCodeSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty PremiumCodeSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PremiumCodeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PremiumCodeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), premiumCodePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"premium_codes\".* FROM \"premium_codes\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, premiumCodePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PremiumCodeSlice")
	}

	*o = slice

	return nil
}

// PremiumCodeExistsG checks if the PremiumCode row exists.
func PremiumCodeExistsG(ctx context.Context, iD int64) (bool, error) {
	return PremiumCodeExists(ctx, boil.GetContextDB(), iD)
}

// PremiumCodeExists checks if the PremiumCode row exists.
func PremiumCodeExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"premium_codes\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if premium_codes exists")
	}

	return exists, nil
}
