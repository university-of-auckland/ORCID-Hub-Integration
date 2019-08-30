package main

import (
	"fmt"
	"strings"
)

// Identity - the user identity record.
type Identity struct {
	EmailAddress string `json:"emailAddress"`
	Emails       []struct {
		Email       string `json:"email"`
		LastUpdated string `json:"lastUpdated"`
		Type        string `json:"type"`
		TypeID      string `json:"typeId"`
		Verified    bool   `json:"verified"`
	} `json:"emails"`
	ExtIds []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"extIds,omitempty"`
	ID  int    `json:"id"`
	Upi string `json:"upi"`
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
	if id.EmailAddress == "" || id.Upi == "" {
		return
	}
	var tokens []Token
	orcid := id.GetORCID()

	if orcid != "" {
		err := oh.get("api/v1/tokens/"+orcid, &tokens)
		if err != nil {
			log.Error(err)
		} else if len(tokens) > 0 {
			goto TOKEN_FOUND
		}
	}
	if id.Upi != "" || id.EmailAddress != "" || id.Emails != nil {
		otherIDs := make([]string, len(id.Emails)+2)
		otherIDs[0] = id.Upi + "@auckland.ac.nz"
		otherIDs[1] = id.EmailAddress
		for i, a := range id.Emails {
			otherIDs[i+2] = a.Email
		}
		for _, oid := range otherIDs {
			if oid != "" {
				err := oh.get("api/v1/tokens/"+oid, &tokens)
				if err != nil {
					log.Error(err)
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

	err := api.put(fmt.Sprintf("identity/integrations/v3/identity/%d/identifier/ORCID", id.ID), map[string]string{"identifier": ORCID}, &resp)
	if err != nil {
		log.Error("failed to update or add ORCID: ", err)
	}
}
