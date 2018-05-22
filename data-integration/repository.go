package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "yawoen"

// DOCNAME the name of the document
const DOCNAME = "companies"

// GetCompanies returns the list of Companies
func (r Repository) GetCompanies() Companies {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(DOCNAME)
	results := Companies{}
	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}
	return results
}

// AddCompany inserts a Company in the DB
func (r Repository) AddCompany(company Company) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	company.ID = bson.NewObjectId()
	session.DB(DBNAME).C(DOCNAME).Insert(company)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// AddCompanies inserts companies in the DB
func (r Repository) AddCompanies(companies Companies) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	for _, company := range companies {
		company.ID = bson.NewObjectId()
		session.DB(DBNAME).C(DOCNAME).Insert(company)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	return true
}

// UpdateCompany updates a Company in the DB
func (r Repository) UpdateCompany(company Company) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	session.DB(DBNAME).C(DOCNAME).UpdateId(company.ID, company)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// DeleteCompany deletes a Company (not used for now)
func (r Repository) DeleteCompany(id string) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return "NOT FOUND"
	}
	// Grab id
	oid := bson.ObjectIdHex(id)
	// Remove user
	if err = session.DB(DBNAME).C(DOCNAME).RemoveId(oid); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}
	// Write status
	return "OK"
}

// FindByName ...
func (r Repository) FindByName(name string) Company {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	results := Company{}
	c := session.DB(DBNAME).C(DOCNAME)
	c.Find(bson.M{"name": strings.ToUpper(name)}).One(&results)
	return results
}

func (r Repository) FindByNameAndZip(name, zip string) Companies {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	results := Companies{}
	c := session.DB(DBNAME).C(DOCNAME)
	search := bson.RegEx{name+".*", ""}
	c.Find(bson.M{"name": search, "addresszip" : zip }).All(&results)
	return results
}
