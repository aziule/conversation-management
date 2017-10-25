package mongo

import "gopkg.in/mgo.v2"

// Params is the struct containing connection parameters to MongoDB
type Params struct {
	DbName string
	DbHost string
	DbUser string
	DbPass string
}

// Db is a struct embedding the MongoDB session and the connection params
type Db struct {
	Session *mgo.Session
	Params  Params
}
