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

// Country is an object representing the database table.
type Country struct {
	ID        string      `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ImportID  null.String `boil:"import_id" json:"import_id,omitempty" toml:"import_id" yaml:"import_id,omitempty"`
	Name      string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Iso       null.String `boil:"iso" json:"iso,omitempty" toml:"iso" yaml:"iso,omitempty"`

	R *countryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L countryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CountryColumns = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	ImportID  string
	Name      string
	Iso       string
}{
	ID:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	ImportID:  "import_id",
	Name:      "name",
	Iso:       "iso",
}

var CountryTableColumns = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	ImportID  string
	Name      string
	Iso       string
}{
	ID:        "countries.id",
	CreatedAt: "countries.created_at",
	UpdatedAt: "countries.updated_at",
	ImportID:  "countries.import_id",
	Name:      "countries.name",
	Iso:       "countries.iso",
}

// Generated where

var CountryWhere = struct {
	ID        whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
	ImportID  whereHelpernull_String
	Name      whereHelperstring
	Iso       whereHelpernull_String
}{
	ID:        whereHelperstring{field: "\"countries\".\"id\""},
	CreatedAt: whereHelpertime_Time{field: "\"countries\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"countries\".\"updated_at\""},
	ImportID:  whereHelpernull_String{field: "\"countries\".\"import_id\""},
	Name:      whereHelperstring{field: "\"countries\".\"name\""},
	Iso:       whereHelpernull_String{field: "\"countries\".\"iso\""},
}

// CountryRels is where relationship names are stored.
var CountryRels = struct {
	Leaderboards  string
	PlayerBuckets string
	Teams         string
}{
	Leaderboards:  "Leaderboards",
	PlayerBuckets: "PlayerBuckets",
	Teams:         "Teams",
}

// countryR is where relationships are stored.
type countryR struct {
	Leaderboards  LeaderboardSlice  `boil:"Leaderboards" json:"Leaderboards" toml:"Leaderboards" yaml:"Leaderboards"`
	PlayerBuckets PlayerBucketSlice `boil:"PlayerBuckets" json:"PlayerBuckets" toml:"PlayerBuckets" yaml:"PlayerBuckets"`
	Teams         TeamSlice         `boil:"Teams" json:"Teams" toml:"Teams" yaml:"Teams"`
}

// NewStruct creates a new relationship struct
func (*countryR) NewStruct() *countryR {
	return &countryR{}
}

func (r *countryR) GetLeaderboards() LeaderboardSlice {
	if r == nil {
		return nil
	}
	return r.Leaderboards
}

func (r *countryR) GetPlayerBuckets() PlayerBucketSlice {
	if r == nil {
		return nil
	}
	return r.PlayerBuckets
}

func (r *countryR) GetTeams() TeamSlice {
	if r == nil {
		return nil
	}
	return r.Teams
}

// countryL is where Load methods for each relationship are stored.
type countryL struct{}

var (
	countryAllColumns            = []string{"id", "created_at", "updated_at", "import_id", "name", "iso"}
	countryColumnsWithoutDefault = []string{"id", "created_at", "updated_at", "name"}
	countryColumnsWithDefault    = []string{"import_id", "iso"}
	countryPrimaryKeyColumns     = []string{"id"}
	countryGeneratedColumns      = []string{}
)

type (
	// CountrySlice is an alias for a slice of pointers to Country.
	// This should almost always be used instead of []Country.
	CountrySlice []*Country

	countryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	countryType                 = reflect.TypeOf(&Country{})
	countryMapping              = queries.MakeStructMapping(countryType)
	countryPrimaryKeyMapping, _ = queries.BindMapping(countryType, countryMapping, countryPrimaryKeyColumns)
	countryInsertCacheMut       sync.RWMutex
	countryInsertCache          = make(map[string]insertCache)
	countryUpdateCacheMut       sync.RWMutex
	countryUpdateCache          = make(map[string]updateCache)
	countryUpsertCacheMut       sync.RWMutex
	countryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single country record from the query.
func (q countryQuery) One(exec boil.Executor) (*Country, error) {
	o := &Country{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: failed to execute a one query for countries")
	}

	return o, nil
}

// All returns all Country records from the query.
func (q countryQuery) All(exec boil.Executor) (CountrySlice, error) {
	var o []*Country

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "schema: failed to assign all query results to Country slice")
	}

	return o, nil
}

// Count returns the count of all Country records in the query.
func (q countryQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to count countries rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q countryQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "schema: failed to check if countries exists")
	}

	return count > 0, nil
}

// Leaderboards retrieves all the leaderboard's Leaderboards with an executor.
func (o *Country) Leaderboards(mods ...qm.QueryMod) leaderboardQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"leaderboards\".\"country_id\"=?", o.ID),
	)

	return Leaderboards(queryMods...)
}

// PlayerBuckets retrieves all the player_bucket's PlayerBuckets with an executor.
func (o *Country) PlayerBuckets(mods ...qm.QueryMod) playerBucketQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"player_bucket\".\"country_id\"=?", o.ID),
	)

	return PlayerBuckets(queryMods...)
}

// Teams retrieves all the team's Teams with an executor.
func (o *Country) Teams(mods ...qm.QueryMod) teamQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"teams\".\"country_id\"=?", o.ID),
	)

	return Teams(queryMods...)
}

// LoadLeaderboards allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (countryL) LoadLeaderboards(e boil.Executor, singular bool, maybeCountry interface{}, mods queries.Applicator) error {
	var slice []*Country
	var object *Country

	if singular {
		var ok bool
		object, ok = maybeCountry.(*Country)
		if !ok {
			object = new(Country)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCountry))
			}
		}
	} else {
		s, ok := maybeCountry.(*[]*Country)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCountry))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &countryR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &countryR{}
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
		qm.From(`leaderboards`),
		qm.WhereIn(`leaderboards.country_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load leaderboards")
	}

	var resultSlice []*Leaderboard
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice leaderboards")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on leaderboards")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for leaderboards")
	}

	if singular {
		object.R.Leaderboards = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &leaderboardR{}
			}
			foreign.R.Country = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.CountryID) {
				local.R.Leaderboards = append(local.R.Leaderboards, foreign)
				if foreign.R == nil {
					foreign.R = &leaderboardR{}
				}
				foreign.R.Country = local
				break
			}
		}
	}

	return nil
}

// LoadPlayerBuckets allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (countryL) LoadPlayerBuckets(e boil.Executor, singular bool, maybeCountry interface{}, mods queries.Applicator) error {
	var slice []*Country
	var object *Country

	if singular {
		var ok bool
		object, ok = maybeCountry.(*Country)
		if !ok {
			object = new(Country)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCountry))
			}
		}
	} else {
		s, ok := maybeCountry.(*[]*Country)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCountry))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &countryR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &countryR{}
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
		qm.From(`player_bucket`),
		qm.WhereIn(`player_bucket.country_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load player_bucket")
	}

	var resultSlice []*PlayerBucket
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice player_bucket")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on player_bucket")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for player_bucket")
	}

	if singular {
		object.R.PlayerBuckets = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &playerBucketR{}
			}
			foreign.R.Country = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.CountryID) {
				local.R.PlayerBuckets = append(local.R.PlayerBuckets, foreign)
				if foreign.R == nil {
					foreign.R = &playerBucketR{}
				}
				foreign.R.Country = local
				break
			}
		}
	}

	return nil
}

// LoadTeams allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (countryL) LoadTeams(e boil.Executor, singular bool, maybeCountry interface{}, mods queries.Applicator) error {
	var slice []*Country
	var object *Country

	if singular {
		var ok bool
		object, ok = maybeCountry.(*Country)
		if !ok {
			object = new(Country)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeCountry))
			}
		}
	} else {
		s, ok := maybeCountry.(*[]*Country)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeCountry)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeCountry))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &countryR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &countryR{}
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
		qm.From(`teams`),
		qm.WhereIn(`teams.country_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load teams")
	}

	var resultSlice []*Team
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice teams")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on teams")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for teams")
	}

	if singular {
		object.R.Teams = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &teamR{}
			}
			foreign.R.Country = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.CountryID) {
				local.R.Teams = append(local.R.Teams, foreign)
				if foreign.R == nil {
					foreign.R = &teamR{}
				}
				foreign.R.Country = local
				break
			}
		}
	}

	return nil
}

// AddLeaderboards adds the given related objects to the existing relationships
// of the country, optionally inserting them as new records.
// Appends related to o.R.Leaderboards.
// Sets related.R.Country appropriately.
func (o *Country) AddLeaderboards(exec boil.Executor, insert bool, related ...*Leaderboard) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.CountryID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"leaderboards\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"country_id"}),
				strmangle.WhereClause("\"", "\"", 2, leaderboardPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.CountryID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &countryR{
			Leaderboards: related,
		}
	} else {
		o.R.Leaderboards = append(o.R.Leaderboards, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &leaderboardR{
				Country: o,
			}
		} else {
			rel.R.Country = o
		}
	}
	return nil
}

// SetLeaderboards removes all previously related items of the
// country replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Country's Leaderboards accordingly.
// Replaces o.R.Leaderboards with related.
// Sets related.R.Country's Leaderboards accordingly.
func (o *Country) SetLeaderboards(exec boil.Executor, insert bool, related ...*Leaderboard) error {
	query := "update \"leaderboards\" set \"country_id\" = null where \"country_id\" = $1"
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
		for _, rel := range o.R.Leaderboards {
			queries.SetScanner(&rel.CountryID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Country = nil
		}
		o.R.Leaderboards = nil
	}

	return o.AddLeaderboards(exec, insert, related...)
}

// RemoveLeaderboards relationships from objects passed in.
// Removes related items from R.Leaderboards (uses pointer comparison, removal does not keep order)
// Sets related.R.Country.
func (o *Country) RemoveLeaderboards(exec boil.Executor, related ...*Leaderboard) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.CountryID, nil)
		if rel.R != nil {
			rel.R.Country = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("country_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Leaderboards {
			if rel != ri {
				continue
			}

			ln := len(o.R.Leaderboards)
			if ln > 1 && i < ln-1 {
				o.R.Leaderboards[i] = o.R.Leaderboards[ln-1]
			}
			o.R.Leaderboards = o.R.Leaderboards[:ln-1]
			break
		}
	}

	return nil
}

// AddPlayerBuckets adds the given related objects to the existing relationships
// of the country, optionally inserting them as new records.
// Appends related to o.R.PlayerBuckets.
// Sets related.R.Country appropriately.
func (o *Country) AddPlayerBuckets(exec boil.Executor, insert bool, related ...*PlayerBucket) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.CountryID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"player_bucket\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"country_id"}),
				strmangle.WhereClause("\"", "\"", 2, playerBucketPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.CountryID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &countryR{
			PlayerBuckets: related,
		}
	} else {
		o.R.PlayerBuckets = append(o.R.PlayerBuckets, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &playerBucketR{
				Country: o,
			}
		} else {
			rel.R.Country = o
		}
	}
	return nil
}

// SetPlayerBuckets removes all previously related items of the
// country replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Country's PlayerBuckets accordingly.
// Replaces o.R.PlayerBuckets with related.
// Sets related.R.Country's PlayerBuckets accordingly.
func (o *Country) SetPlayerBuckets(exec boil.Executor, insert bool, related ...*PlayerBucket) error {
	query := "update \"player_bucket\" set \"country_id\" = null where \"country_id\" = $1"
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
		for _, rel := range o.R.PlayerBuckets {
			queries.SetScanner(&rel.CountryID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Country = nil
		}
		o.R.PlayerBuckets = nil
	}

	return o.AddPlayerBuckets(exec, insert, related...)
}

// RemovePlayerBuckets relationships from objects passed in.
// Removes related items from R.PlayerBuckets (uses pointer comparison, removal does not keep order)
// Sets related.R.Country.
func (o *Country) RemovePlayerBuckets(exec boil.Executor, related ...*PlayerBucket) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.CountryID, nil)
		if rel.R != nil {
			rel.R.Country = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("country_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.PlayerBuckets {
			if rel != ri {
				continue
			}

			ln := len(o.R.PlayerBuckets)
			if ln > 1 && i < ln-1 {
				o.R.PlayerBuckets[i] = o.R.PlayerBuckets[ln-1]
			}
			o.R.PlayerBuckets = o.R.PlayerBuckets[:ln-1]
			break
		}
	}

	return nil
}

// AddTeams adds the given related objects to the existing relationships
// of the country, optionally inserting them as new records.
// Appends related to o.R.Teams.
// Sets related.R.Country appropriately.
func (o *Country) AddTeams(exec boil.Executor, insert bool, related ...*Team) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.CountryID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"teams\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"country_id"}),
				strmangle.WhereClause("\"", "\"", 2, teamPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.CountryID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &countryR{
			Teams: related,
		}
	} else {
		o.R.Teams = append(o.R.Teams, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &teamR{
				Country: o,
			}
		} else {
			rel.R.Country = o
		}
	}
	return nil
}

// SetTeams removes all previously related items of the
// country replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Country's Teams accordingly.
// Replaces o.R.Teams with related.
// Sets related.R.Country's Teams accordingly.
func (o *Country) SetTeams(exec boil.Executor, insert bool, related ...*Team) error {
	query := "update \"teams\" set \"country_id\" = null where \"country_id\" = $1"
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
		for _, rel := range o.R.Teams {
			queries.SetScanner(&rel.CountryID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Country = nil
		}
		o.R.Teams = nil
	}

	return o.AddTeams(exec, insert, related...)
}

// RemoveTeams relationships from objects passed in.
// Removes related items from R.Teams (uses pointer comparison, removal does not keep order)
// Sets related.R.Country.
func (o *Country) RemoveTeams(exec boil.Executor, related ...*Team) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.CountryID, nil)
		if rel.R != nil {
			rel.R.Country = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("country_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Teams {
			if rel != ri {
				continue
			}

			ln := len(o.R.Teams)
			if ln > 1 && i < ln-1 {
				o.R.Teams[i] = o.R.Teams[ln-1]
			}
			o.R.Teams = o.R.Teams[:ln-1]
			break
		}
	}

	return nil
}

// Countries retrieves all the records using an executor.
func Countries(mods ...qm.QueryMod) countryQuery {
	mods = append(mods, qm.From("\"countries\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"countries\".*"})
	}

	return countryQuery{q}
}

// FindCountry retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCountry(exec boil.Executor, iD string, selectCols ...string) (*Country, error) {
	countryObj := &Country{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"countries\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, countryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "schema: unable to select from countries")
	}

	return countryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Country) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("schema: no countries provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(countryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	countryInsertCacheMut.RLock()
	cache, cached := countryInsertCache[key]
	countryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			countryAllColumns,
			countryColumnsWithDefault,
			countryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(countryType, countryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(countryType, countryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"countries\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"countries\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "schema: unable to insert into countries")
	}

	if !cached {
		countryInsertCacheMut.Lock()
		countryInsertCache[key] = cache
		countryInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Country.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Country) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	key := makeCacheKey(columns, nil)
	countryUpdateCacheMut.RLock()
	cache, cached := countryUpdateCache[key]
	countryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			countryAllColumns,
			countryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("schema: unable to update countries, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"countries\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, countryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(countryType, countryMapping, append(wl, countryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "schema: unable to update countries row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by update for countries")
	}

	if !cached {
		countryUpdateCacheMut.Lock()
		countryUpdateCache[key] = cache
		countryUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q countryQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all for countries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected for countries")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CountrySlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), countryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"countries\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, countryPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to update all in country slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to retrieve rows affected all in update all country")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Country) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("schema: no countries provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	nzDefaults := queries.NonZeroDefaultSet(countryColumnsWithDefault, o)

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

	countryUpsertCacheMut.RLock()
	cache, cached := countryUpsertCache[key]
	countryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			countryAllColumns,
			countryColumnsWithDefault,
			countryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			countryAllColumns,
			countryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("schema: unable to upsert countries, could not build update column list")
		}

		ret := strmangle.SetComplement(countryAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(countryPrimaryKeyColumns) == 0 {
				return errors.New("schema: unable to upsert countries, could not build conflict column list")
			}

			conflict = make([]string, len(countryPrimaryKeyColumns))
			copy(conflict, countryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"countries\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(countryType, countryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(countryType, countryMapping, ret)
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
		return errors.Wrap(err, "schema: unable to upsert countries")
	}

	if !cached {
		countryUpsertCacheMut.Lock()
		countryUpsertCache[key] = cache
		countryUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Country record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Country) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("schema: no Country provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), countryPrimaryKeyMapping)
	sql := "DELETE FROM \"countries\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete from countries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by delete for countries")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q countryQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("schema: no countryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from countries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for countries")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CountrySlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), countryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"countries\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, countryPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "schema: unable to delete all from country slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "schema: failed to get rows affected by deleteall for countries")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Country) Reload(exec boil.Executor) error {
	ret, err := FindCountry(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CountrySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CountrySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), countryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"countries\".* FROM \"countries\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, countryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "schema: unable to reload all in CountrySlice")
	}

	*o = slice

	return nil
}

// CountryExists checks if the Country row exists.
func CountryExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"countries\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "schema: unable to check if countries exists")
	}

	return exists, nil
}

// Exists checks if the Country row exists.
func (o *Country) Exists(exec boil.Executor) (bool, error) {
	return CountryExists(exec, o.ID)
}
