package wysci

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func typeMetadata(rows *sql.Rows) ([]string, []string, error) {
	names, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	typenames := make([]string, len(types))
	for i := range types {
		typenames[i] = types[i].DatabaseTypeName()
	}

	return names, typenames, nil
}

func createBuffer(typenames []string) []interface{} {
	buffer := make([]interface{}, len(typenames))
	// "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL", "INT", "BIGINT"
	for i, t := range typenames {
		switch {
		case t == "VARCHAR" || t == "TEXT" || t == "NVARCHAR":
			buffer[i] = new(string)
		case t == "DECIMAL":
			buffer[i] = new(float64)
		case t == "BOOL":
			buffer[i] = new(bool)
		case t == "INT" || t == "INT4":
			buffer[i] = new(int32)
		case t == "BIGINT":
			buffer[i] = new(int64)
		}
	}
	return buffer
}

func castString(x interface{}) (string, error) {
	strptr, ok := x.(*string)
	if !ok {
		return "", fmt.Errorf("Expected string ptr but got %v", reflect.TypeOf(x))
	}

	return *strptr, nil
}

func castInt32(x interface{}) (int32, error) {
	int32ptr, ok := x.(*int32)
	if !ok {
		return 0, fmt.Errorf("Expected int ptr but got %v", reflect.TypeOf(x))
	}
	return *int32ptr, nil
}

func scanRow(rows *sql.Rows, buffer []interface{}, typenames []string) []string {
	rows.Scan(buffer...)
	parts := []string{}

	for idx, v := range buffer {
		t := typenames[idx]
		switch {
		case t == "VARCHAR" || t == "TEXT" || t == "NVCHAR":
			s, err := castString(v)
			if err != nil {
				log.Printf("Unable to cast to string: %v", err)
			}
			parts = append(parts, fmt.Sprintf("%v", s))
		case t == "INT4":
			i, err := castInt32(v)
			if err != nil {
				log.Printf("Unable to cast to int32: %v", err)
			}
			parts = append(parts, fmt.Sprintf("%d", i))
		}
	}

	return parts
}

func appendCSVLine(w http.ResponseWriter, row []string) {
	w.Write([]byte(strings.Join(row, ",")))
	w.Write([]byte("\r\n"))
}

// ConfigureEndpoints configures the service endpoints
func ConfigureEndpoints(config *Configuration, conn *sql.DB) (*httprouter.Router, error) {
	router := httprouter.New()

	for name, endpoint := range config.Endpoints {
		query := config.Queries[endpoint.QueryConfig]

		router.GET(fmt.Sprintf("/api/v1/%s", name), func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			result, err := conn.Query(query.SQL)
			if err != nil {
				log.Printf("Failed to execute query: %v", err)
				w.WriteHeader(500)
				return
			}
			defer result.Close()

			columns, typenames, err := typeMetadata(result)
			if err != nil {
				log.Printf("Failed to get result metadata: %v", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Add("Content-Type", "text/csv")
			appendCSVLine(w, columns)
			values := createBuffer(typenames)
			for result.Next() {
				result.Scan(values...)
				parts := scanRow(result, values, typenames)
				appendCSVLine(w, parts)
			}
		})
	}

	return router, nil
}
