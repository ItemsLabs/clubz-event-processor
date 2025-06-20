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

func testChatMessages(t *testing.T) {
	t.Parallel()

	query := ChatMessages()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testChatMessagesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
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

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChatMessagesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ChatMessages().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChatMessagesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ChatMessageSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChatMessagesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ChatMessageExists(tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ChatMessage exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ChatMessageExists to return true, but got false.")
	}
}

func testChatMessagesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	chatMessageFound, err := FindChatMessage(tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if chatMessageFound == nil {
		t.Error("want a record, got nil")
	}
}

func testChatMessagesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ChatMessages().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testChatMessagesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ChatMessages().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testChatMessagesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chatMessageOne := &ChatMessage{}
	chatMessageTwo := &ChatMessage{}
	if err = randomize.Struct(seed, chatMessageOne, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}
	if err = randomize.Struct(seed, chatMessageTwo, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = chatMessageOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = chatMessageTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ChatMessages().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testChatMessagesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	chatMessageOne := &ChatMessage{}
	chatMessageTwo := &ChatMessage{}
	if err = randomize.Struct(seed, chatMessageOne, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}
	if err = randomize.Struct(seed, chatMessageTwo, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = chatMessageOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = chatMessageTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testChatMessagesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChatMessagesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(chatMessageColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChatMessageToOneChatRoomUsingRoom(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local ChatMessage
	var foreign ChatRoom

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, chatRoomDBTypes, false, chatRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatRoom struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.RoomID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Room().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChatMessageSlice{&local}
	if err = local.L.LoadRoom(tx, false, (*[]*ChatMessage)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Room == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Room = nil
	if err = local.L.LoadRoom(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Room == nil {
		t.Error("struct should have been eager loaded")
	}

}

func testChatMessageToOneUserUsingSender(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local ChatMessage
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, chatMessageDBTypes, false, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SenderID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Sender().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChatMessageSlice{&local}
	if err = local.L.LoadSender(tx, false, (*[]*ChatMessage)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Sender == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Sender = nil
	if err = local.L.LoadSender(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Sender == nil {
		t.Error("struct should have been eager loaded")
	}

}

func testChatMessageToOneSetOpChatRoomUsingRoom(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a ChatMessage
	var b, c ChatRoom

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, chatMessageDBTypes, false, strmangle.SetComplement(chatMessagePrimaryKeyColumns, chatMessageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, chatRoomDBTypes, false, strmangle.SetComplement(chatRoomPrimaryKeyColumns, chatRoomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, chatRoomDBTypes, false, strmangle.SetComplement(chatRoomPrimaryKeyColumns, chatRoomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ChatRoom{&b, &c} {
		err = a.SetRoom(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Room != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.RoomChatMessages[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.RoomID != x.ID {
			t.Error("foreign key was wrong value", a.RoomID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.RoomID))
		reflect.Indirect(reflect.ValueOf(&a.RoomID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.RoomID != x.ID {
			t.Error("foreign key was wrong value", a.RoomID, x.ID)
		}
	}
}
func testChatMessageToOneSetOpUserUsingSender(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a ChatMessage
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, chatMessageDBTypes, false, strmangle.SetComplement(chatMessagePrimaryKeyColumns, chatMessageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetSender(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Sender != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.SenderChatMessages[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SenderID != x.ID {
			t.Error("foreign key was wrong value", a.SenderID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SenderID))
		reflect.Indirect(reflect.ValueOf(&a.SenderID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SenderID != x.ID {
			t.Error("foreign key was wrong value", a.SenderID, x.ID)
		}
	}
}

func testChatMessagesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
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

func testChatMessagesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ChatMessageSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testChatMessagesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ChatMessages().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	chatMessageDBTypes = map[string]string{`ID`: `uuid`, `Message`: `text`, `RoomID`: `uuid`, `SenderID`: `uuid`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`}
	_                  = bytes.MinRead
)

func testChatMessagesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(chatMessagePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(chatMessageAllColumns) == len(chatMessagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testChatMessagesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(chatMessageAllColumns) == len(chatMessagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ChatMessage{}
	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, chatMessageDBTypes, true, chatMessagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(chatMessageAllColumns, chatMessagePrimaryKeyColumns) {
		fields = chatMessageAllColumns
	} else {
		fields = strmangle.SetComplement(
			chatMessageAllColumns,
			chatMessagePrimaryKeyColumns,
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

	slice := ChatMessageSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testChatMessagesUpsert(t *testing.T) {
	t.Parallel()

	if len(chatMessageAllColumns) == len(chatMessagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ChatMessage{}
	if err = randomize.Struct(seed, &o, chatMessageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ChatMessage: %s", err)
	}

	count, err := ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, chatMessageDBTypes, false, chatMessagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChatMessage struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ChatMessage: %s", err)
	}

	count, err = ChatMessages().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
