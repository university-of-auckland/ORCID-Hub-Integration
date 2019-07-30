package main

import (
	"fmt"
	"log"
	"strings"
)

// Identity - the user identity record.
type Identity struct {
	Addresses []struct {
		CareOf      string `json:"careOf"`
		City        string `json:"city"`
		Country     string `json:"country"`
		CountryID   string `json:"countryId"`
		Dpid        string `json:"dpid"`
		LastUpdated string `json:"lastUpdated"`
		Line1       string `json:"line1"`
		Line2       string `json:"line2"`
		Line3       string `json:"line3"`
		Postcode    string `json:"postcode"`
		State       string `json:"state"`
		StateID     string `json:"stateId"`
		Suburb      string `json:"suburb"`
		Type        string `json:"type"`
		TypeID      string `json:"typeId"`
	} `json:"addresses"`
	Citizenships []struct {
		Country   string `json:"country"`
		CountryID string `json:"countryId"`
	} `json:"citizenships"`
	CountryOfBirth struct {
		Country   string `json:"country"`
		CountryID string `json:"countryId"`
	} `json:"countryOfBirth"`
	Deceased struct {
		Comments         string `json:"comments"`
		Date             string `json:"date"`
		Dead             bool   `json:"dead"`
		DeathCertificate string `json:"deathCertificate"`
		Place            string `json:"place"`
	} `json:"deceased"`
	DisabilityInfo struct {
		Disabilities []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"disabilities"`
		DisabilityPermanent bool `json:"disabilityPermanent"`
		DisabilitySupport   bool `json:"disabilitySupport"`
		IsDisabled          bool `json:"isDisabled"`
	} `json:"disabilityInfo"`
	DisplayName  string `json:"displayName"`
	Dob          string `json:"dob"`
	EmailAddress string `json:"emailAddress"`
	Emails       []struct {
		Email       string `json:"email"`
		LastUpdated string `json:"lastUpdated"`
		Type        string `json:"type"`
		TypeID      string `json:"typeId"`
		Verified    bool   `json:"verified"`
	} `json:"emails"`
	EmergencyContacts []struct {
		EmailAddress string `json:"emailAddress"`
		LastUpdated  string `json:"lastUpdated"`
		Name         string `json:"name"`
		Phones       []struct {
			AreaCode    string `json:"areaCode"`
			CountryCode string `json:"countryCode"`
			Extension   string `json:"extension"`
			LastUpdated string `json:"lastUpdated"`
			Number      string `json:"number"`
			Type        string `json:"type"`
			TypeID      string `json:"typeId"`
		} `json:"phones"`
		Relationship string `json:"relationship"`
	} `json:"emergencyContacts"`
	Ethnicities []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"ethnicities"`
	ExtIds []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"extIds"`
	FirstName       string   `json:"firstName"`
	Gender          string   `json:"gender"`
	Groups          []string `json:"groups"`
	ID              int      `json:"id"`
	IDPhotoExists   bool     `json:"idPhotoExists"`
	IwiAffiliations []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"iwiAffiliations"`
	LastName   string `json:"lastName"`
	MergedToID string `json:"mergedToId"`
	Mobile     string `json:"mobile"`
	Names      []struct {
		First       string `json:"first"`
		Last        string `json:"last"`
		LastUpdated string `json:"lastUpdated"`
		Middle      string `json:"middle"`
		Suffix      string `json:"suffix"`
		Title       string `json:"title"`
		Type        string `json:"type"`
	} `json:"names"`
	Phones []struct {
		AreaCode    string `json:"areaCode"`
		CountryCode string `json:"countryCode"`
		Extension   string `json:"extension"`
		LastUpdated string `json:"lastUpdated"`
		Number      string `json:"number"`
		Type        string `json:"type"`
		TypeID      string `json:"typeId"`
	} `json:"phones"`
	PreviousIds     []string `json:"previousIds"`
	PrimaryIdentity bool     `json:"primaryIdentity"`
	Residency       string   `json:"residency"`
	Resolved        bool     `json:"resolved"`
	Upi             string   `json:"upi"`
	WhenUpdated     string   `json:"whenUpdated"`
}

// GetORCID returns the principal part of ORCID iD
// if the identify has an ORCID.
func (id *Identity) GetORCID() string {
	if id.ExtIds == nil {
		return ""
	}
	for _, eid := range id.ExtIds {
		if eid.Type == "ORCID" {
			parts := strings.Split(eid.ID, "/")
			return parts[len(parts)-1]
		}
	}
	return ""
}

// GetOrcidAccessToken gets the ORCID API token to verify that the user
// has granted access to the university.
func (id *Identity) GetOrcidAccessToken() (token Token, ok bool) {
	var tokens []Token
	orcid := id.GetORCID()
	if orcid != "" {
		err := oh.Get("api/v1/tokens/"+orcid, &tokens)
		if err != nil {
			log.Println("ERROR: ", err)
		} else if len(tokens) > 0 {
			goto TOKEN_FOUND
		}
	}
	{
		otherIDs := make([]string, len(id.Emails)+2)
		otherIDs[0] = id.Upi + "@auckland.ac.nz"
		otherIDs[1] = id.EmailAddress
		for i, a := range id.Emails {
			otherIDs[i+2] = a.Email
		}
		for _, oid := range otherIDs {
			if oid != "" {
				err := oh.Get("api/v1/tokens/"+oid, &tokens)
				if err != nil {
					log.Println("ERROR: ", err)
				} else if len(tokens) > 0 {
					goto TOKEN_FOUND
				}
			}
		}
	}
	return
TOKEN_FOUND:
	for _, token := range tokens {
		if strings.Contains(token.Scopes, "update") {
			return token, true
		}
	}
	return
}

// updateOrcid updates the user ORCID iD.
func (id *Identity) updateOrcid(ORCID string) {
	currentORCID := id.GetORCID()
	if currentORCID != "" {
		if ORCID != currentORCID {
			// TODO
		}
		return
	}

	wg.Add(1)
	defer func() {
		wg.Done()
	}()

	// Add ORCID ID if the user doesn't have one
	var resp struct {
		StatusCode string `json:"statusCode"`
	}

	err := api.Put(fmt.Sprintf("identity/integrations/v3/identity/%d/identifier/ORCID", id.ID), map[string]string{"identifier": ORCID}, &resp)
	if err != nil {
		log.Println("ERROR: Failed to update or add ORCID: ", err)
	}
}
