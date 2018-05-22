package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
)

//Util ...
type Util struct {
	Repository Repository
}

func (u Util) loadCsvData() {
	companies := Companies{}
	q1Catalog, err := os.Open("csv/q1_catalog.csv")
	if err != nil {
		log.Printf(
			"%s\t%s\t%s\t%s",
			"Error on load CSV data: ",
			"Loading file",
			"Check",
			time.Since(time.Now()),
		)
		return
	}
	defer q1Catalog.Close()
	reader := csv.NewReader(q1Catalog)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	for i, record := range records {
		if i == 0 {
			continue
		}
		company := Company{Name: strings.ToUpper(record[0]), AddressZip: record[1]}
		companies = append(companies, company)
	}
	success := u.Repository.AddCompanies(companies)
	if !success {
		log.Printf(
			"%s\t%s\t%s\t%s",
			"Error on load CSV data: ",
			"Persisting loaded records in MongoDB",
			"Check",
			time.Since(time.Now()),
		)
		return
	}
}

func (u Util) mergeCompanies(companies Companies) bool {

	for _, company := range companies {
		retrievedCompany := u.Repository.FindByName(company.Name)
		if retrievedCompany.Name != "" {
			retrievedCompany.Website = company.Website
			retrievedCompany.AddressZip = company.AddressZip
			success := u.Repository.UpdateCompany(retrievedCompany)
			if !success {
				log.Printf(
					"%s\t%s\t%s\t%s",
					"Error on merge data: ",
					"User: ",
					retrievedCompany.Name,
					time.Since(time.Now()),
				)
				return false
			}
		}
	}
	return true
}
