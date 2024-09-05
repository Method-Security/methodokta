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

func EnumerateApplication(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.ApplicationReport, error) {
	resources := methodokta.ApplicationReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.ApplicationReport{}, err
	}

	// Fetch all Applications
	getAppsCmd := client.ApplicationAPI.ListApplications(ctx).Expand("")
	allApps, err := fetchListApplicationsWithRetry(getAppsCmd, sleep)
	if err != nil {
		return &methodokta.ApplicationReport{}, err
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
			getGroupsCmd := client.ApplicationGroupsAPI.ListApplicationGroupAssignments(ctx, uid)
			allGroups, err := fetchListApplicationGroupAssignmentsWithRetry(getGroupsCmd, sleep)
			if err != nil {
				errors = append(errors, err.Error())
			} else {
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
			}

			// Application Users
			getUsersCmd := client.ApplicationUsersAPI.ListApplicationUsers(ctx, uid)
			allUsers, err := fetchListApplicationUsersWithRetry(getUsersCmd, sleep)
			if err != nil {
				errors = append(errors, err.Error())
			} else {
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

func fetchListApplicationsWithRetry(cmd okta.ApiListApplicationsRequest, sleep time.Duration) ([]okta.ListApplications200ResponseInner, error) {
	var allApps []okta.ListApplications200ResponseInner
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		apps, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
			continue
		}
		sleepExp = sleep
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor = parsedURL.Query().Get("after")
		hasNextPage = resp.HasNextPage()
		allApps = append(allApps, apps...)
	}
	return allApps, nil
}

func fetchListApplicationGroupAssignmentsWithRetry(cmd okta.ApiListApplicationGroupAssignmentsRequest, sleep time.Duration) ([]okta.ApplicationGroupAssignment, error) {
	var allGroups []okta.ApplicationGroupAssignment
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		groups, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
			continue
		}
		sleepExp = sleep
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor = parsedURL.Query().Get("after")
		hasNextPage = resp.HasNextPage()
		allGroups = append(allGroups, groups...)
	}
	return allGroups, nil
}

func fetchListApplicationUsersWithRetry(cmd okta.ApiListApplicationUsersRequest, sleep time.Duration) ([]okta.AppUser, error) {
	var allUsers []okta.AppUser
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
			continue
		}
		sleepExp = sleep
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor = parsedURL.Query().Get("after")
		hasNextPage = resp.HasNextPage()
		allUsers = append(allUsers, users...)
	}
	return allUsers, nil
}

func retry(sleep time.Duration, err error) bool {
	if err.Error() != "too many requests" {
		return false
	}
	time.Sleep(sleep)
	return true
}
