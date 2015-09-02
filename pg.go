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
	"time"

	// Registers the Postgres driver with the SQL package
	_ "github.com/lib/pq"
)

// PostgresqlTimestampFormat is a Postgresql-accepted timestamp layout
const PostgresqlTimestampFormat = "2006-01-02 15:04:05 MST"

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
	// PsqlConf is a filepath to a configuration file
	PsqlConf = flag.String("psql-conf", "pgconf.json", "Configuration file for DB connection")
	// DBConn represents an open connection to a Postgres DB
	DBConn *sql.DB
)

// DBActor defines an object which implements basic data operations
type DBActor interface {
	Create(*sql.DB) error
	Retrieve(*sql.DB) error
	Update(*sql.DB) error
	Delete(*sql.DB) error
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

// OpenDBConnection opens and returns a connection to the Postgresql DB
func OpenDBConnection(conf *PostgresConfig) *sql.DB {
	conn, err := sql.Open("postgres", conf.buildPGURL())
	if err != nil {
		log.Fatalf("Can't connect to db: %s", err)
	}
	log.Printf("Database connection has been established; %s", conf.SafeString())
	// Assign to global connection
	DBConn = conn
	return DBConn
}

// LoadPostgresConfig opens a PostgresConfig from a file
func LoadPostgresConfig() (conf *PostgresConfig) {
	conf = &PostgresConfig{} // Use a blank configuration
	f, err := os.Open(*PsqlConf)
	if err != nil || os.IsNotExist(err) {
		log.Fatalf("No database configuration file found at %s", *PsqlConf)
	} else {
		err = json.NewDecoder(f).Decode(conf)
		if err != nil {
			// Let's not kid our users and act like the file was OK
			log.Fatalf("Failed to read %s; %s", *PsqlConf, err)
		}
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
	return conf
}

// ToPsqlTimestamp formats a time.Time into a Postgresql-acceptable string.
func ToPsqlTimestamp(ts time.Time) string {
	return ts.Format(PostgresqlTimestampFormat)
}
