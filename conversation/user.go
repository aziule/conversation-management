package conversation

import "gopkg.in/mgo.v2/bson"

// User is the main user model shared across the different platforms
type User struct {
	Id   bson.ObjectId `bson:"_id"`
	FbId string        `bson:"fbid"`
}
