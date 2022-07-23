package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Hostname string
	Port     string
	User     string
	Password string
	DBName   string
	SSL      string
}

type Singleton struct {
	Db *sql.DB
}

var (
	once     sync.Once
	Instance *Singleton
)

func NewConfig(host, port, user, pass, name, ssl string) *Config {
	return &Config{
		Hostname: host,
		Port:     port,
		User:     user,
		Password: pass,
		DBName:   name,
		SSL:      ssl,
	}
}

func GetInstance(config *Config) *Singleton {
	once.Do(func() {
		db, err := ConnectToDb(config)
		if err != nil {
			fmt.Printf("ConnectToDb Error: %v\n", err)
			os.Exit(0)
		}
		Instance = &Singleton{
			Db: db,
		}
	})
	return Instance
}

// ConnectToDb creates a connection to the DB and returns it TODO
func ConnectToDb(config *Config) (*sql.DB, error) {
	// Uses the url pattern to escape special characters in username or password
	// POSTGRESQL_URL='postgres://user:password@host:port/db_name?sslmode=disable/require'
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", url.QueryEscape(config.User),
		url.QueryEscape(config.Password), config.Hostname, config.Port, config.DBName, config.SSL)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return db, err
	}
	
	time.Sleep(60 * time.Second)
	if err = db.Ping(); err != nil {
		return db, err
	}
	db.SetConnMaxLifetime(30 * time.Second)
	return db, nil
}
