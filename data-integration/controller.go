package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//Controller ...
type Controller struct {
	Repository Repository
	Util       Util
}

// Index GET /v1/companies
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	companies := c.Repository.GetCompanies() // list of all companies
	log.Println(companies)
	data, _ := json.Marshal(companies)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddCompany POST /
func (c *Controller) AddCompany(w http.ResponseWriter, r *http.Request) {
	var company Company
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error AddCompany", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddCompany", err)
	}
	if err := json.Unmarshal(body, &company); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddCompany unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	success := c.Repository.AddCompany(company)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// UpdateCompany PUT /
func (c *Controller) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var company Company
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error UpdateCompany", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdateCompany", err)
	}
	if err := json.Unmarshal(body, &company); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateCompany unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	success := c.Repository.UpdateCompany(company)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// DeleteCompany DELETE /
func (c *Controller) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]                                      // param id
	if err := c.Repository.DeleteCompany(id); err != "" { // delete a album by id
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// MergeCompany POST /
func (c *Controller) MergeCompany(w http.ResponseWriter, r *http.Request) {
	companies := Companies{}
	file, fileHandler, err := r.FormFile("csv")
	if err != nil {
		log.Fatalln("Error MergeCompany", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(fileHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Error MergeCompany", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	uploadedFile, _ := os.Open(fileHandler.Filename)
	reader := csv.NewReader(bufio.NewReader(uploadedFile))
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Error MergeCompany", err)
		w.WriteHeader(http.StatusBadRequest)
		//w.Write
		return
	}
	if len(records) == 0 {
		log.Fatalln("Empty file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		company := Company{Name: record[0], AddressZip: record[1], Website: record[2]}
		companies = append(companies, company)
	}
	success := c.Util.mergeCompanies(companies)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return

}

// SearchCompanies SearchCompanies /
func (c *Controller) SearchCompanies(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	zip := r.URL.Query().Get("zip")
	companies := c.Repository.FindByNameAndZip(name, zip)
	data, _ := json.Marshal(companies)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
