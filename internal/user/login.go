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

func EnumerateLogin(ctx context.Context, userFlag string, applicationFlag string, days int, logLimit int, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.LoginReport, error) {
	log := svc1log.FromContext(ctx)

	resources := methodokta.LoginReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.LoginReport{}, err
	}

	// Query parameters
	since := time.Now().AddDate(0, 0, -days)
	limit := int32(1000)
	query := buildSystemLogQuery(userFlag, applicationFlag)
	filter := "eventType eq \"user.authentication.sso\""

	// Fetch System Logs for Application SSO Logins
	log.Info("Total Login Events")
	loginEventCmd := client.SystemLogAPI.ListLogEvents(ctx).Q(query).Filter(filter).Since(since).Limit(limit)
	allLogs, err := fetchLoginEventsWithRetry(ctx, loginEventCmd, logLimit, sleep)
	if err != nil {
		return &methodokta.LoginReport{}, err
	}

	// Loop through Logs to find recent login
	var numLoginMap = make(map[string]int)
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
			numLoginMap[key] = 1
		} else {
			if val.Before(date) {
				recentLoginMap[key] = date
			}
			numLoginMap[key]++
		}
	}
	var loginList []*methodokta.Login
	for key, date := range recentLoginMap {
		parts := strings.Split(key, "|")
		userUID, userEmail, applicationUID, applicationName := parts[0], parts[1], parts[2], parts[3]
		login := methodokta.Login{
			User:        &methodokta.UserInfo{Uid: userUID, Email: userEmail},
			Application: &methodokta.ApplicationInfo{Uid: applicationUID, Name: applicationName},
			Count:       numLoginMap[key],
			TimeFrame:   days,
			Last:        date,
			ScanDate:    time.Now().UTC(),
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

func fetchLoginEventsWithRetry(ctx context.Context, cmd okta.ApiListLogEventsRequest, logLimit int, sleep time.Duration) ([]okta.LogEvent, error) {
	log := svc1log.FromContext(ctx)

	var allLogs []okta.LogEvent
	sleepExp := sleep
	pastCursor := "-1"
	currentCursor := ""
	for pastCursor != currentCursor {
		logs, resp, err := cmd.After(currentCursor).Execute()
		if err != nil {
			log.Info("Login Event", svc1log.SafeParam("sleep", sleepExp))
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
			continue
		}
		sleepExp = sleep
		pastCursor = currentCursor
		parsedURL, _ := url.Parse(resp.NextPage())
		currentCursor = parsedURL.Query().Get("after")

		if logLimit > 0 && len(allLogs)+len(logs) >= logLimit {
			remaining := logLimit - len(allLogs)
			allLogs = append(allLogs, logs[:remaining]...)
			log.Info("Users", svc1log.SafeParam("count", len(allLogs)))
			return allLogs, nil
		}

		allLogs = append(allLogs, logs...)
		log.Info("Login Event", svc1log.SafeParam("count", len(allLogs)))
	}
	return allLogs, nil
}

func retry(sleep time.Duration, err error) bool {
	if err.Error() != "too many requests" {
		return false
	}
	time.Sleep(sleep)
	return true
}

func buildSystemLogQuery(userFlag, applicationFlag string) string {
	query := ""
	if userFlag != "" {
		query += ", " + userFlag
	}
	if applicationFlag != "" {
		query += ", " + applicationFlag
	}
	return query
}
