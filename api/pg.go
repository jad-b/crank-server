package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	// Registers the Postgres driver with the SQL package
	_ "github.com/lib/pq"
)

const postgresConfigFile = "pgconf.json"

// Shared PGConn singleton
var PGConn *sql.DB

//PGConfig is the minimal config needed to connect to a Postgres database.
type pgConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
	SSLMode  string `json:"sslmode"`
}

func (conf *pgConfig) buildPGURL() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   conf.Host,
		Path:   conf.Database,
	}
	q := u.Query()
	q.Set("sslmode", conf.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func (conf *pgConfig) SafeString() string {
	return fmt.Sprintf("postgres://%s@%s/%s", conf.User, conf.Host, conf.Database)
}

// TODO(jdb) Convert to GetDBConnection
func init() {
	pgConf := loadPGConfig()
	PGConn, err := sql.Open("postgres", pgConf.buildPGURL())
	if err != nil {
		log.Fatalf("Can't connect to db: %s", err)
	}
	log.Printf("Database connection has been established; %s", pgConf.SafeString())
}

func loadPGConfig() (conf *pgConfig) {
	f, err := os.Open(postgresConfigFile)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Failed to load Postgres config from %s\nError:\n\t%s", postgresConfigFile, err)
	}
	// Read values in from json
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(conf); err != nil {
		log.Fatal(err)
	}
	return
}
