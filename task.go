package main

type Task struct {
	ID          int      `json:"id"`
	CompletedAt string   `json:"completed-at"`
	CreatedAt   string   `json:"created-at"`
	ExpiresAt   string   `json:"expires-at"`
	Filename    string   `json:"filename"`
	Status      string   `json:"status"`
	Type        string   `json:"task-type"`
	Records     []Record `json:"records"`
}

type Record struct {
	ID                  int    `json:"id"`
	AffiliationType     string `json:"affiliation-type"`
	City                string `json:"city"`
	Country             string `json:"country"`
	Department          string `json:"department"`
	DisambiguatedID     string `json:"disambiguated-id"`
	DisambiguatedSource string `json:"disambiguated-source"`
	Email               string `json:"email"`
	EndDate             string `json:"end-date"`
	ExternalID          string `json:"external-id"`
	FirstName           string `json:"first-name"`
	IsActive            bool   `json:"is-active"`
	LastName            string `json:"last-name"`
	Orcid               string `json:"orcid"`
	Organisation        string `json:"organisation"`
	ProcessedAt         string `json:"processed-at"`
	PutCode             string `json:"put-code"`
	Role                string `json:"role"`
	StartDate           string `json:"start-date"`
	State               string `json:"state"`
	Status              string `json:"status"`
}
