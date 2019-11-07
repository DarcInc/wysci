package wysci

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	simpleTable = `create table if not exists test_simple (
		id int,
		name varchar,
		some_date date
	);`

	basicTypesTable = `create table if not exists test_basic_types (
		sample_int int,
		sample_large_int bigint,
		sample_numeric0 numeric(10,0),
		sample_numeric2 numeric(10,2),
		sample_float float,
		sample_money money,
		sample_fixed_char char(3),
		sample_varchar30 varchar(30),
		sample_varchar varchar,
		sample_text text
	);`

	dateTimeTables = `create table if not exists date_time_types (
		sample_date date,
		sample_time time,
		sample_datetime timestamp
	);`

	addSimpleData = `insert into test_simple values (1, 'hello world', '2019-01-01');
		insert into test_simple values (2, NULL, '2019-01-02');
		insert into test_simple values (3, 'date is null', NULL);
		insert into test_simple values (4, 'embedded,comma', '2019-01-04');
		insert into test_simple values (5, 'embedded	tab', '2019-01-05');`

	addSampleData = `insert into test_basic_types values (1, 1, 1.0, 1.0, 1.0, '$1.00', 'hel', 'Hello world', 'hello world', 'hello world');
		insert into test_basic_types values (null, null, null, null, null, null, null, 'foo,bar', '"Hello World"', '"Hello, World"');
		insert into date_time_types values (current_date, current_time, current_timestamp);
		insert into date_time_types values (null, null, null);`

	cleanupTables = `drop table if exists test_basic_types;
		drop table if exists date_time_types;
		drop table if exists test_simple;`
)

var testConn *sql.DB

func setup(conn *sql.DB) error {
	tables := []string{
		simpleTable,
		basicTypesTable,
		dateTimeTables,
	}

	datas := []string{
		addSimpleData,
		addSampleData,
	}

	for _, v := range tables {
		_, err := conn.Exec(v)
		if err != nil {
			log.Printf("Failed to set up table: \"%s\"", v)
			return err
		}
	}

	for _, v := range datas {
		_, err := conn.Exec(v)
		if err != nil {
			log.Printf("Failed to set up data: \"%s\"", v)
			return err
		}
	}

	return nil
}

func teardown(conn *sql.DB) error {
	_, err := conn.Exec(cleanupTables)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	connStr := os.Getenv("CONNSTR")
	var err error

	log.SetFormatter(&log.JSONFormatter{})

	testConn, err = sql.Open("sqlite3", connStr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = setup(testConn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	err = teardown(testConn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}
