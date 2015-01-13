package torque

import "testing"

func TestBuildPGURL(t *testing.T) {
	pgConf := &PostgresConfig{
		User:     "Dude",
		Password: "Duuuuuuude",
		Database: "lebowski",
		Host:     "localhost:2345",
		SSLMode:  "verify-ca",
	}
	exp := "postgres://Dude:Duuuuuuude@localhost:2345/lebowski?sslmode=verify-ca"
	output := pgConf.buildPGURL()
	if exp != output {
		t.Errorf("\n%s != \n%s", output, exp)
	}
}

func TestCreateTable(t *testing.T) {

}
