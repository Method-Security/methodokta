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

func EnumerateLogin(ctx context.Context, userFlag string, applicationFlag string, oktaConfig *okta.Configuration) (*methodokta.LoginReport, error) {
	resources := methodokta.LoginReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.LoginReport{}, err
	}

	// Fetch System Logs for Application Logins
	q := "user.authentication.sso, "
	if userFlag != "" {
		q += userFlag + ", "
	}
	if applicationFlag != "" {
		q += applicationFlag + ", "
	}
	var allLogs []okta.LogEvent
	logs, resp, err := client.SystemLogAPI.ListLogEvents(ctx).Q(q).Limit(1000).Execute()

	if err != nil {
		return &methodokta.LoginReport{}, err
	}
	allLogs = append(allLogs, logs...)
	for resp.HasNextPage() {
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor := parsedURL.Query().Get("after")
		logs, resp, err = client.SystemLogAPI.ListLogEvents(ctx).Q(q).After(cursor).Limit(1000).Execute()

		if err != nil {
			return &methodokta.LoginReport{}, err
		}
		allLogs = append(allLogs, logs...)
		time.Sleep(1000 * time.Millisecond)
	}

	// Loop through Logs to find recent login
	var recentLoginMap = make(map[string]time.Time)
	for _, l := range allLogs {
		data, _ := l.MarshalJSON()
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		// Get Application data
		target := result["target"]
		application := target.([]interface{})[0]
		applicationUID := application.(map[string]interface{})["id"].(string)
		applicationName := application.(map[string]interface{})["alternateId"].(string)

		// Get User data
		user := result["actor"]
		userUID := user.(map[string]interface{})["id"].(string)
		userEmail := user.(map[string]interface{})["alternateId"].(string)

		// Get date
		dateStr := result["published"].(string)
		date, _ := time.Parse(time.RFC3339, dateStr)

		key := userUID + "|" + userEmail + "|" + applicationUID + "|" + applicationName
		if val, exists := recentLoginMap[key]; !exists {
			recentLoginMap[key] = date
		} else {
			if val.Before(date) {
				recentLoginMap[key] = date
			}
		}
	}
	var loginList []*methodokta.Login
	for key, date := range recentLoginMap {
		parts := strings.Split(key, "|")
		userUID, userEmail, applicationUID, applicationName := parts[0], parts[1], parts[2], parts[3]
		login := methodokta.Login{
			User:        &methodokta.UserInfo{Uid: userUID, Email: userEmail},
			Application: &methodokta.ApplicationInfo{Uid: applicationUID, Name: applicationName},
			Date:        date,
		}
		loginList = append(loginList, &login)
	}

	resources = methodokta.LoginReport{
		Logins: loginList,
		Org:    *org.Id,
		Errors: errors,
	}
	return &resources, nil
}
