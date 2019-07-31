package main

// Token - ORCID API access token
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IssueTime    string `json:"issue_time"`
	RefreshToken string `json:"refresh_token"`
	Scopes       string `json:"scopes"`
	Email        string `json:"email"`
	EPPN         string `json:"eppn"`
	ORCID        string `json:"orcid"`
}
