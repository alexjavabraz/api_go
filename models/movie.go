package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Card struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	IndicePk    string        `bson:"indicePk" json:"indicePk"`
	CoverImage  string        `bson:"indiceSk" json:"indiceSk"`
	Description string        `bson:"accountId" json:"accountId"`
}
