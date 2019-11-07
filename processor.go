package wysci

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

var safeEOL = []byte("\r\n")

// The Formatter interface is implemented by types that format responses.
// The Format function expects a slice of nullable strings that have been
// populated by a previous call to Scan.  It will then write those values to
// the io.Writer and return the count of the bytes written and any error.
type Formatter interface {
	Format(values []sql.NullString, w io.Writer) (int, error)
}

// CSVFormatter implements the Formatter interface to format CSV output.
// It outputs delimited format database columns.  The delimiter and the
// value for NULL strings can be customized by setting the respective
// fields.
type CSVFormatter struct {
	Delimiter, NullString string
	query                 Query
	didPrintHeaders       bool
}

// NewCSVFormatter creates a new CSV Formatter.
// The delimiter is defaulted to a comma and the NullString is defaulted to
// the empty string.  The rows parameter is the output of a database query.
func NewCSVFormatter(q Query) (*CSVFormatter, error) {
	formatter := &CSVFormatter{}
	formatter.Delimiter = ","
	formatter.NullString = ""
	formatter.query = q

	return formatter, nil
}

func (c *CSVFormatter) writeRow(row []string, w io.Writer) (int, error) {
	re, err := regexp.Compile(fmt.Sprintf("[\"%s]", c.Delimiter))
	if err != nil {
		log.WithField("message", err.Error()).Errorf("Failed to compile embedded quotes expression: %v", err)
		return 0, err
	}

	replQuotes, err := regexp.Compile("\"")
	if err != nil {
		log.WithField("message", err.Error()).Errorf("Failed to compile quote replacement expression: %v", err)
		return 0, err
	}

	doubleQuotes := []byte("\"\"")
	b := new(bytes.Buffer)

	for i := 0; i < len(row); i++ {
		byteRow := []byte(row[i])
		if re.Match(byteRow) {
			byteRow = replQuotes.ReplaceAll(byteRow, doubleQuotes)
			byteRow = append(append([]byte{'"'}, byteRow...), '"')
		}

		_, err := b.Write(byteRow)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to write buffer: %v", err)
			return 0, err
		}

		if i < len(row)-1 {
			_, err = b.Write([]byte(c.Delimiter))
			if err != nil {
				log.WithField("message", err.Error()).Errorf("Failed to write buffer: %v", err)
				return 0, err
			}
		}
	}

	_, err = b.Write(safeEOL)
	if err != nil {
		log.WithField("message", err.Error()).Errorf("Failed to write buffer: %v", err)
		return 0, err
	}

	n, err := w.Write(b.Bytes())
	if err != nil {
		log.WithField("message", err.Error()).Errorf("Failed to write buffer: %v", err)
		return n, err
	}

	return n, nil
}

func (c *CSVFormatter) writeHeaders(w io.Writer) (int, error) {
	bytesWritten, err := c.writeRow(c.query.columns, w)
	if err != nil {
		return bytesWritten, err
	}

	c.didPrintHeaders = true
	return bytesWritten, nil
}

// Format formats the CSV response.
// It implements the RowFormatter interface for the CSVFormatter type.
func (c *CSVFormatter) Format(values []sql.NullString, w io.Writer) (int, error) {
	var bytesWritten int
	var err error

	if !c.didPrintHeaders {
		bytesWritten, err = c.writeHeaders(w)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to write headers: %v", err)
			return bytesWritten, err
		}
	}

	row := make([]string, len(c.query.columns))
	for i, v := range values {
		if v.Valid {
			row[i] = v.String
		} else {
			row[i] = c.NullString
		}
	}

	nextBytes, err := c.writeRow(row, w)
	bytesWritten += nextBytes
	if err != nil {
		log.WithField("message", err.Error()).Errorf("Failed to write response row: %v", err)
		return bytesWritten, err
	}

	return bytesWritten, nil
}

// ColumnCount returns the columns in the CSV formatter.
// It implements the ColumnCounter interface for the CSVFormatter type.
func (c *CSVFormatter) ColumnCount() int {
	return len(c.query.columns)
}

// QueryProcessor translates results and passes them to a formatter.
// The QueryProcessor is responsible for extracting the database query
// and passing the results to the formatter for output.
type QueryProcessor struct {
	RowFormatter Formatter
}

// Process processes the results to pass to the formatter.
// It is responsible for Scanning the resulting rows and then passing
// those rows to the formatter for final output.
func (qp QueryProcessor) Process(query Query, w io.Writer) (int, error) {
	var err error
	startTime := time.Now()

	log.WithField("startTime", startTime).Info("Start processing query")
	defer log.WithFields(log.Fields{
		"finished": time.Now(),
		"duration": time.Since(startTime),
	}).Info("Finished processing query")

	if qp.RowFormatter == nil {
		qp.RowFormatter, err = NewCSVFormatter(query)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to generate new CSV formatter: %v", err)
			return 0, err
		}
	}

	nCols := len(query.Columns())
	buffer := make([]sql.NullString, nCols)
	scanLine := make([]interface{}, nCols)
	for i := 0; i < len(buffer); i++ {
		scanLine[i] = &buffer[i]
	}

	totalBytes := 0
	for query.result.Next() {
		err = query.result.Scan(scanLine...)
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to scan result row: %v", err)
			return totalBytes, err
		}

		bytesFormatted, err := qp.RowFormatter.Format(buffer, w)
		totalBytes += bytesFormatted
		if err != nil {
			log.WithField("message", err.Error()).Errorf("Failed to format result row: %v", err)
			return totalBytes, err
		}
	}

	return totalBytes, nil
}
