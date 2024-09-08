package org

import (
	"context"
	"fmt"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateOrg(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.OrgReport, error) {
	resources := methodokta.OrgReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(context.Background()).Execute()
	if err != nil {
		errors = append(errors, err.Error())
		time.Sleep(sleep)
		org, _, err = client.OrgSettingAPI.GetOrgSettings(context.Background()).Execute()
		if err != nil {
			return &methodokta.OrgReport{}, err
		}
	}

	// Org data
	var url string
	if org.Subdomain != nil {
		url = fmt.Sprintf("https://%s.okta.com", *org.Subdomain)
	} else {
		return &methodokta.OrgReport{}, err
	}
	status := *org.Status
	statusEnum, err := methodokta.NewStatusTypeFromString(status)
	if err != nil {
		statusEnum, _ = methodokta.NewStatusTypeFromString("UNKNOWN")
	}
	orgInfo := methodokta.OrgInfo{
		Uid:         *org.Id,
		CompanyName: org.CompanyName,
		Url:         url,
		Status:      statusEnum,
		Created:     *org.Created,
	}

	resources = methodokta.OrgReport{
		Organization: &orgInfo,
		Org:          *org.Id,
		Errors:       errors,
	}

	return &resources, nil

}
