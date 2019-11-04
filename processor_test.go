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
	bytesWritten, err := c.writeRow([]string{"foo", "bar"}, b)
	if err != nil {
		t.Error(err)
	}
	if bytesWritten != len("foo,bar\r\n") {
		t.Errorf("Expected %d bytes written but got %d", len("foo,bar\r\n"), bytesWritten)
	}
}

func TestCSVFormatterEmbeddedDelimiters(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: ",",
	}

	b := new(bytes.Buffer)
	bytesWritten, err := c.writeRow([]string{"foo", "embedded,comma"}, b)
	if err != nil {
		t.Error(err)
	}

	if bytesWritten != len("foo,\"embedded,comma\"\r\n") {
		t.Errorf("Expected %d bytes written but got %d", len("foo,\"embedded,comma\"\r\n"), bytesWritten)
	}

	re := regexp.MustCompile("\"embedded,comma\"")
	if !re.Match(b.Bytes()) {
		t.Error("Expected expression to match embedded comma")
	}
}

func TestCSVFormatterEmbeddedDelimiters_Tab(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: "\t",
	}

	b := new(bytes.Buffer)
	bytesWritten, err := c.writeRow([]string{"foo", `embedded	tab`}, b)
	if err != nil {
		t.Error(err)
	}

	if bytesWritten != len("foo,\"embedded\ttab\"\r\n") {
		t.Errorf("Expected %d bytes written but got %d", len("foo,\"embedded\ttab\"\r\n"), bytesWritten)
	}

	re := regexp.MustCompile(`"embedded	tab"`)
	if !re.Match(b.Bytes()) {
		t.Error("Expected expression to match embedded tab")
	}
}

func TestCSVFormatterEmbeddedDelimiters_Bar(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: "|",
	}

	b := new(bytes.Buffer)
	bytesWritten, err := c.writeRow([]string{"foo", `embedded|bar`}, b)
	if err != nil {
		t.Error(err)
	}

	if bytesWritten != len("foo,\"embedded|bar\"\r\n") {
		t.Errorf("Expected %d bytes written but got %d", len("foo,\"embedded|bar\"\r\n"), bytesWritten)
	}

	re := regexp.MustCompile(`"embedded\|bar"`)
	if !re.Match(b.Bytes()) {
		t.Error("Expected expression to match embedded bar")
	}
}

func TestCSVFormatterEmbeddedQuotes(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: ",",
	}

	b := new(bytes.Buffer)
	bytesWritten, err := c.writeRow([]string{"foo", "\"embedded quotes\""}, b)
	if err != nil {
		t.Error(err)
	}

	if bytesWritten != len("foo,\"\"\"embedded quotes\"\"\"\r\n") {
		t.Errorf("Exepcted %d bytes written but got %d", len("foo,\"\"\"embedded quotes\"\"\"\r\n"), bytesWritten)
	}

	re := regexp.MustCompile("\"\"\"embedded quotes\"\"\"")
	if !re.Match(b.Bytes()) {
		t.Error("Expected expression to match embedded quotes")
	}
}

func TestCSVFormatWriteHeader(t *testing.T) {
	c := &CSVFormatter{
		Delimiter: ",",
		query: Query{
			columns: []string{"foo", "bar"},
		},
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
	query, err := ExecuteQuery(testConn, "select * from test_simple where id = 1")
	if err != nil {
		t.Fatal(err)
	}
	defer query.Close()

	cols := query.Columns()

	buffer := make([]sql.NullString, len(cols))
	scanRow := make([]interface{}, len(buffer))
	for i := 0; i < len(cols); i++ {
		scanRow[i] = &buffer[i]
	}

	query.Result().Next()

	err = query.Result().Scan(scanRow...)
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)

	csv, err := NewCSVFormatter(query)
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
	query, err := ExecuteQuery(testConn, "select * from test_simple where id = 1")
	if err != nil {
		t.Fatal(err)
	}
	defer query.Close()

	output := new(bytes.Buffer)

	qp := QueryProcessor{}
	byteCount, err := qp.Process(query, output)
	if err != nil {
		t.Error(err)
	}
	if byteCount < 1 {
		t.Errorf("Expected byte count > 1 but got %d", byteCount)
	}
}

func TestProcessorWithEmbeddedDelimiters(t *testing.T) {
	query, err := ExecuteQuery(testConn, "select * from test_simple where id = 4")
	if err != nil {
		t.Fatal(err)
	}
	defer query.Close()

	output := new(bytes.Buffer)

	qp := QueryProcessor{}
	_, err = qp.Process(query, output)
	if err != nil {
		t.Error(err)
	}

	re := regexp.MustCompile("\"embedded,comma\"")
	if !re.Match(output.Bytes()) {
		t.Error("Failed to match double quotes around embedded delimiter")
	}
}
