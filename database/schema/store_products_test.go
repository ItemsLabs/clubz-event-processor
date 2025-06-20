// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package schema

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testStoreProducts(t *testing.T) {
	t.Parallel()

	query := StoreProducts()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testStoreProductsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testStoreProductsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := StoreProducts().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testStoreProductsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := StoreProductSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testStoreProductsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := StoreProductExists(tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if StoreProduct exists: %s", err)
	}
	if !e {
		t.Errorf("Expected StoreProductExists to return true, but got false.")
	}
}

func testStoreProductsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	storeProductFound, err := FindStoreProduct(tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if storeProductFound == nil {
		t.Error("want a record, got nil")
	}
}

func testStoreProductsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = StoreProducts().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testStoreProductsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := StoreProducts().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testStoreProductsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	storeProductOne := &StoreProduct{}
	storeProductTwo := &StoreProduct{}
	if err = randomize.Struct(seed, storeProductOne, storeProductDBTypes, false, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}
	if err = randomize.Struct(seed, storeProductTwo, storeProductDBTypes, false, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = storeProductOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = storeProductTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := StoreProducts().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testStoreProductsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	storeProductOne := &StoreProduct{}
	storeProductTwo := &StoreProduct{}
	if err = randomize.Struct(seed, storeProductOne, storeProductDBTypes, false, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}
	if err = randomize.Struct(seed, storeProductTwo, storeProductDBTypes, false, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = storeProductOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = storeProductTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testStoreProductsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testStoreProductsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(storeProductColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testStoreProductToManyProductStoreProductTransactions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a StoreProduct
	var b, c StoreProductTransaction

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, storeProductTransactionDBTypes, false, storeProductTransactionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, storeProductTransactionDBTypes, false, storeProductTransactionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.ProductID = a.ID
	c.ProductID = a.ID

	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.ProductStoreProductTransactions().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ProductID == b.ProductID {
			bFound = true
		}
		if v.ProductID == c.ProductID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StoreProductSlice{&a}
	if err = a.L.LoadProductStoreProductTransactions(tx, false, (*[]*StoreProduct)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ProductStoreProductTransactions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ProductStoreProductTransactions = nil
	if err = a.L.LoadProductStoreProductTransactions(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ProductStoreProductTransactions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testStoreProductToManyAddOpProductStoreProductTransactions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a StoreProduct
	var b, c, d, e StoreProductTransaction

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, storeProductDBTypes, false, strmangle.SetComplement(storeProductPrimaryKeyColumns, storeProductColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*StoreProductTransaction{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, storeProductTransactionDBTypes, false, strmangle.SetComplement(storeProductTransactionPrimaryKeyColumns, storeProductTransactionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*StoreProductTransaction{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddProductStoreProductTransactions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ProductID {
			t.Error("foreign key was wrong value", a.ID, first.ProductID)
		}
		if a.ID != second.ProductID {
			t.Error("foreign key was wrong value", a.ID, second.ProductID)
		}

		if first.R.Product != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Product != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ProductStoreProductTransactions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ProductStoreProductTransactions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ProductStoreProductTransactions().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStoreProductsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testStoreProductsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := StoreProductSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testStoreProductsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := StoreProducts().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	storeProductDBTypes = map[string]string{`ID`: `uuid`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`, `StoreProductID`: `character varying`, `AppleProductID`: `character varying`, `GoogleProductID`: `character varying`, `Description`: `text`, `Price`: `double precision`, `Currency`: `character varying`, `Active`: `boolean`, `ProductType`: `character varying`}
	_                   = bytes.MinRead
)

func testStoreProductsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(storeProductPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(storeProductAllColumns) == len(storeProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testStoreProductsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(storeProductAllColumns) == len(storeProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &StoreProduct{}
	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, storeProductDBTypes, true, storeProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(storeProductAllColumns, storeProductPrimaryKeyColumns) {
		fields = storeProductAllColumns
	} else {
		fields = strmangle.SetComplement(
			storeProductAllColumns,
			storeProductPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := StoreProductSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testStoreProductsUpsert(t *testing.T) {
	t.Parallel()

	if len(storeProductAllColumns) == len(storeProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := StoreProduct{}
	if err = randomize.Struct(seed, &o, storeProductDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert StoreProduct: %s", err)
	}

	count, err := StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, storeProductDBTypes, false, storeProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize StoreProduct struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert StoreProduct: %s", err)
	}

	count, err = StoreProducts().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
