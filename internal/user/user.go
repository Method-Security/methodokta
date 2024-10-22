package user

import (
	"context"
	"encoding/json"
	"log"
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
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.UserReport{}, err
	}

	// Print Total Users
	printUsersCmd := client.GroupAPI.ListGroups(ctx)
	err = printUserTotal(printUsersCmd)
	if err != nil {
		return &methodokta.UserReport{}, err
	}

	// Fetch all Users
	getUserscmd := client.UserAPI.ListUsers(ctx)
	allUsers, err := fetchUsersWithRetry(getUserscmd, sleep)
	if err != nil {
		return &methodokta.UserReport{}, err
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
		getAppsCmd := client.UserAPI.ListAppLinks(ctx, *u.Id)
		allAppLinks, err := fetchListAppLinksWithRetry(getAppsCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, a := range allAppLinks {
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
		getGroupsCmd := client.UserAPI.ListUserGroups(ctx, *u.Id)
		allGroups, err := fetchListUserGroupsWithRetry(getGroupsCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, g := range allGroups {
				group := methodokta.GroupInfo{Uid: *g.Id, Name: *g.Profile.Name}
				user.Groups = append(user.Groups, &group)
			}
		}

		// User roles
		getRolesCmd := client.RoleAssignmentAPI.ListAssignedRolesForUser(ctx, *u.Id)
		roles, err := fetchListAssignedRolesForUserWithRetry(getRolesCmd, sleep)
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

func printUserTotal(cmd okta.ApiListGroupsRequest) error {
	groups, _, err := cmd.Q("everyone").Expand("stats").Execute()
	if err != nil {
		return err
	}

	for _, group := range groups {
		if embeddedStats, ok := group.Embedded["stats"]; ok {
			if usersCount, ok := embeddedStats["usersCount"].(float64); ok {
				log.Printf("Total Users: %d", int(usersCount))
			} else {
				log.Print("usersCount not found in group stats")
			}
		} else {
			log.Print("_embedded.stats not found")
		}
	}
	return nil
}

func fetchUsersWithRetry(cmd okta.ApiListUsersRequest, sleep time.Duration) ([]okta.User, error) {
	var allUsers []okta.User
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Printf("USERS: fetchUsers sleep - %v", sleepExp)
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
		log.Printf("USERS: fetchUsers count - %v", len(allUsers))
	}
	return allUsers, nil
}

func fetchListAppLinksWithRetry(cmd okta.ApiListAppLinksRequest, sleep time.Duration) ([]okta.AppLink, error) {
	sleepExp := sleep
	apps, _, err := cmd.Execute()
	for err != nil {
		apps, _, err = cmd.Execute()
		if err != nil {
			log.Printf("APPS: fetchListAppLinks sleep - %v", sleepExp)
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	log.Printf("APPS: fetchListAppLinks count - %v", len(apps))
	return apps, nil
}

func fetchListUserGroupsWithRetry(cmd okta.ApiListUserGroupsRequest, sleep time.Duration) ([]okta.Group, error) {
	var allGroups []okta.Group
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		groups, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Printf("GROUPS: fetchListUserGroups sleep - %v", sleepExp)
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
		log.Printf("GROUPS: fetchListUserGroups count - %v", len(allGroups))
	}
	return allGroups, nil
}

func fetchListAssignedRolesForUserWithRetry(cmd okta.ApiListAssignedRolesForUserRequest, sleep time.Duration) ([]okta.Role, error) {
	sleepExp := sleep
	roles, _, err := cmd.Execute()
	for err != nil {
		roles, _, err = cmd.Execute()
		if err != nil {
			log.Printf("ROLES: fetchListAssignedRolesForUser sleep - %v", sleepExp)
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	log.Printf("ROLES: fetchListAssignedRolesForUser count - %v", len(roles))
	return roles, nil
}
