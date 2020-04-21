package models

import "gopkg.in/mgo.v2/bson"

// Advertisment model
type Advertisment struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Advertiser  bson.ObjectId `json:"advertiser" bson:"advertiser"`
	CompanyName string        `json:"companyName" bson:"companyName"`
	Website     string        `json:"website" bson:"website"`
}
