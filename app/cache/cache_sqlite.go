// +build !appengine,!bolt,sqlite

package cache

import (
	"database/sql"
	"net/http"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./cache.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ID (key TEXT NOT NULL PRIMARY KEY, value TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}
}

type Cache struct {
	mu sync.RWMutex
}

type RCache struct {
	*Cache
	r *http.Request
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) New(r *http.Request) *RCache {
	return &RCache{Cache: c, r: r}
}

func (c *RCache) Set(id string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newID, err := translate(id)
	if err != nil {
		return "", err
	}

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO ID(key, value) VALUES(?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newID, id)
	if err != nil {
		return "", err
	}

	tx.Commit()

	return newID, nil
}

func (c *RCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stmt, err := db.Prepare("SELECT value FROM ID WHERE key = ?")
	if err != nil {
		return "", false
	}
	defer stmt.Close()

	var oldID string
	err = stmt.QueryRow(id).Scan(&oldID)
	if err != nil {
		return "", false
	}

	return oldID, true
}
