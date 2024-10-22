package group

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func EnumerateGroup(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.GroupReport, error) {
	log := svc1log.FromContext(ctx)

	resources := methodokta.GroupReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.GroupReport{}, err
	}

	// Fetch all Groups
	log.Info("Total Groups")
	getGroupsCmd := client.GroupAPI.ListGroups(ctx)
	allGroups, err := fetchListGroupsWithRetry(ctx, getGroupsCmd, sleep)
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

		log.Info("List Applications + Roles + Users for Group", svc1log.SafeParam("Name", *g.Profile.Name))

		// Group Applications
		getAppsCmd := client.GroupAPI.ListAssignedApplicationsForGroup(ctx, *g.Id)
		allApps, err := fetchListAssignedApplicationsForGroupWithRetry(ctx, getAppsCmd, sleep)
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
		allRoles, err := fetchListGroupAssignedRolesWithRetry(ctx, getRolesCmd, sleep)
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
		allUsers, err := fetchListGroupUsersWithRetry(ctx, getUsersCmd, sleep)
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

func fetchListGroupsWithRetry(ctx context.Context, cmd okta.ApiListGroupsRequest, sleep time.Duration) ([]okta.Group, error) {
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

func fetchListAssignedApplicationsForGroupWithRetry(ctx context.Context, cmd okta.ApiListAssignedApplicationsForGroupRequest, sleep time.Duration) ([]okta.ListApplications200ResponseInner, error) {
	log := svc1log.FromContext(ctx)

	var allApps []okta.ListApplications200ResponseInner
	var changePage bool
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		users, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Info("Applications", svc1log.SafeParam("sleep", sleepExp))
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
			log.Info("Applications", svc1log.SafeParam("count", len(allApps)))
		}
	}
	return allApps, nil
}

func fetchListGroupAssignedRolesWithRetry(ctx context.Context, cmd okta.ApiListGroupAssignedRolesRequest, sleep time.Duration) ([]okta.Role, error) {
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

func fetchListGroupUsersWithRetry(ctx context.Context, cmd okta.ApiListGroupUsersRequest, sleep time.Duration) ([]okta.GroupMember, error) {
	log := svc1log.FromContext(ctx)

	var allUsers []okta.GroupMember
	var changePage bool
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
			log.Info("Users", svc1log.SafeParam("count", len(allUsers)))
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
