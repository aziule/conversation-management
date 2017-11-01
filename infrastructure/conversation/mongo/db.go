// @todo: move to a dedicated package outside conversation
// @todo: init indexes (unique, etc.)
package mongo

import "gopkg.in/mgo.v2"

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
