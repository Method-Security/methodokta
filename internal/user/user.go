package user

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateUser(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.UserReport, error) {
	resources := methodokta.UserReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(context.Background()).Execute()
	if err != nil {
		return &methodokta.UserReport{}, err
	}

	// Fetch all users
	var allUsers []okta.User
	users, resp, err := client.UserAPI.ListUsers(ctx).Execute()
	if err != nil {
		errors = append(errors, err.Error())
		time.Sleep(sleep)
		users, resp, err = client.UserAPI.ListUsers(ctx).Execute()
		if err != nil {
			return &methodokta.UserReport{}, err
		}
	}
	allUsers = append(allUsers, users...)
	for resp.HasNextPage() {
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor := parsedURL.Query().Get("after")
		users, resp, err = client.UserAPI.ListUsers(ctx).After(cursor).Execute()
		if err != nil {
			errors = append(errors, err.Error())
			time.Sleep(sleep)
			users, resp, err = client.UserAPI.ListUsers(ctx).After(cursor).Execute()
			if err != nil {
				return &methodokta.UserReport{}, err
			}
		}
		allUsers = append(allUsers, users...)
	}

	// Loop through Users
	var userList []*methodokta.User
	for _, u := range allUsers {

		// User data
		var firstname *string
		if u.Profile.FirstName.IsSet() {
			firstname = u.Profile.FirstName.Get()
		}
		var lastname *string
		if u.Profile.FirstName.IsSet() {
			lastname = u.Profile.LastName.Get()
		}
		var PasswordChanged *time.Time
		if u.PasswordChanged.IsSet() {
			PasswordChanged = u.PasswordChanged.Get()
		}
		status := *u.Status
		statusEnum, err := methodokta.NewStatusTypeFromString(status)
		if err != nil {
			statusEnum, _ = methodokta.NewStatusTypeFromString("UNKNOWN")
		}
		user := methodokta.User{
			Uid:             *u.Id,
			Firstname:       firstname,
			Lastname:        lastname,
			Email:           *u.Profile.Email,
			Status:          statusEnum,
			Created:         *u.Created,
			PasswordChanged: PasswordChanged,
		}

		// User applications
		appLinks, _, err := client.UserAPI.ListAppLinks(ctx, *u.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, a := range appLinks {
				data, _ := a.MarshalJSON()
				var result map[string]interface{}
				err = json.Unmarshal(data, &result)
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					appName := result["label"].(string)
					if strings.Contains(appName, "Google Workspace") {
						appName = "Google Workspace"
					}
					application := methodokta.ApplicationInfo{Uid: result["appInstanceId"].(string), Name: appName}

					user.Applications = append(user.Applications, &application)
				}
			}
		}

		// User groups
		var allGroups []okta.Group
		groups, resp, err := client.UserAPI.ListUserGroups(ctx, *u.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			allGroups = append(allGroups, groups...)
			for resp.HasNextPage() {
				parsedURL, _ := url.Parse(resp.NextPage())
				cursor := parsedURL.Query().Get("after")
				groups, resp, err = client.UserAPI.ListUserGroups(ctx, *u.Id).After(cursor).Execute()
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					allGroups = append(allGroups, groups...)
				}
			}
		}
		for _, g := range allGroups {
			group := methodokta.GroupInfo{Uid: *g.Id, Name: *g.Profile.Name}
			user.Groups = append(user.Groups, &group)
		}

		// User roles
		roles, _, err := client.RoleAssignmentAPI.ListAssignedRolesForUser(ctx, *u.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, r := range roles {
				roleType := *r.Type
				roleEnum, err := methodokta.NewRoleTypeFromString(roleType)
				if err != nil {
					roleEnum, _ = methodokta.NewRoleTypeFromString("UNKNOWN")
				} else {
					role := methodokta.RoleInfo{
						Uid:         *r.Id,
						Name:        r.Label,
						Type:        &roleEnum,
						Description: r.Description,
						Created:     r.Created,
					}
					user.Roles = append(user.Roles, &role)
				}
			}
		}

		// Complete User data List
		userList = append(userList, &user)
	}

	resources = methodokta.UserReport{
		Users:  userList,
		Org:    *org.Id,
		Errors: errors,
	}
	return &resources, nil
}
