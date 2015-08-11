package torque

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	// Registers the Postgres driver with the SQL package
	_ "github.com/lib/pq"
)

const postgresConfigFile = "pgconf.json"

var (
	// PsqlHost is PostgresQL hostname
	PsqlHost = flag.String("psql-host", "", "Hostname of Postgres DB")
	// PsqlPort is the PostgresQL port
	PsqlPort = flag.String("psql-port", "", "Port of Postgres DB")
	// PsqlUser is the Postgresql database user
	PsqlUser = flag.String("psql-user", "", "Postgresql user")
	// PsqlPassword is the Postgresql database user's password
	PsqlPassword = flag.String("psql-password", "", "Postgresql password")
	// PsqlDB is the Postgresql databse name
	PsqlDB = flag.String("psql-db", "", "Postgresql DB")
	// DBConn represents an open connection to a Postgres DB
	DBConn *sql.DB
)

// DBResource defines an object which implements basic data operations
type DBResource interface {
	DBCreate() error
	DBRetrieve() error
	DBUpdate() error
	DBDelete() error
}

// PostgresConfig is the minimal config needed to connect to a Postgres database.
// Shared DBConn singleton
type PostgresConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
	SSLMode  string `json:"sslmode"`
}

func (conf *PostgresConfig) buildPGURL() string {
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

// SafeString outptus the Postgresl URL in a security-safe manner
func (conf *PostgresConfig) SafeString() string {
	return fmt.Sprintf("postgres://%s@%s/%s", conf.User, conf.Host, conf.Database)
}

// GetDBConnection opens and returns a connection to the Postgresql DB
func GetDBConnection(conf *PostgresConfig) *sql.DB {
	DBConn, err := sql.Open("postgres", conf.buildPGURL())
	if err != nil {
		log.Fatalf("Can't connect to db: %s", err)
	}
	log.Printf("Database connection has been established; %s", conf.SafeString())
	return DBConn
}

// LoadPGConfig opens a PostgresConfig from a file
func LoadPGConfig() (conf *PostgresConfig) {
	f, err := os.Open(postgresConfigFile)
	if err == nil && !os.IsNotExist(err) {
		// Read values in from json
		decoder := json.NewDecoder(f)
		err = decoder.Decode(conf)
	}
	if err != nil {
		conf = &PostgresConfig{}
	}
	// Overwrite config file with command-line variables
	if *PsqlUser != "" {
		conf.User = *PsqlUser
	}
	if *PsqlHost != "" && *PsqlPort != "" {
		conf.Host = net.JoinHostPort(*PsqlHost, *PsqlPort)
	}
	if *PsqlPassword != "" {
		conf.Password = *PsqlPassword
	}
	if *PsqlDB != "" {
		conf.Database = *PsqlDB
	}
	return
}
