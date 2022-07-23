package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	_ "github.com/lib/pq"
)

type Config struct {
	Hostname string
	Port     string
	User     string
	Password string
	DBName   string
	SSL      string //disable/require
	//SchemaName string
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
		//SchemaName: schema,
	}
}

func GetInstance(config *Config) *Singleton {
	once.Do(func() {
		db, err := ConnectToDb(config)
		if err != nil {
			// TODO - EXIT!
		}
		Instance = &Singleton{
			Db: db,
		}
	})
	return Instance
}

// ConnectToDb creates a connection to the DB and returns it TODO
func ConnectToDb(config *Config) (*sql.DB, error) {
	//schema := ""
	//if config.SchemaName != "" {
	//	schema = fmt.Sprintf("search_path=%s&", config.SchemaName)
	//}
	// use the url pattern to escape special characters in username or password
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%ssslmode=%s", url.QueryEscape(config.User),
		url.QueryEscape(config.Password), config.Hostname, config.Port, config.DBName, "", config.SSL)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return db, err
	}
	//db.SetConnMaxLifetime(30 * time.Second)
	return db, nil
}
