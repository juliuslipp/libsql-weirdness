package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/libsql/go-libsql"
	_ "modernc.org/sqlite"
	"path/filepath"
	"time"
)

type Client = sqlx.DB

type ClientConfig struct {
	SyncInterval time.Duration
	JournalMode  bool

	DatabaseUrl       string
	DatabaseAuthToken *string
}

func NewClientConfig(databaseUrl string, databaseAuthToken *string) ClientConfig {
	return ClientConfig{
		SyncInterval:      1 * time.Second,
		DatabaseUrl:       databaseUrl,
		DatabaseAuthToken: databaseAuthToken,
	}
}

func NewClient(config ClientConfig) (*Client, error) {
	dbUrl := config.DatabaseUrl
	if config.DatabaseAuthToken != nil {
		dbUrl = dbUrl + "?authToken=" + *config.DatabaseAuthToken
	}

	dbConn, err := sqlx.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func NewEmbeddedClient(config ClientConfig) (*Client, error) {
	if true {
		return NewClient(config)
	}

	absPath, _ := filepath.Abs("./sqlite.db")
	connector, err := libsql.NewEmbeddedReplicaConnectorWithAutoSync(absPath, config.DatabaseUrl, *config.DatabaseAuthToken, config.SyncInterval)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(sql.OpenDB(connector), "sqlite")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
