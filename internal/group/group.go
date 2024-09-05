package group

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateGroup(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.GroupReport, error) {
	resources := methodokta.GroupReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.GroupReport{}, err
	}

	// Fetch all Groups
	getUsersCmd := client.GroupAPI.ListGroups(ctx)
	allGroups, err := fetchListGroupsWithRetry(getUsersCmd, sleep)
	if err != nil {
		return &methodokta.GroupReport{}, err
	}

	// Loop through Groups
	var groupList []*methodokta.Group
	for _, g := range allGroups {

		// Group data
		groupType := *g.Type
		groupTypeEnum, err := methodokta.NewGroupTypeFromString(groupType)
		if err != nil {
			groupTypeEnum, _ = methodokta.NewGroupTypeFromString("UNKNOWN")
		}
		group := methodokta.Group{
			Uid:         *g.Id,
			Name:        *g.Profile.Name,
			Type:        groupTypeEnum,
			Description: g.Profile.Description,
			Created:     *g.Created,
		}

		// Group Applications
		getAppsCmd := client.GroupAPI.ListAssignedApplicationsForGroup(ctx, *g.Id)
		allApps, err := fetchListAssignedApplicationsForGroupWithRetry(getAppsCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, a := range allApps {
				data, _ := a.MarshalJSON()
				var result map[string]interface{}
				err = json.Unmarshal(data, &result)
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					application := methodokta.ApplicationInfo{Uid: result["id"].(string), Name: result["label"].(string)}
					group.Applications = append(group.Applications, &application)
				}
			}
		}

		// Group Roles
		getRolesCmd := client.RoleAssignmentAPI.ListGroupAssignedRoles(ctx, *g.Id)
		allRoles, err := fetchListGroupAssignedRolesWithRetry(getRolesCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, r := range allRoles {
				role := methodokta.RoleInfo{Uid: *r.Id, Name: r.Label}
				group.Roles = append(group.Roles, &role)
			}
		}

		// Group Users
		getUsersCmd := client.GroupAPI.ListGroupUsers(ctx, *g.Id)
		allUsers, err := fetchListGroupUsersWithRetry(getUsersCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {

			for _, u := range allUsers {
				user := methodokta.UserInfo{Uid: *u.Id, Email: *u.Profile.Email}
				group.Users = append(group.Users, &user)
			}
		}

		groupList = append(groupList, &group)
	}

	resources = methodokta.GroupReport{
		Groups: groupList,
		Org:    *org.Id,
		Errors: errors,
	}

	return &resources, nil
}

func fetchListGroupsWithRetry(cmd okta.ApiListGroupsRequest, sleep time.Duration) ([]okta.Group, error) {
	var allGroups []okta.Group
	var changePage bool
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		groups, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			changePage = false
			sleepExp *= 2
		} else {
			changePage = true
		}
		if changePage {
			sleepExp = sleep
			parsedURL, _ := url.Parse(resp.NextPage())
			cursor = parsedURL.Query().Get("after")
			hasNextPage = resp.HasNextPage()
			allGroups = append(allGroups, groups...)
		}
	}
	return allGroups, nil
}

func fetchListAssignedApplicationsForGroupWithRetry(cmd okta.ApiListAssignedApplicationsForGroupRequest, sleep time.Duration) ([]okta.ListApplications200ResponseInner, error) {
	var allApps []okta.ListApplications200ResponseInner
	var changePage bool
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			changePage = false
			sleepExp *= 2
		} else {
			changePage = true
		}
		if changePage {
			sleepExp = sleep
			parsedURL, _ := url.Parse(resp.NextPage())
			cursor = parsedURL.Query().Get("after")
			hasNextPage = resp.HasNextPage()
			allApps = append(allApps, users...)
		}
	}
	return allApps, nil
}

func fetchListGroupAssignedRolesWithRetry(cmd okta.ApiListGroupAssignedRolesRequest, sleep time.Duration) ([]okta.Role, error) {
	sleepExp := sleep
	roles, _, err := cmd.Execute()
	for err != nil {
		roles, _, err = cmd.Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	return roles, nil
}

func fetchListGroupUsersWithRetry(cmd okta.ApiListGroupUsersRequest, sleep time.Duration) ([]okta.GroupMember, error) {
	var allUsers []okta.GroupMember
	var changePage bool
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			changePage = false
			sleepExp *= 2
		} else {
			changePage = true
		}
		if changePage {
			sleepExp = sleep
			parsedURL, _ := url.Parse(resp.NextPage())
			cursor = parsedURL.Query().Get("after")
			hasNextPage = resp.HasNextPage()
			allUsers = append(allUsers, users...)
		}
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
