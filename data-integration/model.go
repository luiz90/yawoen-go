package main

import "gopkg.in/mgo.v2/bson"

// Company Entity
type Company struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string        `json:"name"`
	AddressZip string        `json:"zip"`
	Website    string        `json:"website"`
}

// Companies ...
type Companies []Company
