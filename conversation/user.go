package conversation

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id   bson.ObjectId `bson:"_id"`
	FbId string        `bson:"fbid"`
}
