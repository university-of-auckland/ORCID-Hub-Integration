package main

import "time"

type Identity struct {
	Department   string `json:"department"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Emails       []struct {
		Email       string    `json:"email"`
		LastUpdated time.Time `json:"lastUpdated"`
		Order       int       `json:"order"`
		Type        string    `json:"type"`
		Verified    bool      `json:"verified"`
	} `json:"emails"`
	ExtIds []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"extIds"`
	FirstName string `json:"firstName"`
	JobTitle  string `json:"jobTitle"`
	LastName  string `json:"lastName"`
	Mobile    string `json:"mobile"`
	Names     []struct {
		First  string `json:"first"`
		Last   string `json:"last"`
		Middle string `json:"middle"`
		Suffix string `json:"suffix"`
		Title  string `json:"title"`
		Type   string `json:"type"`
	} `json:"names"`
	Upi string `json:"upi"`
}
