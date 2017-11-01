// @todo: init indexes (unique, etc.)
package db

import (
	"gopkg.in/mgo.v2"
	"time"
)

// Params is the struct containing connection parameters to MongoDB
type DbParams struct {
	DbName string
	DbHost string
	DbUser string
	DbPass string
}

// Db is a struct embedding the MongoDB session and the connection params
type Db struct {
	Session *mgo.Session
	Params  DbParams
}

// NewSession clones the current session and returns it
func (db *Db) NewSession() *mgo.Session {
	return db.Session.Clone()
}

// Close closes the session
func (db *Db) Close() {
	db.Session.Close()
}

// CreateSession creates a new Db session, stores it inside
// a new Db struct, and returns the struct.
func CreateSession(dbParams DbParams) (*Db, error) {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{dbParams.DbHost},
		Database: dbParams.DbName,
		Timeout:  2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return &Db{session, dbParams}, nil
}
