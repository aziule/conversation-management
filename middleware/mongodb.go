// @todo: remove - here only for inspiration
package middleware

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type Options struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

type MongoDB struct {
	currentDb *mgo.Database
	options   *Options
}

type MongoDBSession struct {
	*mgo.Session
	*Options
}

func New(options *Options) *MongoDB {
	return &MongoDB{
		options: options,
	}
}

func (db *MongoDB) NewSession() *MongoDBSession {
	mongoOptions := db.options

	// set default DialTimeout value
	if mongoOptions.DialTimeout <= 0 {
		mongoOptions.DialTimeout = 1 * time.Minute
	}

	session, err := mgo.DialWithTimeout(mongoOptions.ServerName, mongoOptions.DialTimeout)
	if err != nil {
		panic(err)
	}
	db.currentDb = session.DB(mongoOptions.DatabaseName)
	return &MongoDBSession{session, mongoOptions}
}

func (session *MongoDBSession) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	s := session.Clone()
	defer s.Close()

	db := &MongoDB{
		currentDb: s.DB(session.DatabaseName),
	}

	fmt.Println(db)

	next(w, req)
}
