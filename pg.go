package torque

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
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
	PsqlDB = flag.String("psql-db", "", "Postgresqlx.DB")
	// PsqlConf is a filepath to a configuration file
	PsqlConf = flag.String("psql-conf", "pgconf.json", "Configuration file for DB connection")
	// DB represents an open connection to a Postgres DB
	DB *sqlx.DB
)

// DBActor defines an object which implements basic data operations
type DBActor interface {
	Create(*sqlx.DB) error
	Retrieve(*sqlx.DB) error
	Update(*sqlx.DB) error
	Delete(*sqlx.DB) error
}

// PostgresConfig is the minimal config needed to connect to a Postgres database.
// Shared DB singleton
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

// OpenDBConnection opens and returns a connection to the Postgresqlx.DB
func OpenDBConnection(conf *PostgresConfig) *sqlx.DB {
	db := sqlx.MustConnect("postgres", conf.buildPGURL())
	log.Printf("Database connection has been established; %s", conf.SafeString())
	// Assign to global connection
	DB = db
	return DB
}

// LoadPostgresConfig opens a PostgresConfig from a file
func LoadPostgresConfig(confFile string) (conf *PostgresConfig) {
	conf = &PostgresConfig{} // Use a blank configuration
	absConfPath, err := filepath.Abs(confFile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(absConfPath)
	if err != nil || os.IsNotExist(err) {
		log.Printf("No database configuration file found at %s", absConfPath)
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
