// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package schema

import (
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
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// DjangoCeleryBeatSolarschedule is an object representing the database table.
type DjangoCeleryBeatSolarschedule struct {
	ID        int           `boil:"id" json:"id" toml:"id" yaml:"id"`
	Event     string        `boil:"event" json:"event" toml:"event" yaml:"event"`
	Latitude  types.Decimal `boil:"latitude" json:"latitude" toml:"latitude" yaml:"latitude"`
	Longitude types.Decimal `boil:"longitude" json:"longitude" toml:"longitude" yaml:"longitude"`

	R *djangoCeleryBeatSolarscheduleR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L djangoCeleryBeatSolarscheduleL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var DjangoCeleryBeatSolarscheduleColumns = struct {
	ID        string
	Event     string
	Latitude  string
	Longitude string
}{
	ID:        "id",
	Event:     "event",
	Latitude:  "latitude",
	Longitude: "longitude",
}

var DjangoCeleryBeatSolarscheduleTableColumns = struct {
	ID        string
	Event     string
	Latitude  string
	Longitude string
}{
	ID:        "django_celery_beat_solarschedule.id",
	Event:     "django_celery_beat_solarschedule.event",
	Latitude:  "django_celery_beat_solarschedule.latitude",
	Longitude: "django_celery_beat_solarschedule.longitude",
}

// Generated where

type whereHelpertypes_Decimal struct{ field string }

func (w whereHelpertypes_Decimal) EQ(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_Decimal) NEQ(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_Decimal) LT(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_Decimal) LTE(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_Decimal) GT(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_Decimal) GTE(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var DjangoCeleryBeatSolarscheduleWhere = struct {
	ID        whereHelperint
	Event     whereHelperstring
	Latitude  whereHelpertypes_Decimal
	Longitude whereHelpertypes_Decimal
}{
	ID:        whereHelperint{field: "\"django_celery_beat_solarschedule\".\"id\""},
	Event:     whereHelperstring{field: "\"django_celery_beat_solarschedule\".\"event\""},
	Latitude:  whereHelpertypes_Decimal{field: "\"django_celery_beat_solarschedule\".\"latitude\""},
	Longitude: whereHelpertypes_Decimal{field: "\"django_celery_beat_solarschedule\".\"longitude\""},
}

// DjangoCeleryBeatSolarscheduleRels is where relationship names are stored.
var DjangoCeleryBeatSolarscheduleRels = struct {
	SolarDjangoCeleryBeatPeriodictasks string
}{
	SolarDjangoCeleryBeatPeriodictasks: "SolarDjangoCeleryBeatPeriodictasks",
}

// djangoCeleryBeatSolarscheduleR is where relationships are stored.
type djangoCeleryBeatSolarscheduleR struct {
	SolarDjangoCeleryBeatPeriodictasks DjangoCeleryBeatPeriodictaskSlice `boil:"SolarDjangoCeleryBeatPeriodictasks" json:"SolarDjangoCeleryBeatPeriodictasks" toml:"SolarDjangoCeleryBeatPeriodictasks" yaml:"SolarDjangoCeleryBeatPeriodictasks"`
}

// NewStruct creates a new relationship struct
func (*djangoCeleryBeatSolarscheduleR) NewStruct() *djangoCeleryBeatSolarscheduleR {
	return &djangoCeleryBeatSolarscheduleR{}
}

func (r *djangoCeleryBeatSolarscheduleR) GetSolarDjangoCeleryBeatPeriodictasks() DjangoCeleryBeatPeriodictaskSlice {
	if r == nil {
		return nil
	}
	return r.SolarDjangoCeleryBeatPeriodictasks
}

// djangoCeleryBeatSolarscheduleL is where Load methods for each relationship are stored.
type djangoCeleryBeatSolarscheduleL struct{}

var (
	djangoCeleryBeatSolarscheduleAllColumns            = []string{"id", "event", "latitude", "longitude"}
	djangoCeleryBeatSolarscheduleColumnsWithoutDefault = []string{"event", "latitude", "longitude"}
	djangoCeleryBeatSolarscheduleColumnsWithDefault    = []string{"id"}
	djangoCeleryBeatSolarschedulePrimaryKeyColumns     = []string{"id"}
	djangoCeleryBeatSolarscheduleGeneratedColumns      = []string{}
)

type (
	// DjangoCeleryBeatSolarscheduleSlice is an alias for a slice of pointers to DjangoCeleryBeatSolarschedule.
	// This should almost always be used instead of []DjangoCeleryBeatSolarschedule.
	DjangoCeleryBeatSolarscheduleSlice []*DjangoCeleryBeatSolarschedule

	djangoCeleryBeatSolarscheduleQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	djangoCeleryBeatSolarscheduleType                 = reflect.TypeOf(&DjangoCeleryBeatSolarschedule{})
	djangoCeleryBeatSolarscheduleMapping              = queries.MakeStructMapping(djangoCeleryBeatSolarscheduleType)
	djangoCeleryBeatSolarschedulePrimaryKeyMapping, _ = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, djangoCeleryBeatSolarschedulePrimaryKeyColumns)
	djangoCeleryBeatSolarscheduleInsertCacheMut       sync.RWMutex
	djangoCeleryBeatSolarscheduleInsertCache          = make(map[string]insertCache)
	djangoCeleryBeatSolarscheduleUpdateCacheMut       sync.RWMutex
	djangoCeleryBeatSolarscheduleUpdateCache          = make(map[string]updateCache)
	djangoCeleryBeatSolarscheduleUpsertCacheMut       sync.RWMutex
	djangoCeleryBeatSolarscheduleUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single djangoCeleryBeatSolarschedule record from the query.
func (q djangoCeleryBeatSolarscheduleQuery) One(exec boil.Executor) (*DjangoCeleryBeatSolarschedule, error) {
	o := &DjangoCeleryBeatSolarschedule{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: failed to execute a one query for django_celery_beat_solarschedule")
	}

	return o, nil
}

// All returns all DjangoCeleryBeatSolarschedule records from the query.
func (q djangoCeleryBeatSolarscheduleQuery) All(exec boil.Executor) (DjangoCeleryBeatSolarscheduleSlice, error) {
	var o []*DjangoCeleryBeatSolarschedule

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "schema: failed to assign all query results to DjangoCeleryBeatSolarschedule slice")
	}

	return o, nil
}

// Count returns the count of all DjangoCeleryBeatSolarschedule records in the query.
func (q djangoCeleryBeatSolarscheduleQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to count django_celery_beat_solarschedule rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q djangoCeleryBeatSolarscheduleQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "schema: failed to check if django_celery_beat_solarschedule exists")
	}

	return count > 0, nil
}

// SolarDjangoCeleryBeatPeriodictasks retrieves all the django_celery_beat_periodictask's DjangoCeleryBeatPeriodictasks with an executor via solar_id column.
func (o *DjangoCeleryBeatSolarschedule) SolarDjangoCeleryBeatPeriodictasks(mods ...qm.QueryMod) djangoCeleryBeatPeriodictaskQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"django_celery_beat_periodictask\".\"solar_id\"=?", o.ID),
	)

	return DjangoCeleryBeatPeriodictasks(queryMods...)
}

// LoadSolarDjangoCeleryBeatPeriodictasks allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (djangoCeleryBeatSolarscheduleL) LoadSolarDjangoCeleryBeatPeriodictasks(e boil.Executor, singular bool, maybeDjangoCeleryBeatSolarschedule interface{}, mods queries.Applicator) error {
	var slice []*DjangoCeleryBeatSolarschedule
	var object *DjangoCeleryBeatSolarschedule

	if singular {
		var ok bool
		object, ok = maybeDjangoCeleryBeatSolarschedule.(*DjangoCeleryBeatSolarschedule)
		if !ok {
			object = new(DjangoCeleryBeatSolarschedule)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeDjangoCeleryBeatSolarschedule)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeDjangoCeleryBeatSolarschedule))
			}
		}
	} else {
		s, ok := maybeDjangoCeleryBeatSolarschedule.(*[]*DjangoCeleryBeatSolarschedule)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeDjangoCeleryBeatSolarschedule)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeDjangoCeleryBeatSolarschedule))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &djangoCeleryBeatSolarscheduleR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &djangoCeleryBeatSolarscheduleR{}
			}
			args[obj.ID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`django_celery_beat_periodictask`),
		qm.WhereIn(`django_celery_beat_periodictask.solar_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load django_celery_beat_periodictask")
	}

	var resultSlice []*DjangoCeleryBeatPeriodictask
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice django_celery_beat_periodictask")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on django_celery_beat_periodictask")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for django_celery_beat_periodictask")
	}

	if singular {
		object.R.SolarDjangoCeleryBeatPeriodictasks = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &djangoCeleryBeatPeriodictaskR{}
			}
			foreign.R.Solar = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.SolarID) {
				local.R.SolarDjangoCeleryBeatPeriodictasks = append(local.R.SolarDjangoCeleryBeatPeriodictasks, foreign)
				if foreign.R == nil {
					foreign.R = &djangoCeleryBeatPeriodictaskR{}
				}
				foreign.R.Solar = local
				break
			}
		}
	}

	return nil
}

// AddSolarDjangoCeleryBeatPeriodictasks adds the given related objects to the existing relationships
// of the django_celery_beat_solarschedule, optionally inserting them as new records.
// Appends related to o.R.SolarDjangoCeleryBeatPeriodictasks.
// Sets related.R.Solar appropriately.
func (o *DjangoCeleryBeatSolarschedule) AddSolarDjangoCeleryBeatPeriodictasks(exec boil.Executor, insert bool, related ...*DjangoCeleryBeatPeriodictask) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.SolarID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"django_celery_beat_periodictask\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"solar_id"}),
				strmangle.WhereClause("\"", "\"", 2, djangoCeleryBeatPeriodictaskPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.SolarID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &djangoCeleryBeatSolarscheduleR{
			SolarDjangoCeleryBeatPeriodictasks: related,
		}
	} else {
		o.R.SolarDjangoCeleryBeatPeriodictasks = append(o.R.SolarDjangoCeleryBeatPeriodictasks, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &djangoCeleryBeatPeriodictaskR{
				Solar: o,
			}
		} else {
			rel.R.Solar = o
		}
	}
	return nil
}

// SetSolarDjangoCeleryBeatPeriodictasks removes all previously related items of the
// django_celery_beat_solarschedule replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Solar's SolarDjangoCeleryBeatPeriodictasks accordingly.
// Replaces o.R.SolarDjangoCeleryBeatPeriodictasks with related.
// Sets related.R.Solar's SolarDjangoCeleryBeatPeriodictasks accordingly.
func (o *DjangoCeleryBeatSolarschedule) SetSolarDjangoCeleryBeatPeriodictasks(exec boil.Executor, insert bool, related ...*DjangoCeleryBeatPeriodictask) error {
	query := "update \"django_celery_beat_periodictask\" set \"solar_id\" = null where \"solar_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.SolarDjangoCeleryBeatPeriodictasks {
			queries.SetScanner(&rel.SolarID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Solar = nil
		}
		o.R.SolarDjangoCeleryBeatPeriodictasks = nil
	}

	return o.AddSolarDjangoCeleryBeatPeriodictasks(exec, insert, related...)
}

// RemoveSolarDjangoCeleryBeatPeriodictasks relationships from objects passed in.
// Removes related items from R.SolarDjangoCeleryBeatPeriodictasks (uses pointer comparison, removal does not keep order)
// Sets related.R.Solar.
func (o *DjangoCeleryBeatSolarschedule) RemoveSolarDjangoCeleryBeatPeriodictasks(exec boil.Executor, related ...*DjangoCeleryBeatPeriodictask) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.SolarID, nil)
		if rel.R != nil {
			rel.R.Solar = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("solar_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.SolarDjangoCeleryBeatPeriodictasks {
			if rel != ri {
				continue
			}

			ln := len(o.R.SolarDjangoCeleryBeatPeriodictasks)
			if ln > 1 && i < ln-1 {
				o.R.SolarDjangoCeleryBeatPeriodictasks[i] = o.R.SolarDjangoCeleryBeatPeriodictasks[ln-1]
			}
			o.R.SolarDjangoCeleryBeatPeriodictasks = o.R.SolarDjangoCeleryBeatPeriodictasks[:ln-1]
			break
		}
	}

	return nil
}

// DjangoCeleryBeatSolarschedules retrieves all the records using an executor.
func DjangoCeleryBeatSolarschedules(mods ...qm.QueryMod) djangoCeleryBeatSolarscheduleQuery {
	mods = append(mods, qm.From("\"django_celery_beat_solarschedule\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"django_celery_beat_solarschedule\".*"})
	}

	return djangoCeleryBeatSolarscheduleQuery{q}
}

// FindDjangoCeleryBeatSolarschedule retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDjangoCeleryBeatSolarschedule(exec boil.Executor, iD int, selectCols ...string) (*DjangoCeleryBeatSolarschedule, error) {
	djangoCeleryBeatSolarscheduleObj := &DjangoCeleryBeatSolarschedule{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"django_celery_beat_solarschedule\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, djangoCeleryBeatSolarscheduleObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: unable to select from django_celery_beat_solarschedule")
	}

	return djangoCeleryBeatSolarscheduleObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *DjangoCeleryBeatSolarschedule) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("schema: no django_celery_beat_solarschedule provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(djangoCeleryBeatSolarscheduleColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	djangoCeleryBeatSolarscheduleInsertCacheMut.RLock()
	cache, cached := djangoCeleryBeatSolarscheduleInsertCache[key]
	djangoCeleryBeatSolarscheduleInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			djangoCeleryBeatSolarscheduleAllColumns,
			djangoCeleryBeatSolarscheduleColumnsWithDefault,
			djangoCeleryBeatSolarscheduleColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"django_celery_beat_solarschedule\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"django_celery_beat_solarschedule\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "schema: unable to insert into django_celery_beat_solarschedule")
	}

	if !cached {
		djangoCeleryBeatSolarscheduleInsertCacheMut.Lock()
		djangoCeleryBeatSolarscheduleInsertCache[key] = cache
		djangoCeleryBeatSolarscheduleInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the DjangoCeleryBeatSolarschedule.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *DjangoCeleryBeatSolarschedule) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	djangoCeleryBeatSolarscheduleUpdateCacheMut.RLock()
	cache, cached := djangoCeleryBeatSolarscheduleUpdateCache[key]
	djangoCeleryBeatSolarscheduleUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			djangoCeleryBeatSolarscheduleAllColumns,
			djangoCeleryBeatSolarschedulePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("schema: unable to update django_celery_beat_solarschedule, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"django_celery_beat_solarschedule\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, djangoCeleryBeatSolarschedulePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, append(wl, djangoCeleryBeatSolarschedulePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update django_celery_beat_solarschedule row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by update for django_celery_beat_solarschedule")
	}

	if !cached {
		djangoCeleryBeatSolarscheduleUpdateCacheMut.Lock()
		djangoCeleryBeatSolarscheduleUpdateCache[key] = cache
		djangoCeleryBeatSolarscheduleUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q djangoCeleryBeatSolarscheduleQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all for django_celery_beat_solarschedule")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected for django_celery_beat_solarschedule")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DjangoCeleryBeatSolarscheduleSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("schema: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), djangoCeleryBeatSolarschedulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"django_celery_beat_solarschedule\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, djangoCeleryBeatSolarschedulePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all in djangoCeleryBeatSolarschedule slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected all in update all djangoCeleryBeatSolarschedule")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *DjangoCeleryBeatSolarschedule) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("schema: no django_celery_beat_solarschedule provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(djangoCeleryBeatSolarscheduleColumnsWithDefault, o)

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

	djangoCeleryBeatSolarscheduleUpsertCacheMut.RLock()
	cache, cached := djangoCeleryBeatSolarscheduleUpsertCache[key]
	djangoCeleryBeatSolarscheduleUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			djangoCeleryBeatSolarscheduleAllColumns,
			djangoCeleryBeatSolarscheduleColumnsWithDefault,
			djangoCeleryBeatSolarscheduleColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			djangoCeleryBeatSolarscheduleAllColumns,
			djangoCeleryBeatSolarschedulePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("schema: unable to upsert django_celery_beat_solarschedule, could not build update column list")
		}

		ret := strmangle.SetComplement(djangoCeleryBeatSolarscheduleAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(djangoCeleryBeatSolarschedulePrimaryKeyColumns) == 0 {
				return errors.New("schema: unable to upsert django_celery_beat_solarschedule, could not build conflict column list")
			}

			conflict = make([]string, len(djangoCeleryBeatSolarschedulePrimaryKeyColumns))
			copy(conflict, djangoCeleryBeatSolarschedulePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"django_celery_beat_solarschedule\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(djangoCeleryBeatSolarscheduleType, djangoCeleryBeatSolarscheduleMapping, ret)
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "schema: unable to upsert django_celery_beat_solarschedule")
	}

	if !cached {
		djangoCeleryBeatSolarscheduleUpsertCacheMut.Lock()
		djangoCeleryBeatSolarscheduleUpsertCache[key] = cache
		djangoCeleryBeatSolarscheduleUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single DjangoCeleryBeatSolarschedule record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *DjangoCeleryBeatSolarschedule) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("schema: no DjangoCeleryBeatSolarschedule provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), djangoCeleryBeatSolarschedulePrimaryKeyMapping)
	sql := "DELETE FROM \"django_celery_beat_solarschedule\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete from django_celery_beat_solarschedule")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by delete for django_celery_beat_solarschedule")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q djangoCeleryBeatSolarscheduleQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("schema: no djangoCeleryBeatSolarscheduleQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from django_celery_beat_solarschedule")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for django_celery_beat_solarschedule")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DjangoCeleryBeatSolarscheduleSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), djangoCeleryBeatSolarschedulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"django_celery_beat_solarschedule\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, djangoCeleryBeatSolarschedulePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from djangoCeleryBeatSolarschedule slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for django_celery_beat_solarschedule")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *DjangoCeleryBeatSolarschedule) Reload(exec boil.Executor) error {
	ret, err := FindDjangoCeleryBeatSolarschedule(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DjangoCeleryBeatSolarscheduleSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := DjangoCeleryBeatSolarscheduleSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), djangoCeleryBeatSolarschedulePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"django_celery_beat_solarschedule\".* FROM \"django_celery_beat_solarschedule\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, djangoCeleryBeatSolarschedulePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "schema: unable to reload all in DjangoCeleryBeatSolarscheduleSlice")
	}

	*o = slice

	return nil
}

// DjangoCeleryBeatSolarscheduleExists checks if the DjangoCeleryBeatSolarschedule row exists.
func DjangoCeleryBeatSolarscheduleExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"django_celery_beat_solarschedule\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "schema: unable to check if django_celery_beat_solarschedule exists")
	}

	return exists, nil
}

// Exists checks if the DjangoCeleryBeatSolarschedule row exists.
func (o *DjangoCeleryBeatSolarschedule) Exists(exec boil.Executor) (bool, error) {
	return DjangoCeleryBeatSolarscheduleExists(exec, o.ID)
}
