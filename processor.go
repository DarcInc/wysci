package wysci

import (
	"database/sql"
	"io"
	"strings"
)

var safeEOL = []byte("\r\n")

// RowFormatter is implemented by types that format responses.
type RowFormatter interface {
	Format(values []sql.NullString, w io.Writer) (int, error)
}

// ColumnCounter is implemented by types that return the number of output columns
type ColumnCounter interface {
	ColumnCount() int
}

// Formatter is responsible for formatting database results
type Formatter interface {
	RowFormatter
	ColumnCounter
}

// CSVFormatter implements the Formatter interface to format CSV output
type CSVFormatter struct {
	Delimiter, NullString string
	columns               []string
	types                 []*sql.ColumnType
	didPrintHeaders       bool
}

// NewCSVFormatter creates a new CSV Formatter
func NewCSVFormatter(rows *sql.Rows) (*CSVFormatter, error) {
	formatter := &CSVFormatter{}
	formatter.Delimiter = ","
	formatter.NullString = ""

	var err error

	formatter.columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	formatter.types, err = rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	return formatter, nil
}

func (c *CSVFormatter) writeRow(row string, w io.Writer) (int, error) {
	bytesWritten, err := w.Write([]byte(row))
	if err != nil {
		return 0, err
	}

	nextBytes, err := w.Write(safeEOL)
	bytesWritten += nextBytes
	if err != nil {
		return bytesWritten, err
	}

	return bytesWritten, nil
}

func (c *CSVFormatter) writeHeaders(w io.Writer) (int, error) {
	bytesWritten, err := c.writeRow(strings.Join(c.columns, c.Delimiter), w)
	if err != nil {
		return bytesWritten, err
	}

	c.didPrintHeaders = true
	return bytesWritten, nil
}

// Format formats the CSV response
func (c *CSVFormatter) Format(values []sql.NullString, w io.Writer) (int, error) {
	var bytesWritten int
	var err error

	if !c.didPrintHeaders {
		bytesWritten, err = c.writeHeaders(w)
		if err != nil {
			return bytesWritten, err
		}
	}

	row := make([]string, len(c.columns))
	for i, v := range values {
		if v.Valid {
			row[i] = v.String
		} else {
			row[i] = c.NullString
		}
	}

	nextBytes, err := c.writeRow(strings.Join(row, c.Delimiter), w)
	bytesWritten += nextBytes
	if err != nil {
		return bytesWritten, err
	}

	return bytesWritten, nil
}

// ColumnCount returns the columns in the CSV formatter
func (c *CSVFormatter) ColumnCount() int {
	return len(c.columns)
}

// QueryProcessor translates results and passes them to a formatter.
type QueryProcessor struct {
	RowFormatter Formatter
}

// Process processes the results to pass to the formatter
func (qp QueryProcessor) Process(rows *sql.Rows, w io.Writer) (int, error) {
	var err error

	if qp.RowFormatter == nil {
		qp.RowFormatter, err = NewCSVFormatter(rows)
		if err != nil {
			return 0, err
		}
	}

	buffer := make([]sql.NullString, qp.RowFormatter.ColumnCount())
	scanLine := make([]interface{}, qp.RowFormatter.ColumnCount())
	for i := 0; i < len(buffer); i++ {
		scanLine[i] = &buffer[i]
	}

	totalBytes := 0
	for rows.Next() {
		err = rows.Scan(scanLine...)
		if err != nil {
			return totalBytes, err
		}

		bytesFormatted, err := qp.RowFormatter.Format(buffer, w)
		totalBytes += bytesFormatted
		if err != nil {
			return totalBytes, err
		}
	}

	return totalBytes, nil
}
