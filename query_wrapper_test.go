package wysci

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestDBTypeString(t *testing.T) {
	x := DBText

	if x.String() != "Text" {
		t.Errorf("Expected 'Text' but got %s", x.String())
	}

	if x = DBNumber; x.String() != "Number" {
		t.Errorf("Expected 'Number' but got %s", x.String())
	}

	if x = DBDate; x.String() != "Date" {
		t.Errorf("Expected 'Date' but got %s", x.String())
	}

	if x = DBBytes; x.String() != "Bytes" {
		t.Errorf("Expected 'Bytes' but got %s", x.String())
	}

	if x = DBTime; x.String() != "Time" {
		t.Errorf("Expected 'Time' but got %s", x.String())
	}

	if x = DBUnknown; x.String() != "Unknown" {
		t.Errorf("Expected 'Unknown' but got %s", x.String())
	}
}

func TestExecuteQuery(t *testing.T) {
	q, err := ExecuteQuery(testConn, "select id, name, some_date from test_simple")
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()

	cols := q.Columns()
	if len(cols) != 3 {
		t.Errorf("Expected 3 columns for query columns but got %d", len(cols))
	}

	if cols[0] != "id" {
		t.Errorf("Expected first column to be 'id' but got %s", cols[0])
	}

	if cols[1] != "name" {
		t.Errorf("Expected second column to be 'name' but got %s", cols[1])
	}

	if cols[2] != "some_date" {
		t.Errorf("Expected third column to be 'some_date' but got %s", cols[2])
	}
}

type testSimpleIter struct {
	RowCount int
}

func (t *testSimpleIter) Iteration(row []interface{}) (bool, error) {
	t.RowCount++
	return false, nil
}

type testStopIter struct {
	RowCount int
}

func (t *testStopIter) Iteration(row []interface{}) (bool, error) {
	t.RowCount++
	if t.RowCount > 2 {
		return true, nil
	}
	return false, nil
}

func TestSimpleForEach(t *testing.T) {
	q, err := ExecuteQuery(testConn, "select id from test_simple")
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()

	si := &testSimpleIter{}

	err = q.ForEach(si)
	if err != nil {
		t.Error(err)
	}

	if si.RowCount != 5 {
		t.Errorf("Expected 5 results but got %d", si.RowCount)
	}
}

func TestStopForEach(t *testing.T) {
	q, err := ExecuteQuery(testConn, "select id from test_simple")
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()

	si := &testStopIter{}
	q.ForEach(si)
	if err != nil {
		t.Error(err)
	}

	if si.RowCount != 3 {
		t.Errorf("Expected 3 results but got %d", si.RowCount)
	}
}

type testSimpleAccum struct {
}

func (ts *testSimpleAccum) Accumulate(current interface{}, row []interface{}) (interface{}, error) {
	total, ok := current.(int)
	if !ok {
		return nil, fmt.Errorf("Failed to cast current to int, %v instead", reflect.TypeOf(current))
	}

	idStr, ok := row[0].(*sql.NullString)
	if !ok {
		return nil, fmt.Errorf("Failed to cast buffer element to NullString, %v instead", reflect.TypeOf(row[0]))
	}

	if !idStr.Valid {
		return nil, fmt.Errorf("Expected a value for 'id' but got nil instead")
	}

	nextVal, err := strconv.Atoi(idStr.String)
	if err != nil {
		return nil, err
	}

	total += nextVal

	return total, nil
}

func TestSimpleAccum(t *testing.T) {
	q, err := ExecuteQuery(testConn, "select id from test_simple")
	if err != nil {
		t.Fatal(err)
	}

	s := &testSimpleAccum{}
	val, err := q.Accumulate(s, 0)
	if err != nil {
		t.Fatal(err)
	}

	total, ok := val.(int)
	if !ok {
		t.Errorf("Expected accumulator to return int but got %v instead", reflect.TypeOf(val))
	}

	if total != 15 {
		t.Errorf("Expected 15 but got %d", total)
	}
}
