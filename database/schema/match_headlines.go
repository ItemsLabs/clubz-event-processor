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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// MatchHeadline is an object representing the database table.
type MatchHeadline struct {
	ID          int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title       string      `boil:"title" json:"title" toml:"title" yaml:"title"`
	Description string      `boil:"description" json:"description" toml:"description" yaml:"description"`
	Images      string      `boil:"images" json:"images" toml:"images" yaml:"images"`
	Type        string      `boil:"type" json:"type" toml:"type" yaml:"type"`
	ScreenType  int         `boil:"screen_type" json:"screen_type" toml:"screen_type" yaml:"screen_type"`
	MatchID     string      `boil:"match_id" json:"match_id" toml:"match_id" yaml:"match_id"`
	ImageType   null.String `boil:"image_type" json:"image_type,omitempty" toml:"image_type" yaml:"image_type,omitempty"`

	R *matchHeadlineR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L matchHeadlineL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MatchHeadlineColumns = struct {
	ID          string
	Title       string
	Description string
	Images      string
	Type        string
	ScreenType  string
	MatchID     string
	ImageType   string
}{
	ID:          "id",
	Title:       "title",
	Description: "description",
	Images:      "images",
	Type:        "type",
	ScreenType:  "screen_type",
	MatchID:     "match_id",
	ImageType:   "image_type",
}

var MatchHeadlineTableColumns = struct {
	ID          string
	Title       string
	Description string
	Images      string
	Type        string
	ScreenType  string
	MatchID     string
	ImageType   string
}{
	ID:          "match_headlines.id",
	Title:       "match_headlines.title",
	Description: "match_headlines.description",
	Images:      "match_headlines.images",
	Type:        "match_headlines.type",
	ScreenType:  "match_headlines.screen_type",
	MatchID:     "match_headlines.match_id",
	ImageType:   "match_headlines.image_type",
}

// Generated where

var MatchHeadlineWhere = struct {
	ID          whereHelperint
	Title       whereHelperstring
	Description whereHelperstring
	Images      whereHelperstring
	Type        whereHelperstring
	ScreenType  whereHelperint
	MatchID     whereHelperstring
	ImageType   whereHelpernull_String
}{
	ID:          whereHelperint{field: "\"match_headlines\".\"id\""},
	Title:       whereHelperstring{field: "\"match_headlines\".\"title\""},
	Description: whereHelperstring{field: "\"match_headlines\".\"description\""},
	Images:      whereHelperstring{field: "\"match_headlines\".\"images\""},
	Type:        whereHelperstring{field: "\"match_headlines\".\"type\""},
	ScreenType:  whereHelperint{field: "\"match_headlines\".\"screen_type\""},
	MatchID:     whereHelperstring{field: "\"match_headlines\".\"match_id\""},
	ImageType:   whereHelpernull_String{field: "\"match_headlines\".\"image_type\""},
}

// MatchHeadlineRels is where relationship names are stored.
var MatchHeadlineRels = struct {
	Match string
}{
	Match: "Match",
}

// matchHeadlineR is where relationships are stored.
type matchHeadlineR struct {
	Match *Match `boil:"Match" json:"Match" toml:"Match" yaml:"Match"`
}

// NewStruct creates a new relationship struct
func (*matchHeadlineR) NewStruct() *matchHeadlineR {
	return &matchHeadlineR{}
}

func (r *matchHeadlineR) GetMatch() *Match {
	if r == nil {
		return nil
	}
	return r.Match
}

// matchHeadlineL is where Load methods for each relationship are stored.
type matchHeadlineL struct{}

var (
	matchHeadlineAllColumns            = []string{"id", "title", "description", "images", "type", "screen_type", "match_id", "image_type"}
	matchHeadlineColumnsWithoutDefault = []string{"title", "description", "images", "type", "screen_type", "match_id"}
	matchHeadlineColumnsWithDefault    = []string{"id", "image_type"}
	matchHeadlinePrimaryKeyColumns     = []string{"id"}
	matchHeadlineGeneratedColumns      = []string{}
)

type (
	// MatchHeadlineSlice is an alias for a slice of pointers to MatchHeadline.
	// This should almost always be used instead of []MatchHeadline.
	MatchHeadlineSlice []*MatchHeadline

	matchHeadlineQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	matchHeadlineType                 = reflect.TypeOf(&MatchHeadline{})
	matchHeadlineMapping              = queries.MakeStructMapping(matchHeadlineType)
	matchHeadlinePrimaryKeyMapping, _ = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, matchHeadlinePrimaryKeyColumns)
	matchHeadlineInsertCacheMut       sync.RWMutex
	matchHeadlineInsertCache          = make(map[string]insertCache)
	matchHeadlineUpdateCacheMut       sync.RWMutex
	matchHeadlineUpdateCache          = make(map[string]updateCache)
	matchHeadlineUpsertCacheMut       sync.RWMutex
	matchHeadlineUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single matchHeadline record from the query.
func (q matchHeadlineQuery) One(exec boil.Executor) (*MatchHeadline, error) {
	o := &MatchHeadline{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: failed to execute a one query for match_headlines")
	}

	return o, nil
}

// All returns all MatchHeadline records from the query.
func (q matchHeadlineQuery) All(exec boil.Executor) (MatchHeadlineSlice, error) {
	var o []*MatchHeadline

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "schema: failed to assign all query results to MatchHeadline slice")
	}

	return o, nil
}

// Count returns the count of all MatchHeadline records in the query.
func (q matchHeadlineQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to count match_headlines rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q matchHeadlineQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "schema: failed to check if match_headlines exists")
	}

	return count > 0, nil
}

// Match pointed to by the foreign key.
func (o *MatchHeadline) Match(mods ...qm.QueryMod) matchQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.MatchID),
	}

	queryMods = append(queryMods, mods...)

	return Matches(queryMods...)
}

// LoadMatch allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (matchHeadlineL) LoadMatch(e boil.Executor, singular bool, maybeMatchHeadline interface{}, mods queries.Applicator) error {
	var slice []*MatchHeadline
	var object *MatchHeadline

	if singular {
		var ok bool
		object, ok = maybeMatchHeadline.(*MatchHeadline)
		if !ok {
			object = new(MatchHeadline)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeMatchHeadline)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeMatchHeadline))
			}
		}
	} else {
		s, ok := maybeMatchHeadline.(*[]*MatchHeadline)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeMatchHeadline)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeMatchHeadline))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &matchHeadlineR{}
		}
		args[object.MatchID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &matchHeadlineR{}
			}

			args[obj.MatchID] = struct{}{}

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
		qm.From(`matches`),
		qm.WhereIn(`matches.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Match")
	}

	var resultSlice []*Match
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Match")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for matches")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for matches")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Match = foreign
		if foreign.R == nil {
			foreign.R = &matchR{}
		}
		foreign.R.MatchHeadlines = append(foreign.R.MatchHeadlines, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.MatchID == foreign.ID {
				local.R.Match = foreign
				if foreign.R == nil {
					foreign.R = &matchR{}
				}
				foreign.R.MatchHeadlines = append(foreign.R.MatchHeadlines, local)
				break
			}
		}
	}

	return nil
}

// SetMatch of the matchHeadline to the related item.
// Sets o.R.Match to related.
// Adds o to related.R.MatchHeadlines.
func (o *MatchHeadline) SetMatch(exec boil.Executor, insert bool, related *Match) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"match_headlines\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"match_id"}),
		strmangle.WhereClause("\"", "\"", 2, matchHeadlinePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.MatchID = related.ID
	if o.R == nil {
		o.R = &matchHeadlineR{
			Match: related,
		}
	} else {
		o.R.Match = related
	}

	if related.R == nil {
		related.R = &matchR{
			MatchHeadlines: MatchHeadlineSlice{o},
		}
	} else {
		related.R.MatchHeadlines = append(related.R.MatchHeadlines, o)
	}

	return nil
}

// MatchHeadlines retrieves all the records using an executor.
func MatchHeadlines(mods ...qm.QueryMod) matchHeadlineQuery {
	mods = append(mods, qm.From("\"match_headlines\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"match_headlines\".*"})
	}

	return matchHeadlineQuery{q}
}

// FindMatchHeadline retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMatchHeadline(exec boil.Executor, iD int, selectCols ...string) (*MatchHeadline, error) {
	matchHeadlineObj := &MatchHeadline{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"match_headlines\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, matchHeadlineObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: unable to select from match_headlines")
	}

	return matchHeadlineObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MatchHeadline) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("schema: no match_headlines provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(matchHeadlineColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	matchHeadlineInsertCacheMut.RLock()
	cache, cached := matchHeadlineInsertCache[key]
	matchHeadlineInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			matchHeadlineAllColumns,
			matchHeadlineColumnsWithDefault,
			matchHeadlineColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"match_headlines\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"match_headlines\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "schema: unable to insert into match_headlines")
	}

	if !cached {
		matchHeadlineInsertCacheMut.Lock()
		matchHeadlineInsertCache[key] = cache
		matchHeadlineInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the MatchHeadline.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MatchHeadline) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	matchHeadlineUpdateCacheMut.RLock()
	cache, cached := matchHeadlineUpdateCache[key]
	matchHeadlineUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			matchHeadlineAllColumns,
			matchHeadlinePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("schema: unable to update match_headlines, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"match_headlines\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, matchHeadlinePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, append(wl, matchHeadlinePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "schema: unable to update match_headlines row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by update for match_headlines")
	}

	if !cached {
		matchHeadlineUpdateCacheMut.Lock()
		matchHeadlineUpdateCache[key] = cache
		matchHeadlineUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q matchHeadlineQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all for match_headlines")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected for match_headlines")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MatchHeadlineSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), matchHeadlinePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"match_headlines\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, matchHeadlinePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all in matchHeadline slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected all in update all matchHeadline")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MatchHeadline) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("schema: no match_headlines provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(matchHeadlineColumnsWithDefault, o)

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

	matchHeadlineUpsertCacheMut.RLock()
	cache, cached := matchHeadlineUpsertCache[key]
	matchHeadlineUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			matchHeadlineAllColumns,
			matchHeadlineColumnsWithDefault,
			matchHeadlineColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			matchHeadlineAllColumns,
			matchHeadlinePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("schema: unable to upsert match_headlines, could not build update column list")
		}

		ret := strmangle.SetComplement(matchHeadlineAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(matchHeadlinePrimaryKeyColumns) == 0 {
				return errors.New("schema: unable to upsert match_headlines, could not build conflict column list")
			}

			conflict = make([]string, len(matchHeadlinePrimaryKeyColumns))
			copy(conflict, matchHeadlinePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"match_headlines\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(matchHeadlineType, matchHeadlineMapping, ret)
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
		return errors.Wrap(err, "schema: unable to upsert match_headlines")
	}

	if !cached {
		matchHeadlineUpsertCacheMut.Lock()
		matchHeadlineUpsertCache[key] = cache
		matchHeadlineUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single MatchHeadline record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MatchHeadline) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("schema: no MatchHeadline provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), matchHeadlinePrimaryKeyMapping)
	sql := "DELETE FROM \"match_headlines\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete from match_headlines")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by delete for match_headlines")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q matchHeadlineQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("schema: no matchHeadlineQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from match_headlines")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for match_headlines")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MatchHeadlineSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), matchHeadlinePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"match_headlines\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, matchHeadlinePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from matchHeadline slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for match_headlines")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *MatchHeadline) Reload(exec boil.Executor) error {
	ret, err := FindMatchHeadline(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MatchHeadlineSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MatchHeadlineSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), matchHeadlinePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"match_headlines\".* FROM \"match_headlines\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, matchHeadlinePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "schema: unable to reload all in MatchHeadlineSlice")
	}

	*o = slice

	return nil
}

// MatchHeadlineExists checks if the MatchHeadline row exists.
func MatchHeadlineExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"match_headlines\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "schema: unable to check if match_headlines exists")
	}

	return exists, nil
}

// Exists checks if the MatchHeadline row exists.
func (o *MatchHeadline) Exists(exec boil.Executor) (bool, error) {
	return MatchHeadlineExists(exec, o.ID)
}
