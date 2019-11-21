package wysci

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

// DBType is used to provide high-level type information about result columns.
// Rather than worry about the different types of text fields, it's usually
// enough to know that a column is a text type.
type DBType int

const (
	// DBNumber is a numeric database column
	DBNumber DBType = 0
	// DBText is a text column
	DBText DBType = 1
	// DBDate is a date column
	DBDate DBType = 2
	// DBTime is a time column
	DBTime DBType = 3
	// DBBytes is a raw byte column
	DBBytes DBType = 4
	// DBUnknown is an unmpaped column type
	DBUnknown DBType = 999
)

// String implements the Stringer interface
func (d DBType) String() string {
	switch d {
	case DBNumber:
		return "Number"
	case DBText:
		return "Text"
	case DBDate:
		return "Date"
	case DBTime:
		return "Time"
	case DBBytes:
		return "Bytes"
	}

	return "Unknown"
}

// Query wraps an SQL query to expose its meta-data.
// Wrapping the query results allows the API to define higher level
// functions for interrogating the query metadata (such as the type).
type Query struct {
	executedQuery string
	result        *sql.Rows
	columns       []string
	types         []*sql.ColumnType
}

func logStartTime(t time.Time, i string) {
	f := log.Fields{
		"startTime": t,
	}

	if i != "" {
		f["requestID"] = i
	}

	log.WithFields(f).Info("Querying Database")
}

func logEndTime(t time.Time, i string) {
	f := log.Fields{
		"duration": time.Since(t),
		"finished": time.Now(),
	}

	if i != "" {
		f["requestID"] = i
	}

	log.WithFields(f).Info("Finished Querying Database")
}

func logError(e error, i string, fmt string, args ...interface{}) {
	t := time.Now()
	f := log.Fields{
		"time":    t.String(),
		"message": e.Error(),
	}

	if i != "" {
		f["requestID"] = i
	}

	log.WithFields(f).Errorf(fmt, args...)
}

// ExecuteQuery executes an SQL query and returns the wrapped results.
func ExecuteQuery(conn *sql.DB, query string, params ...interface{}) (Query, error) {
	startTime := time.Now()
	logStartTime(startTime, "")
	defer logEndTime(startTime, "")

	rows, err := conn.Query(query, params...)
	if err != nil {
		logError(err, "", "Failed to execute query with error: %v", err)
		return Query{}, err
	}

	cols, err := rows.Columns()
	if err != nil {
		logError(err, "", "Failed to get result columns: %v", err)
		rows.Close()
		return Query{}, err
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		logError(err, "", "Failed to get result column types: %v", err)
		rows.Close()
		return Query{}, err
	}

	return Query{
		executedQuery: query,
		result:        rows,
		columns:       cols,
		types:         types,
	}, nil
}

// ExecuteQueryWithContext executes a query with an available context for cancelation
func ExecuteQueryWithContext(ctx context.Context, conn *sql.DB, query string, params ...interface{}) (Query, error) {
	requestID := RequestIDFromContext(ctx)

	startTime := time.Now()
	logStartTime(startTime, requestID)

	defer logEndTime(startTime, requestID)

	// TODO: The database query needs to be "cancellable"
	rows, err := conn.QueryContext(ctx, query, params...)
	if err != nil {
		logError(err, requestID, "Failed to execute query with error: %v", err)
		return Query{}, err
	}

	cols, err := rows.Columns()
	if err != nil {
		logError(err, requestID, "Failed to get result columns: %v", err)
		rows.Close()
		return Query{}, err
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		logError(err, requestID, "Failed to get result column types: %v", err)
		rows.Close()
		return Query{}, err
	}

	return Query{
		executedQuery: query,
		result:        rows,
		columns:       cols,
		types:         types,
	}, nil
}

// Close allows the query to be closed from a defer
func (q Query) Close() error {
	return q.result.Close()
}

// Result returns the raw query result
func (q Query) Result() *sql.Rows {
	return q.result
}

// Columns returns the query result column names
func (q Query) Columns() []string {
	return q.columns
}

// IndexOf returns the index of a column with a given name
func (q Query) IndexOf(colName string) (int, error) {
	for i := range q.columns {
		if colName == q.columns[i] {
			return i, nil
		}
	}
	log.WithField("columnName", colName).Warnf("Failed to find column %s", colName)
	return -1, fmt.Errorf("Failed to find column %s", colName)
}

// Type returns the high-level type given either the index or the name
func (q Query) Type(column interface{}) (DBType, error) {
	name, ok := column.(string)
	idx := 0
	var err error

	if !ok {
		idx, ok = column.(int)
		if !ok {
			panic(fmt.Sprintf("Expected int or string but got %v", reflect.TypeOf(column)))
		}
	} else {
		idx, err = q.IndexOf(name)
	}

	if err != nil {
		return -1, fmt.Errorf("Unknown column %v, %d", column, idx)
	}

	//targetType := q.types[idx].DatabaseTypeName()
	//switch {
	//	case targetType == "VARCHAR" || targetType == "TEXT"
	//}
	return DBUnknown, nil
}

// Iterator is a function that can be passed into the query iterator
type Iterator interface {
	Iteration(row []interface{}) (bool, error)
}

// Accumulator is a function that accumulates a value
type Accumulator interface {
	Accumulate(current interface{}, row []interface{}) (interface{}, error)
}

// MakeBuffer creates a buffer that can process rows from this query
func (q Query) MakeBuffer() []interface{} {
	strings := make([]sql.NullString, len(q.columns))
	result := make([]interface{}, len(q.columns))

	for i, v := range strings {
		result[i] = &v
	}

	return result
}

// ForEach is a function that iterators over each result in the query
func (q Query) ForEach(i Iterator) error {
	startTime := time.Now()
	log.WithField("startTime", startTime).Info("Starting ForEach")
	defer log.WithFields(log.Fields{
		"duration": time.Since(startTime),
		"finished": time.Now(),
	}).Info("Finished ForEach")

	buffer := q.MakeBuffer()

	for q.result.Next() {
		log.Info("Starting next iteration")
		err := q.result.Scan(buffer...)

		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to scan query results: %v", err)
			return err
		}

		stop, err := i.Iteration(buffer)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Error processing result row using iterator: %v", err)
			return err
		}

		if stop {
			break
		}
	}

	return nil
}

// Accumulate executes the accumulator over the results
func (q Query) Accumulate(a Accumulator, starting interface{}) (interface{}, error) {
	startTime := time.Now()
	log.WithField("startTime", startTime).Info("Starting accumulator")
	defer log.WithFields(log.Fields{
		"duration": time.Since(startTime),
		"finished": time.Now(),
	}).Info("Finished accumulator")

	buffer := q.MakeBuffer()

	current := starting
	for q.result.Next() {
		log.Info("Starting next accumulator iteration")

		err := q.result.Scan(buffer...)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to scan query results: %v", err)
			return nil, err
		}

		current, err = a.Accumulate(current, buffer)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Error calling accumulator: %v", err)
			return starting, err
		}
	}

	return current, nil
}
