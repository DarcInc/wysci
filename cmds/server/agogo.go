package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/darcinc/wysci"
	_ "github.com/lib/pq"
)

func finalString(env, cmd, config, def string) string {
	result := def
	if config != "" {
		result = config
	}

	if env != "" {
		result = env
	}

	if cmd != "" {
		result = cmd
	}

	return result
}

func finalInt(env, cmd, config, def int) int {
	result := def
	if config != 0 {
		result = config
	}

	if env != 0 {
		result = cmd
	}

	if cmd != 0 {
		result = cmd
	}

	return result
}

func main() {

	// TODO: Clean up the config code out of main
	cmdConf := flag.String("config", "", "The configuration to load")
	cmdHost := flag.String("dbhost", "", "The database hostname")
	cmdUser := flag.String("dbuser", "", "The database user")
	cmdName := flag.String("database", "", "The database to connect to")
	cmdPass := flag.String("dbpass", "", "The password for the database user")
	cmdPort := flag.Int("dbport", 5432, "The database listening port")

	flag.Parse()

	envConf := os.Getenv("CONFIG")
	envHost := os.Getenv("DBHOST")
	envUser := os.Getenv("DBUSER")
	envName := os.Getenv("DATABASE")
	envPass := os.Getenv("DBPASS")
	envPort, _ := strconv.Atoi(os.Getenv("DBPORT"))

	configFile := finalString(envConf, *cmdConf, "", "wysci.toml")

	config, err := wysci.LoadConfiguration(configFile)
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
	}

	databaseHost := finalString(envHost, *cmdHost, config.Database.DBHost, "localhost")
	databaseUser := finalString(envUser, *cmdUser, config.Database.DBUser, "postgres")
	databasePass := finalString(envPass, *cmdPass, config.Database.DBPass, "")
	databaseName := finalString(envName, *cmdName, config.Database.DBName, "postgres")
	databasePort := finalInt(envPort, *cmdPort, config.Database.DBPort, 5432)

	connString := fmt.Sprintf("user=%s dbname=%s host=%s password=%s port=%d sslmode=disable",
		databaseUser, databaseName, databaseHost, databasePass, databasePort)
	fmt.Println(connString)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		os.Exit(1)
	}

	router, err := wysci.ConfigureEndpoints(config, conn)
	if err != nil {
		log.Printf("Failed to create endpoints: %v", err)
	}

	// TODO: Needs to run using HTTPS
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Connection.Address, config.Connection.Port),
		router,
	)
}
