package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateApplication(ctx context.Context, oktaConfig *okta.Configuration) (*methodokta.ApplicationReport, error) {
	resources := methodokta.ApplicationReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(context.Background()).Execute()
	if err != nil {
		return &methodokta.ApplicationReport{}, err
	}

	// Fetch all Applications
	apps, resp, err := client.ApplicationAPI.ListApplications(ctx).Expand("").Execute()
	if err != nil {
		return &methodokta.ApplicationReport{}, err
	}
	var allApps []okta.ListApplications200ResponseInner
	allApps = append(allApps, apps...)
	for resp.HasNextPage() {
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor := parsedURL.Query().Get("after")
		apps, resp, err = client.ApplicationAPI.ListApplications(ctx).After(cursor).Execute()
		if err != nil {
			return &methodokta.ApplicationReport{}, err
		}
		allApps = append(allApps, apps...)
	}

	// Loop through Applications
	var appList []*methodokta.Application
	for _, a := range allApps {

		// Application data
		data, _ := a.MarshalJSON()
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			uid := result["id"].(string)
			name := result["name"].(string)
			label := result["label"].(string)
			status := result["status"].(string)
			statusEnum, err := methodokta.NewStatusTypeFromString(status)
			if err != nil {
				statusEnum, _ = methodokta.NewStatusTypeFromString("UNKNOWN")
			}
			created, err := time.Parse(time.RFC3339, result["created"].(string))
			if err != nil {
				errors = append(errors, err.Error())
			}
			authMethod := result["signOnMode"].(string)
			authEnum, err := methodokta.NewAuthTypeFromString(authMethod)
			if err != nil {
				authEnum, _ = methodokta.NewAuthTypeFromString("UNKNOWN")
			}
			setting := result["settings"]
			var appURL string
			if setting != nil {
				signOn := setting.(map[string]interface{})["signOn"]
				if signOn != nil {
					if signOn.(map[string]interface{})["ssoAcsUrl"] != nil {
						appURL = signOn.(map[string]interface{})["ssoAcsUrl"].(string)
					}
				}
			}
			var baseURL *string
			if appURL != "" {
				parsedURL, err := url.Parse(appURL)
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					val := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
					baseURL = &val
				}
			}
			application := methodokta.Application{
				Uid:        uid,
				Name:       name,
				Label:      label,
				Url:        baseURL,
				Status:     statusEnum,
				Created:    created,
				AuthMethod: authEnum,
			}

			// Application groups
			var allGroups []okta.ApplicationGroupAssignment
			groups, resp, err := client.ApplicationGroupsAPI.ListApplicationGroupAssignments(ctx, uid).Execute()
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				allGroups = append(allGroups, groups...)
				for resp.HasNextPage() {
					parsedURL, _ := url.Parse(resp.NextPage())
					cursor := parsedURL.Query().Get("after")
					groups, resp, err = client.ApplicationGroupsAPI.ListApplicationGroupAssignments(ctx, uid).After(cursor).Execute()
					if err != nil {
						errors = append(errors, err.Error())
					} else {
						allGroups = append(allGroups, groups...)
					}
				}
			}
			for _, g := range allGroups {
				group, _, err := client.GroupAPI.GetGroup(ctx, *g.Id).Execute()
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					name := group.Profile.Name
					group := methodokta.GroupInfo{Uid: *g.Id, Name: *name}
					application.Groups = append(application.Groups, &group)
				}
			}

			// Application Users
			var allUsers []okta.AppUser
			users, resp, err := client.ApplicationUsersAPI.ListApplicationUsers(ctx, uid).Execute()
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				allUsers = append(allUsers, users...)
				for resp.HasNextPage() {
					parsedURL, _ := url.Parse(resp.NextPage())
					cursor := parsedURL.Query().Get("after")
					users, resp, err = client.ApplicationUsersAPI.ListApplicationUsers(ctx, uid).After(cursor).Execute()
					if err != nil {
						errors = append(errors, err.Error())
					} else {
						allUsers = append(allUsers, users...)
					}
				}
			}
			for _, u := range allUsers {
				user, _, err := client.UserAPI.GetUser(ctx, *u.Id).Execute()
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					email := user.Profile.Email
					user := methodokta.UserInfo{Uid: *u.Id, Email: *email}
					application.Users = append(application.Users, &user)
				}
			}

			appList = append(appList, &application)
		}
	}

	resources = methodokta.ApplicationReport{
		Applications: appList,
		Org:          *org.Id,
		Errors:       errors,
	}

	return &resources, nil
}
