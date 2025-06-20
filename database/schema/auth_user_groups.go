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
	"github.com/volatiletech/strmangle"
)

// AuthUserGroup is an object representing the database table.
type AuthUserGroup struct {
	ID      int `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID  int `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	GroupID int `boil:"group_id" json:"group_id" toml:"group_id" yaml:"group_id"`

	R *authUserGroupR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L authUserGroupL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AuthUserGroupColumns = struct {
	ID      string
	UserID  string
	GroupID string
}{
	ID:      "id",
	UserID:  "user_id",
	GroupID: "group_id",
}

var AuthUserGroupTableColumns = struct {
	ID      string
	UserID  string
	GroupID string
}{
	ID:      "auth_user_groups.id",
	UserID:  "auth_user_groups.user_id",
	GroupID: "auth_user_groups.group_id",
}

// Generated where

var AuthUserGroupWhere = struct {
	ID      whereHelperint
	UserID  whereHelperint
	GroupID whereHelperint
}{
	ID:      whereHelperint{field: "\"auth_user_groups\".\"id\""},
	UserID:  whereHelperint{field: "\"auth_user_groups\".\"user_id\""},
	GroupID: whereHelperint{field: "\"auth_user_groups\".\"group_id\""},
}

// AuthUserGroupRels is where relationship names are stored.
var AuthUserGroupRels = struct {
	Group string
	User  string
}{
	Group: "Group",
	User:  "User",
}

// authUserGroupR is where relationships are stored.
type authUserGroupR struct {
	Group *AuthGroup `boil:"Group" json:"Group" toml:"Group" yaml:"Group"`
	User  *AuthUser  `boil:"User" json:"User" toml:"User" yaml:"User"`
}

// NewStruct creates a new relationship struct
func (*authUserGroupR) NewStruct() *authUserGroupR {
	return &authUserGroupR{}
}

func (r *authUserGroupR) GetGroup() *AuthGroup {
	if r == nil {
		return nil
	}
	return r.Group
}

func (r *authUserGroupR) GetUser() *AuthUser {
	if r == nil {
		return nil
	}
	return r.User
}

// authUserGroupL is where Load methods for each relationship are stored.
type authUserGroupL struct{}

var (
	authUserGroupAllColumns            = []string{"id", "user_id", "group_id"}
	authUserGroupColumnsWithoutDefault = []string{"user_id", "group_id"}
	authUserGroupColumnsWithDefault    = []string{"id"}
	authUserGroupPrimaryKeyColumns     = []string{"id"}
	authUserGroupGeneratedColumns      = []string{}
)

type (
	// AuthUserGroupSlice is an alias for a slice of pointers to AuthUserGroup.
	// This should almost always be used instead of []AuthUserGroup.
	AuthUserGroupSlice []*AuthUserGroup

	authUserGroupQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	authUserGroupType                 = reflect.TypeOf(&AuthUserGroup{})
	authUserGroupMapping              = queries.MakeStructMapping(authUserGroupType)
	authUserGroupPrimaryKeyMapping, _ = queries.BindMapping(authUserGroupType, authUserGroupMapping, authUserGroupPrimaryKeyColumns)
	authUserGroupInsertCacheMut       sync.RWMutex
	authUserGroupInsertCache          = make(map[string]insertCache)
	authUserGroupUpdateCacheMut       sync.RWMutex
	authUserGroupUpdateCache          = make(map[string]updateCache)
	authUserGroupUpsertCacheMut       sync.RWMutex
	authUserGroupUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single authUserGroup record from the query.
func (q authUserGroupQuery) One(exec boil.Executor) (*AuthUserGroup, error) {
	o := &AuthUserGroup{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: failed to execute a one query for auth_user_groups")
	}

	return o, nil
}

// All returns all AuthUserGroup records from the query.
func (q authUserGroupQuery) All(exec boil.Executor) (AuthUserGroupSlice, error) {
	var o []*AuthUserGroup

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "schema: failed to assign all query results to AuthUserGroup slice")
	}

	return o, nil
}

// Count returns the count of all AuthUserGroup records in the query.
func (q authUserGroupQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to count auth_user_groups rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q authUserGroupQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "schema: failed to check if auth_user_groups exists")
	}

	return count > 0, nil
}

// Group pointed to by the foreign key.
func (o *AuthUserGroup) Group(mods ...qm.QueryMod) authGroupQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.GroupID),
	}

	queryMods = append(queryMods, mods...)

	return AuthGroups(queryMods...)
}

// User pointed to by the foreign key.
func (o *AuthUserGroup) User(mods ...qm.QueryMod) authUserQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	return AuthUsers(queryMods...)
}

// LoadGroup allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (authUserGroupL) LoadGroup(e boil.Executor, singular bool, maybeAuthUserGroup interface{}, mods queries.Applicator) error {
	var slice []*AuthUserGroup
	var object *AuthUserGroup

	if singular {
		var ok bool
		object, ok = maybeAuthUserGroup.(*AuthUserGroup)
		if !ok {
			object = new(AuthUserGroup)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthUserGroup)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthUserGroup))
			}
		}
	} else {
		s, ok := maybeAuthUserGroup.(*[]*AuthUserGroup)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthUserGroup)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthUserGroup))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authUserGroupR{}
		}
		args[object.GroupID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authUserGroupR{}
			}

			args[obj.GroupID] = struct{}{}

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
		qm.From(`auth_group`),
		qm.WhereIn(`auth_group.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AuthGroup")
	}

	var resultSlice []*AuthGroup
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AuthGroup")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for auth_group")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for auth_group")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Group = foreign
		if foreign.R == nil {
			foreign.R = &authGroupR{}
		}
		foreign.R.GroupAuthUserGroups = append(foreign.R.GroupAuthUserGroups, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.GroupID == foreign.ID {
				local.R.Group = foreign
				if foreign.R == nil {
					foreign.R = &authGroupR{}
				}
				foreign.R.GroupAuthUserGroups = append(foreign.R.GroupAuthUserGroups, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (authUserGroupL) LoadUser(e boil.Executor, singular bool, maybeAuthUserGroup interface{}, mods queries.Applicator) error {
	var slice []*AuthUserGroup
	var object *AuthUserGroup

	if singular {
		var ok bool
		object, ok = maybeAuthUserGroup.(*AuthUserGroup)
		if !ok {
			object = new(AuthUserGroup)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthUserGroup)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthUserGroup))
			}
		}
	} else {
		s, ok := maybeAuthUserGroup.(*[]*AuthUserGroup)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthUserGroup)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthUserGroup))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authUserGroupR{}
		}
		args[object.UserID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authUserGroupR{}
			}

			args[obj.UserID] = struct{}{}

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
		qm.From(`auth_user`),
		qm.WhereIn(`auth_user.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AuthUser")
	}

	var resultSlice []*AuthUser
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AuthUser")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for auth_user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for auth_user")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &authUserR{}
		}
		foreign.R.UserAuthUserGroups = append(foreign.R.UserAuthUserGroups, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &authUserR{}
				}
				foreign.R.UserAuthUserGroups = append(foreign.R.UserAuthUserGroups, local)
				break
			}
		}
	}

	return nil
}

// SetGroup of the authUserGroup to the related item.
// Sets o.R.Group to related.
// Adds o to related.R.GroupAuthUserGroups.
func (o *AuthUserGroup) SetGroup(exec boil.Executor, insert bool, related *AuthGroup) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"auth_user_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"group_id"}),
		strmangle.WhereClause("\"", "\"", 2, authUserGroupPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GroupID = related.ID
	if o.R == nil {
		o.R = &authUserGroupR{
			Group: related,
		}
	} else {
		o.R.Group = related
	}

	if related.R == nil {
		related.R = &authGroupR{
			GroupAuthUserGroups: AuthUserGroupSlice{o},
		}
	} else {
		related.R.GroupAuthUserGroups = append(related.R.GroupAuthUserGroups, o)
	}

	return nil
}

// SetUser of the authUserGroup to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserAuthUserGroups.
func (o *AuthUserGroup) SetUser(exec boil.Executor, insert bool, related *AuthUser) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"auth_user_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, authUserGroupPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &authUserGroupR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &authUserR{
			UserAuthUserGroups: AuthUserGroupSlice{o},
		}
	} else {
		related.R.UserAuthUserGroups = append(related.R.UserAuthUserGroups, o)
	}

	return nil
}

// AuthUserGroups retrieves all the records using an executor.
func AuthUserGroups(mods ...qm.QueryMod) authUserGroupQuery {
	mods = append(mods, qm.From("\"auth_user_groups\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"auth_user_groups\".*"})
	}

	return authUserGroupQuery{q}
}

// FindAuthUserGroup retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAuthUserGroup(exec boil.Executor, iD int, selectCols ...string) (*AuthUserGroup, error) {
	authUserGroupObj := &AuthUserGroup{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"auth_user_groups\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, authUserGroupObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: unable to select from auth_user_groups")
	}

	return authUserGroupObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AuthUserGroup) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("schema: no auth_user_groups provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(authUserGroupColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	authUserGroupInsertCacheMut.RLock()
	cache, cached := authUserGroupInsertCache[key]
	authUserGroupInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			authUserGroupAllColumns,
			authUserGroupColumnsWithDefault,
			authUserGroupColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(authUserGroupType, authUserGroupMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(authUserGroupType, authUserGroupMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"auth_user_groups\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"auth_user_groups\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "schema: unable to insert into auth_user_groups")
	}

	if !cached {
		authUserGroupInsertCacheMut.Lock()
		authUserGroupInsertCache[key] = cache
		authUserGroupInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the AuthUserGroup.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AuthUserGroup) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	authUserGroupUpdateCacheMut.RLock()
	cache, cached := authUserGroupUpdateCache[key]
	authUserGroupUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			authUserGroupAllColumns,
			authUserGroupPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("schema: unable to update auth_user_groups, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"auth_user_groups\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, authUserGroupPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(authUserGroupType, authUserGroupMapping, append(wl, authUserGroupPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "schema: unable to update auth_user_groups row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by update for auth_user_groups")
	}

	if !cached {
		authUserGroupUpdateCacheMut.Lock()
		authUserGroupUpdateCache[key] = cache
		authUserGroupUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q authUserGroupQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all for auth_user_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected for auth_user_groups")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AuthUserGroupSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"auth_user_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, authUserGroupPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all in authUserGroup slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected all in update all authUserGroup")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AuthUserGroup) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("schema: no auth_user_groups provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(authUserGroupColumnsWithDefault, o)

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

	authUserGroupUpsertCacheMut.RLock()
	cache, cached := authUserGroupUpsertCache[key]
	authUserGroupUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			authUserGroupAllColumns,
			authUserGroupColumnsWithDefault,
			authUserGroupColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			authUserGroupAllColumns,
			authUserGroupPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("schema: unable to upsert auth_user_groups, could not build update column list")
		}

		ret := strmangle.SetComplement(authUserGroupAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(authUserGroupPrimaryKeyColumns) == 0 {
				return errors.New("schema: unable to upsert auth_user_groups, could not build conflict column list")
			}

			conflict = make([]string, len(authUserGroupPrimaryKeyColumns))
			copy(conflict, authUserGroupPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"auth_user_groups\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(authUserGroupType, authUserGroupMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(authUserGroupType, authUserGroupMapping, ret)
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
		return errors.Wrap(err, "schema: unable to upsert auth_user_groups")
	}

	if !cached {
		authUserGroupUpsertCacheMut.Lock()
		authUserGroupUpsertCache[key] = cache
		authUserGroupUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single AuthUserGroup record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AuthUserGroup) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("schema: no AuthUserGroup provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), authUserGroupPrimaryKeyMapping)
	sql := "DELETE FROM \"auth_user_groups\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete from auth_user_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by delete for auth_user_groups")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q authUserGroupQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("schema: no authUserGroupQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from auth_user_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for auth_user_groups")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AuthUserGroupSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"auth_user_groups\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, authUserGroupPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from authUserGroup slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for auth_user_groups")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AuthUserGroup) Reload(exec boil.Executor) error {
	ret, err := FindAuthUserGroup(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthUserGroupSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AuthUserGroupSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"auth_user_groups\".* FROM \"auth_user_groups\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, authUserGroupPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "schema: unable to reload all in AuthUserGroupSlice")
	}

	*o = slice

	return nil
}

// AuthUserGroupExists checks if the AuthUserGroup row exists.
func AuthUserGroupExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"auth_user_groups\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "schema: unable to check if auth_user_groups exists")
	}

	return exists, nil
}

// Exists checks if the AuthUserGroup row exists.
func (o *AuthUserGroup) Exists(exec boil.Executor) (bool, error) {
	return AuthUserGroupExists(exec, o.ID)
}
