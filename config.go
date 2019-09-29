package wysci

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// DBConfig is a database connection configuration
type DBConfig struct {
	DBHost string `toml:"host,omitempty"`
	DBUser string `toml:"user,omitempty"`
	DBPass string `toml:"password,omitempty"`
	DBName string `toml:"database,omitempty"`
	DBPort int    `toml:"port,omitempty"`
}

// QueryConfig describes a query to execute
type QueryConfig struct {
	SQL    string `toml:"sql,omitempty"`
	Break  string `toml:"break,omitempty"`
	Params string `toml:"params,omitempty"`
}

// Service describes the service endpoint
type Service struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

// Parameter is an endpoint parameter
type Parameter struct {
	Type     string `toml:"type"`
	Required string `toml:"required"`
	Ordinal  int    `toml:"ordinal"`
}

// Endpoint describe a service endpoint
type Endpoint struct {
	QueryConfig string               `toml:"query"`
	Parameters  map[string]Parameter `toml:"parameters"`
}

// Configuration defines a wysci server
type Configuration struct {
	Database   DBConfig               `toml:"database"`
	Queries    map[string]QueryConfig `toml:"queries"`
	Connection Service                `toml:"connection"`
	Endpoints  map[string]Endpoint    `tomls:"endpoints"`
}

// LoadConfiguration loads the server
func LoadConfiguration(path string) (*Configuration, error) {
	log.Println("Reading TOML")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// config := make(map[string]interface{})
	config := Configuration{}
	_, err = toml.Decode(string(data), &config)
	if err != nil {
		fmt.Printf("Error parsing toml: %v\n", err)
	}

	return &config, nil
}
