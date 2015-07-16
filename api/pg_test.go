package api

import "testing"

func TestBuildPGURL(t *testing.T) {
	pgConf := &pgConfig{
		User:     "Dude",
		Password: "Duuuuuuude",
		Database: "lebowski",
		Host:     "localhost:2345",
		SSLMode:  "verify-ca",
	}
	exp := "postgres://Dude:Duuuuuude@localhost:2345/lebowski?sslmode=verify-ca"
	output := pgConf.buildPGURL()
	if exp != output {
		t.Errorf("%s != %s", output, exp)
	}
}
