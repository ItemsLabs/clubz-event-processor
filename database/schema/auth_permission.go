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

// AuthPermission is an object representing the database table.
type AuthPermission struct {
	ID            int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name          string `boil:"name" json:"name" toml:"name" yaml:"name"`
	ContentTypeID int    `boil:"content_type_id" json:"content_type_id" toml:"content_type_id" yaml:"content_type_id"`
	Codename      string `boil:"codename" json:"codename" toml:"codename" yaml:"codename"`

	R *authPermissionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L authPermissionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AuthPermissionColumns = struct {
	ID            string
	Name          string
	ContentTypeID string
	Codename      string
}{
	ID:            "id",
	Name:          "name",
	ContentTypeID: "content_type_id",
	Codename:      "codename",
}

var AuthPermissionTableColumns = struct {
	ID            string
	Name          string
	ContentTypeID string
	Codename      string
}{
	ID:            "auth_permission.id",
	Name:          "auth_permission.name",
	ContentTypeID: "auth_permission.content_type_id",
	Codename:      "auth_permission.codename",
}

// Generated where

var AuthPermissionWhere = struct {
	ID            whereHelperint
	Name          whereHelperstring
	ContentTypeID whereHelperint
	Codename      whereHelperstring
}{
	ID:            whereHelperint{field: "\"auth_permission\".\"id\""},
	Name:          whereHelperstring{field: "\"auth_permission\".\"name\""},
	ContentTypeID: whereHelperint{field: "\"auth_permission\".\"content_type_id\""},
	Codename:      whereHelperstring{field: "\"auth_permission\".\"codename\""},
}

// AuthPermissionRels is where relationship names are stored.
var AuthPermissionRels = struct {
	ContentType                       string
	PermissionAuthGroupPermissions    string
	PermissionAuthUserUserPermissions string
}{
	ContentType:                       "ContentType",
	PermissionAuthGroupPermissions:    "PermissionAuthGroupPermissions",
	PermissionAuthUserUserPermissions: "PermissionAuthUserUserPermissions",
}

// authPermissionR is where relationships are stored.
type authPermissionR struct {
	ContentType                       *DjangoContentType          `boil:"ContentType" json:"ContentType" toml:"ContentType" yaml:"ContentType"`
	PermissionAuthGroupPermissions    AuthGroupPermissionSlice    `boil:"PermissionAuthGroupPermissions" json:"PermissionAuthGroupPermissions" toml:"PermissionAuthGroupPermissions" yaml:"PermissionAuthGroupPermissions"`
	PermissionAuthUserUserPermissions AuthUserUserPermissionSlice `boil:"PermissionAuthUserUserPermissions" json:"PermissionAuthUserUserPermissions" toml:"PermissionAuthUserUserPermissions" yaml:"PermissionAuthUserUserPermissions"`
}

// NewStruct creates a new relationship struct
func (*authPermissionR) NewStruct() *authPermissionR {
	return &authPermissionR{}
}

func (r *authPermissionR) GetContentType() *DjangoContentType {
	if r == nil {
		return nil
	}
	return r.ContentType
}

func (r *authPermissionR) GetPermissionAuthGroupPermissions() AuthGroupPermissionSlice {
	if r == nil {
		return nil
	}
	return r.PermissionAuthGroupPermissions
}

func (r *authPermissionR) GetPermissionAuthUserUserPermissions() AuthUserUserPermissionSlice {
	if r == nil {
		return nil
	}
	return r.PermissionAuthUserUserPermissions
}

// authPermissionL is where Load methods for each relationship are stored.
type authPermissionL struct{}

var (
	authPermissionAllColumns            = []string{"id", "name", "content_type_id", "codename"}
	authPermissionColumnsWithoutDefault = []string{"name", "content_type_id", "codename"}
	authPermissionColumnsWithDefault    = []string{"id"}
	authPermissionPrimaryKeyColumns     = []string{"id"}
	authPermissionGeneratedColumns      = []string{}
)

type (
	// AuthPermissionSlice is an alias for a slice of pointers to AuthPermission.
	// This should almost always be used instead of []AuthPermission.
	AuthPermissionSlice []*AuthPermission

	authPermissionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	authPermissionType                 = reflect.TypeOf(&AuthPermission{})
	authPermissionMapping              = queries.MakeStructMapping(authPermissionType)
	authPermissionPrimaryKeyMapping, _ = queries.BindMapping(authPermissionType, authPermissionMapping, authPermissionPrimaryKeyColumns)
	authPermissionInsertCacheMut       sync.RWMutex
	authPermissionInsertCache          = make(map[string]insertCache)
	authPermissionUpdateCacheMut       sync.RWMutex
	authPermissionUpdateCache          = make(map[string]updateCache)
	authPermissionUpsertCacheMut       sync.RWMutex
	authPermissionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single authPermission record from the query.
func (q authPermissionQuery) One(exec boil.Executor) (*AuthPermission, error) {
	o := &AuthPermission{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: failed to execute a one query for auth_permission")
	}

	return o, nil
}

// All returns all AuthPermission records from the query.
func (q authPermissionQuery) All(exec boil.Executor) (AuthPermissionSlice, error) {
	var o []*AuthPermission

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "schema: failed to assign all query results to AuthPermission slice")
	}

	return o, nil
}

// Count returns the count of all AuthPermission records in the query.
func (q authPermissionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to count auth_permission rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q authPermissionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "schema: failed to check if auth_permission exists")
	}

	return count > 0, nil
}

// ContentType pointed to by the foreign key.
func (o *AuthPermission) ContentType(mods ...qm.QueryMod) djangoContentTypeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ContentTypeID),
	}

	queryMods = append(queryMods, mods...)

	return DjangoContentTypes(queryMods...)
}

// PermissionAuthGroupPermissions retrieves all the auth_group_permission's AuthGroupPermissions with an executor via permission_id column.
func (o *AuthPermission) PermissionAuthGroupPermissions(mods ...qm.QueryMod) authGroupPermissionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"auth_group_permissions\".\"permission_id\"=?", o.ID),
	)

	return AuthGroupPermissions(queryMods...)
}

// PermissionAuthUserUserPermissions retrieves all the auth_user_user_permission's AuthUserUserPermissions with an executor via permission_id column.
func (o *AuthPermission) PermissionAuthUserUserPermissions(mods ...qm.QueryMod) authUserUserPermissionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"auth_user_user_permissions\".\"permission_id\"=?", o.ID),
	)

	return AuthUserUserPermissions(queryMods...)
}

// LoadContentType allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (authPermissionL) LoadContentType(e boil.Executor, singular bool, maybeAuthPermission interface{}, mods queries.Applicator) error {
	var slice []*AuthPermission
	var object *AuthPermission

	if singular {
		var ok bool
		object, ok = maybeAuthPermission.(*AuthPermission)
		if !ok {
			object = new(AuthPermission)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthPermission))
			}
		}
	} else {
		s, ok := maybeAuthPermission.(*[]*AuthPermission)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthPermission))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authPermissionR{}
		}
		args[object.ContentTypeID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authPermissionR{}
			}

			args[obj.ContentTypeID] = struct{}{}

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
		qm.From(`django_content_type`),
		qm.WhereIn(`django_content_type.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load DjangoContentType")
	}

	var resultSlice []*DjangoContentType
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice DjangoContentType")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for django_content_type")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for django_content_type")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ContentType = foreign
		if foreign.R == nil {
			foreign.R = &djangoContentTypeR{}
		}
		foreign.R.ContentTypeAuthPermissions = append(foreign.R.ContentTypeAuthPermissions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ContentTypeID == foreign.ID {
				local.R.ContentType = foreign
				if foreign.R == nil {
					foreign.R = &djangoContentTypeR{}
				}
				foreign.R.ContentTypeAuthPermissions = append(foreign.R.ContentTypeAuthPermissions, local)
				break
			}
		}
	}

	return nil
}

// LoadPermissionAuthGroupPermissions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (authPermissionL) LoadPermissionAuthGroupPermissions(e boil.Executor, singular bool, maybeAuthPermission interface{}, mods queries.Applicator) error {
	var slice []*AuthPermission
	var object *AuthPermission

	if singular {
		var ok bool
		object, ok = maybeAuthPermission.(*AuthPermission)
		if !ok {
			object = new(AuthPermission)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthPermission))
			}
		}
	} else {
		s, ok := maybeAuthPermission.(*[]*AuthPermission)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthPermission))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authPermissionR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authPermissionR{}
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
		qm.From(`auth_group_permissions`),
		qm.WhereIn(`auth_group_permissions.permission_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load auth_group_permissions")
	}

	var resultSlice []*AuthGroupPermission
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice auth_group_permissions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on auth_group_permissions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for auth_group_permissions")
	}

	if singular {
		object.R.PermissionAuthGroupPermissions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &authGroupPermissionR{}
			}
			foreign.R.Permission = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PermissionID {
				local.R.PermissionAuthGroupPermissions = append(local.R.PermissionAuthGroupPermissions, foreign)
				if foreign.R == nil {
					foreign.R = &authGroupPermissionR{}
				}
				foreign.R.Permission = local
				break
			}
		}
	}

	return nil
}

// LoadPermissionAuthUserUserPermissions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (authPermissionL) LoadPermissionAuthUserUserPermissions(e boil.Executor, singular bool, maybeAuthPermission interface{}, mods queries.Applicator) error {
	var slice []*AuthPermission
	var object *AuthPermission

	if singular {
		var ok bool
		object, ok = maybeAuthPermission.(*AuthPermission)
		if !ok {
			object = new(AuthPermission)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthPermission))
			}
		}
	} else {
		s, ok := maybeAuthPermission.(*[]*AuthPermission)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthPermission)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthPermission))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authPermissionR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authPermissionR{}
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
		qm.From(`auth_user_user_permissions`),
		qm.WhereIn(`auth_user_user_permissions.permission_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load auth_user_user_permissions")
	}

	var resultSlice []*AuthUserUserPermission
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice auth_user_user_permissions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on auth_user_user_permissions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for auth_user_user_permissions")
	}

	if singular {
		object.R.PermissionAuthUserUserPermissions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &authUserUserPermissionR{}
			}
			foreign.R.Permission = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PermissionID {
				local.R.PermissionAuthUserUserPermissions = append(local.R.PermissionAuthUserUserPermissions, foreign)
				if foreign.R == nil {
					foreign.R = &authUserUserPermissionR{}
				}
				foreign.R.Permission = local
				break
			}
		}
	}

	return nil
}

// SetContentType of the authPermission to the related item.
// Sets o.R.ContentType to related.
// Adds o to related.R.ContentTypeAuthPermissions.
func (o *AuthPermission) SetContentType(exec boil.Executor, insert bool, related *DjangoContentType) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"auth_permission\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"content_type_id"}),
		strmangle.WhereClause("\"", "\"", 2, authPermissionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ContentTypeID = related.ID
	if o.R == nil {
		o.R = &authPermissionR{
			ContentType: related,
		}
	} else {
		o.R.ContentType = related
	}

	if related.R == nil {
		related.R = &djangoContentTypeR{
			ContentTypeAuthPermissions: AuthPermissionSlice{o},
		}
	} else {
		related.R.ContentTypeAuthPermissions = append(related.R.ContentTypeAuthPermissions, o)
	}

	return nil
}

// AddPermissionAuthGroupPermissions adds the given related objects to the existing relationships
// of the auth_permission, optionally inserting them as new records.
// Appends related to o.R.PermissionAuthGroupPermissions.
// Sets related.R.Permission appropriately.
func (o *AuthPermission) AddPermissionAuthGroupPermissions(exec boil.Executor, insert bool, related ...*AuthGroupPermission) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PermissionID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"auth_group_permissions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"permission_id"}),
				strmangle.WhereClause("\"", "\"", 2, authGroupPermissionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PermissionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &authPermissionR{
			PermissionAuthGroupPermissions: related,
		}
	} else {
		o.R.PermissionAuthGroupPermissions = append(o.R.PermissionAuthGroupPermissions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &authGroupPermissionR{
				Permission: o,
			}
		} else {
			rel.R.Permission = o
		}
	}
	return nil
}

// AddPermissionAuthUserUserPermissions adds the given related objects to the existing relationships
// of the auth_permission, optionally inserting them as new records.
// Appends related to o.R.PermissionAuthUserUserPermissions.
// Sets related.R.Permission appropriately.
func (o *AuthPermission) AddPermissionAuthUserUserPermissions(exec boil.Executor, insert bool, related ...*AuthUserUserPermission) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PermissionID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"auth_user_user_permissions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"permission_id"}),
				strmangle.WhereClause("\"", "\"", 2, authUserUserPermissionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PermissionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &authPermissionR{
			PermissionAuthUserUserPermissions: related,
		}
	} else {
		o.R.PermissionAuthUserUserPermissions = append(o.R.PermissionAuthUserUserPermissions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &authUserUserPermissionR{
				Permission: o,
			}
		} else {
			rel.R.Permission = o
		}
	}
	return nil
}

// AuthPermissions retrieves all the records using an executor.
func AuthPermissions(mods ...qm.QueryMod) authPermissionQuery {
	mods = append(mods, qm.From("\"auth_permission\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"auth_permission\".*"})
	}

	return authPermissionQuery{q}
}

// FindAuthPermission retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAuthPermission(exec boil.Executor, iD int, selectCols ...string) (*AuthPermission, error) {
	authPermissionObj := &AuthPermission{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"auth_permission\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, authPermissionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: unable to select from auth_permission")
	}

	return authPermissionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AuthPermission) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("schema: no auth_permission provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(authPermissionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	authPermissionInsertCacheMut.RLock()
	cache, cached := authPermissionInsertCache[key]
	authPermissionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			authPermissionAllColumns,
			authPermissionColumnsWithDefault,
			authPermissionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(authPermissionType, authPermissionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(authPermissionType, authPermissionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"auth_permission\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"auth_permission\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "schema: unable to insert into auth_permission")
	}

	if !cached {
		authPermissionInsertCacheMut.Lock()
		authPermissionInsertCache[key] = cache
		authPermissionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the AuthPermission.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AuthPermission) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	authPermissionUpdateCacheMut.RLock()
	cache, cached := authPermissionUpdateCache[key]
	authPermissionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			authPermissionAllColumns,
			authPermissionPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("schema: unable to update auth_permission, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"auth_permission\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, authPermissionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(authPermissionType, authPermissionMapping, append(wl, authPermissionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "schema: unable to update auth_permission row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by update for auth_permission")
	}

	if !cached {
		authPermissionUpdateCacheMut.Lock()
		authPermissionUpdateCache[key] = cache
		authPermissionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q authPermissionQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all for auth_permission")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected for auth_permission")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AuthPermissionSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authPermissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"auth_permission\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, authPermissionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all in authPermission slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected all in update all authPermission")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AuthPermission) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("schema: no auth_permission provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(authPermissionColumnsWithDefault, o)

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

	authPermissionUpsertCacheMut.RLock()
	cache, cached := authPermissionUpsertCache[key]
	authPermissionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			authPermissionAllColumns,
			authPermissionColumnsWithDefault,
			authPermissionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			authPermissionAllColumns,
			authPermissionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("schema: unable to upsert auth_permission, could not build update column list")
		}

		ret := strmangle.SetComplement(authPermissionAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(authPermissionPrimaryKeyColumns) == 0 {
				return errors.New("schema: unable to upsert auth_permission, could not build conflict column list")
			}

			conflict = make([]string, len(authPermissionPrimaryKeyColumns))
			copy(conflict, authPermissionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"auth_permission\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(authPermissionType, authPermissionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(authPermissionType, authPermissionMapping, ret)
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
		return errors.Wrap(err, "schema: unable to upsert auth_permission")
	}

	if !cached {
		authPermissionUpsertCacheMut.Lock()
		authPermissionUpsertCache[key] = cache
		authPermissionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single AuthPermission record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AuthPermission) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("schema: no AuthPermission provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), authPermissionPrimaryKeyMapping)
	sql := "DELETE FROM \"auth_permission\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete from auth_permission")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by delete for auth_permission")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q authPermissionQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("schema: no authPermissionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from auth_permission")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for auth_permission")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AuthPermissionSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authPermissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"auth_permission\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, authPermissionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from authPermission slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for auth_permission")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AuthPermission) Reload(exec boil.Executor) error {
	ret, err := FindAuthPermission(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthPermissionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AuthPermissionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authPermissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"auth_permission\".* FROM \"auth_permission\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, authPermissionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "schema: unable to reload all in AuthPermissionSlice")
	}

	*o = slice

	return nil
}

// AuthPermissionExists checks if the AuthPermission row exists.
func AuthPermissionExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"auth_permission\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "schema: unable to check if auth_permission exists")
	}

	return exists, nil
}

// Exists checks if the AuthPermission row exists.
func (o *AuthPermission) Exists(exec boil.Executor) (bool, error) {
	return AuthPermissionExists(exec, o.ID)
}
