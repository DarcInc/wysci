package wysci

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

// Create a new buffer by allocating
// TODO: need a more complete list of Postgres types (this may be okay for now)
func createBuffer(typenames []string) []interface{} {
	buffer := make([]interface{}, len(typenames))
	// "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL", "INT", "BIGINT"
	for i, t := range typenames {
		switch {
		case t == "VARCHAR" || t == "TEXT" || t == "NVARCHAR" || t == "MONEY":
			buffer[i] = new(sql.NullString)
		case t == "TIMESTAMP" || t == "DATE":
			buffer[i] = new(sql.NullTime)
		case t == "DECIMAL":
			buffer[i] = new(sql.NullFloat64)
		case t == "BOOL":
			buffer[i] = new(sql.NullBool)
		case t == "INT" || t == "INT4":
			buffer[i] = new(sql.NullInt32)
		case t == "BIGINT":
			buffer[i] = new(sql.NullInt64)
		}
	}
	return buffer
}

// Safely cast out a string
func castString(x interface{}) (*sql.NullString, error) {
	strptr, ok := x.(*sql.NullString)
	if !ok {
		return nil, fmt.Errorf("Expected string ptr but got %v", reflect.TypeOf(x))
	}

	return strptr, nil
}

// Safely cast out an int32
func castInt32(x interface{}) (*sql.NullInt32, error) {
	int32ptr, ok := x.(*sql.NullInt32)
	if !ok {
		return nil, fmt.Errorf("Expected int ptr but got %v", reflect.TypeOf(x))
	}

	return int32ptr, nil
}

func castInt64(x interface{}) (*sql.NullInt64, error) {
	int64ptr, ok := x.(*sql.NullInt64)
	if !ok {
		return nil, fmt.Errorf("Expected int64 ptr but got %v", reflect.TypeOf(x))
	}
	return int64ptr, nil
}

func castFloat64(x interface{}) (*sql.NullFloat64, error) {
	float32Ptr, ok := x.(*sql.NullFloat64)
	if !ok {
		return nil, fmt.Errorf("Expected float64 ptr but got %v", reflect.TypeOf(x))
	}
	return float32Ptr, nil
}

func castBool(x interface{}) (*sql.NullBool, error) {
	boolPtr, ok := x.(*sql.NullBool)
	if !ok {
		return nil, fmt.Errorf("Expected bool ptr but got %v", reflect.TypeOf(x))
	}
	return boolPtr, nil
}

func castTime(x interface{}) (*sql.NullTime, error) {
	timePtr, ok := x.(*sql.NullTime)
	if !ok {
		return nil, fmt.Errorf("Expected time ptr but got %v", reflect.TypeOf(x))
	}
	return timePtr, nil
}

// Scan a row into a buffer
func scanRow(rows *sql.Rows, buffer []interface{}, typenames []string) []string {
	rows.Scan(buffer...)
	parts := []string{}

	for idx, v := range buffer {
		t := typenames[idx]
		switch {
		case t == "VARCHAR" || t == "TEXT" || t == "NVCHAR" || t == "MONEY":
			s, err := castString(v)
			if err != nil {
				log.Printf("Unable to cast to string: %v", err)
				continue
			}

			if s.Valid {
				parts = append(parts, s.String)
			} else {
				parts = append(parts, "")
			}
		case t == "TIMESTAMP" || t == "DATE":
			s, err := castTime(v)
			if err != nil {
				log.Printf("Unable to cast time: %v", err)
				continue
			}

			if s.Valid {
				switch {
				case t == "TIMESTAMP":
					parts = append(parts, fmt.Sprintf("%v", s.Time))
				case t == "DATE":
					parts = append(parts, s.Time.Format("01/02/2006"))
				}
			} else {
				parts = append(parts, "")
			}
		case t == "DECIMAL":
			f, err := castFloat64(v)
			if err != nil {
				log.Printf("unable to cast to float: %v", err)
				continue
			}

			if f.Valid {
				parts = append(parts, fmt.Sprintf("%f", f.Float64))
			} else {
				parts = append(parts, "")
			}
		case t == "BOOL":
			b, err := castBool(v)
			if err != nil {
				log.Printf("Unable to cast to bool: %v", err)
				continue
			}

			if b.Valid {
				parts = append(parts, fmt.Sprintf("%v", b.Bool))
			} else {
				parts = append(parts, "")
			}
		case t == "INT4":
			i, err := castInt32(v)
			if err != nil {
				log.Printf("Unable to cast to int32: %v", err)
				continue
			}

			if i.Valid {
				parts = append(parts, fmt.Sprintf("%d", i.Int32))
			} else {
				parts = append(parts, "")
			}
		case t == "BIGINT":
			b, err := castInt64(v)
			if err != nil {
				log.Printf("Unable to cast to int64: %v", err)
				continue
			}

			if b.Valid {
				parts = append(parts, fmt.Sprintf("%d", b.Int64))
			} else {
				parts = append(parts, "")
			}
		}
	}

	return parts
}

// Appends a line to the CSV response
func appendCSVLine(w http.ResponseWriter, row []string) {
	w.Write([]byte(strings.Join(row, ",")))
	w.Write([]byte("\r\n"))
}

// ConfigureEndpoints configures the service endpoints
func ConfigureEndpoints(config *Configuration, conn *sql.DB) (*httprouter.Router, error) {
	router := httprouter.New()

	for name, endpoint := range config.Endpoints {
		query := config.Queries[endpoint.QueryConfig]

		log.Printf("Adding /api/v1/%s", name)
		router.GET(fmt.Sprintf("/api/v1/%s", name), func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			parameters := make([]interface{}, len(endpoint.Parameters))
			for pname, pdesc := range endpoint.Parameters {
				raw := r.URL.Query().Get(pname)
				switch {
				case pdesc.Type == "number":
					val, err := strconv.ParseInt(raw, 10, 64)
					if err != nil {
						log.Printf("Failed to configure %s: %v", name, err)
						continue
					}
					parameters[pdesc.Ordinal-1] = val
				case pdesc.Type == "string":
					parameters[pdesc.Ordinal-1] = raw
				}
			}

			result, err := conn.Query(query.SQL, parameters...)
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

			log.Printf("%v", typenames)

			//w.Header().Add("Content-Type", "text/csv")
			appendCSVLine(w, columns)
			values := createBuffer(typenames)
			for result.Next() {
				parts := scanRow(result, values, typenames)
				appendCSVLine(w, parts)
			}
		})
	}

	return router, nil
}
