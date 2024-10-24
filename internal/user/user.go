package user

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func EnumerateUser(ctx context.Context, limit int, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.UserReport, error) {
	log := svc1log.FromContext(ctx)
	resources := methodokta.UserReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.UserReport{}, err
	}

	// Fetch all Users
	log.Info("Total Users")
	getUserscmd := client.UserAPI.ListUsers(ctx)
	allUsers, err := fetchUsersWithRetry(ctx, getUserscmd, limit, sleep)
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

		log.Info("List Applications + Groups + Roles for User", svc1log.SafeParam("ID", *u.Id))

		// User applications
		getAppsCmd := client.UserAPI.ListAppLinks(ctx, *u.Id)
		allAppLinks, err := fetchListAppLinksWithRetry(ctx, getAppsCmd, sleep)
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
		allGroups, err := fetchListUserGroupsWithRetry(ctx, getGroupsCmd, sleep)
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
		roles, err := fetchListAssignedRolesForUserWithRetry(ctx, getRolesCmd, sleep)
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

func fetchUsersWithRetry(ctx context.Context, cmd okta.ApiListUsersRequest, limit int, sleep time.Duration) ([]okta.User, error) {
	log := svc1log.FromContext(ctx)

	var allUsers []okta.User
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Info("Users", svc1log.SafeParam("sleep", sleepExp))
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

		if limit > 0 && len(allUsers)+len(users) >= limit {
			remaining := limit - len(allUsers)
			allUsers = append(allUsers, users[:remaining]...)
			log.Info("Users", svc1log.SafeParam("count", len(allUsers)))
			return allUsers, nil
		}

		allUsers = append(allUsers, users...)
		log.Info("Users", svc1log.SafeParam("count", len(allUsers)))
	}
	return allUsers, nil
}

func fetchListAppLinksWithRetry(ctx context.Context, cmd okta.ApiListAppLinksRequest, sleep time.Duration) ([]okta.AppLink, error) {
	log := svc1log.FromContext(ctx)

	sleepExp := sleep
	apps, _, err := cmd.Execute()
	for err != nil {
		apps, _, err = cmd.Execute()
		if err != nil {
			log.Info("Applications", svc1log.SafeParam("sleep", sleepExp))
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	log.Info("Applications", svc1log.SafeParam("count", len(apps)))
	return apps, nil
}

func fetchListUserGroupsWithRetry(ctx context.Context, cmd okta.ApiListUserGroupsRequest, sleep time.Duration) ([]okta.Group, error) {
	log := svc1log.FromContext(ctx)

	var allGroups []okta.Group
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		groups, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Info("Groups", svc1log.SafeParam("sleep", sleepExp))
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
		log.Info("Groups", svc1log.SafeParam("count", len(allGroups)))
	}
	return allGroups, nil
}

func fetchListAssignedRolesForUserWithRetry(ctx context.Context, cmd okta.ApiListAssignedRolesForUserRequest, sleep time.Duration) ([]okta.Role, error) {
	log := svc1log.FromContext(ctx)

	sleepExp := sleep
	roles, _, err := cmd.Execute()
	for err != nil {
		roles, _, err = cmd.Execute()
		if err != nil {
			log.Info("Roles", svc1log.SafeParam("sleep", sleepExp))
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	log.Info("Roles", svc1log.SafeParam("count", len(roles)))
	return roles, nil
}
