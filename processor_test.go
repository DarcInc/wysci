package wysci

import (
	"bytes"
	"database/sql"
	"regexp"
	"testing"
)

func TestCSVFormatWriteRow(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: ",",
	}

	b := new(bytes.Buffer)
	bytesWritten, err := c.writeRow("foo,bar", b)
	if err != nil {
		t.Error(err)
	}
	if bytesWritten < len("foo,bar\r\n") {
		t.Errorf("Expected %d bytes written but got %d", len("foo,bar\r\n"), bytesWritten)
	}

}

func TestCSVFormatWriteHeader(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: ",",
		columns:   []string{"foo", "bar"},
	}

	b := new(bytes.Buffer)
	count, err := c.writeHeaders(b)
	if err != nil {
		t.Error(err)
	}
	if count < len("foo,bar") {
		t.Errorf("Expected length to be at least %d but got %d", len("foo,bar"), count)
	}

	if string(b.Bytes()) != "foo,bar\r\n" {
		t.Errorf("Invalid header row: %s", string(b.Bytes()))
	}

	if !c.didPrintHeaders {
		t.Errorf("Headers written but didPrintHeaders guard not set")
	}
}

func TestCSVFormat(t *testing.T) {
	if testConn == nil {
		t.Fatal("Test connection isn't set")
	}
	rows, err := testConn.Query("select * from test_simple where id = 1")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}

	buffer := make([]sql.NullString, len(cols))
	scanRow := make([]interface{}, len(buffer))
	for i := 0; i < len(cols); i++ {
		scanRow[i] = &buffer[i]
	}

	rows.Next()

	err = rows.Scan(scanRow...)
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)

	csv, err := NewCSVFormatter(rows)
	if err != nil {
		t.Fatal(err)
	}
	size, err := csv.Format(buffer, b)
	if size < 1 {
		t.Errorf("Invalid size: %d bytes", size)
	}

	re := regexp.MustCompile("hello world")
	if !re.Match(b.Bytes()) {
		t.Errorf("Did not contain required field")
	}
}

func TestProcessor(t *testing.T) {
	rows, err := testConn.Query("select * from test_simple where id = 1")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	output := new(bytes.Buffer)

	qp := QueryProcessor{}
	byteCount, err := qp.Process(rows, output)
	if err != nil {
		t.Error(err)
	}
	if byteCount < 1 {
		t.Errorf("Expected byte count > 1 but got %d", byteCount)
	}
}
