package group

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateGroup(ctx context.Context, oktaConfig *okta.Configuration) (*methodokta.GroupReport, error) {
	resources := methodokta.GroupReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(context.Background()).Execute()
	if err != nil {
		return &methodokta.GroupReport{}, err
	}

	// Fetch all Groups
	groups, resp, err := client.GroupAPI.ListGroups(ctx).Execute()
	if err != nil {
		return &methodokta.GroupReport{}, err
	}
	var allGroups []okta.Group
	allGroups = append(allGroups, groups...)
	for resp.HasNextPage() {
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor := parsedURL.Query().Get("after")
		groups, resp, err = client.GroupAPI.ListGroups(ctx).After(cursor).Execute()
		if err != nil {
			return &methodokta.GroupReport{}, err
		}
		allGroups = append(allGroups, groups...)
		time.Sleep(1000 * time.Millisecond)
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
		var allApps []okta.ListApplications200ResponseInner
		apps, resp, err := client.GroupAPI.ListAssignedApplicationsForGroup(ctx, *g.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			allApps = append(allApps, apps...)
			for resp.HasNextPage() {
				parsedURL, _ := url.Parse(resp.NextPage())
				cursor := parsedURL.Query().Get("after")
				apps, resp, err = client.GroupAPI.ListAssignedApplicationsForGroup(ctx, *g.Id).After(cursor).Execute()
				if err != nil {
					errors = append(errors, err.Error())
				} else {
					allApps = append(allApps, apps...)
				}
			}
		}
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

		// Group Roles
		roles, _, err := client.RoleAssignmentAPI.ListGroupAssignedRoles(ctx, *g.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, r := range roles {
				role := methodokta.RoleInfo{Uid: *r.Id, Name: r.Label}
				group.Roles = append(group.Roles, &role)
			}
		}

		// Group Users
		var allUsers []okta.GroupMember
		users, resp, err := client.GroupAPI.ListGroupUsers(ctx, *g.Id).Execute()
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			allUsers = append(allUsers, users...)
			for resp.HasNextPage() {
				parsedURL, _ := url.Parse(resp.NextPage())
				cursor := parsedURL.Query().Get("after")
				users, resp, err = client.GroupAPI.ListGroupUsers(ctx, *g.Id).After(cursor).Execute()
				if err != nil {
					errors = append(errors, err.Error())
				}
				allUsers = append(allUsers, users...)
			}
		}
		for _, u := range allUsers {
			user := methodokta.UserInfo{Uid: *u.Id, Email: *u.Profile.Email}
			group.Users = append(group.Users, &user)
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
