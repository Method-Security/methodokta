package device

import (
	"context"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func EnumerateDevice(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.DeviceReport, error) {
	log := svc1log.FromContext(ctx)

	resources := methodokta.DeviceReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.DeviceReport{}, err
	}

	// Fetch all Devices
	log.Info("Total Devices")
	getDevicescmd := client.DeviceAPI.ListDevices(ctx).Expand("")
	allDevices, err := fetchDevicesWithRetry(ctx, getDevicescmd, sleep)
	if err != nil {
		return &methodokta.DeviceReport{}, err
	}

	// Loop through Devices
	var deviceList []*methodokta.Device
	for _, d := range allDevices {

		// Device data
		status := *d.Status
		statusEnum, err := methodokta.NewStatusTypeFromString(status)
		if err != nil {
			statusEnum, _ = methodokta.NewStatusTypeFromString("UNKNOWN")
		}
		device := methodokta.Device{
			Uid:          *d.Id,
			Name:         d.Profile.DisplayName,
			Platform:     d.Profile.Platform,
			Manufacturer: d.Profile.Manufacturer,
			Model:        d.Profile.Model,
			OsVersion:    d.Profile.OsVersion,
			Status:       statusEnum,
			Created:      *d.Created,
		}

		log.Info("Users for Device", svc1log.SafeParam("Name", d.Profile.DisplayName))

		// Device User
		getDeviceUserCmd := client.DeviceAPI.ListDeviceUsers(ctx, *d.Id)
		users, err := fetchDeviceUsersWithRetry(ctx, getDeviceUserCmd, sleep)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			for _, u := range users {
				user := methodokta.UserInfo{Uid: *u.User.Id, Email: *u.User.Profile.Email}
				device.Users = append(device.Users, &user)
			}
		}

		deviceList = append(deviceList, &device)
	}

	resources = methodokta.DeviceReport{
		Devices: deviceList,
		Org:     *org.Id,
		Errors:  errors,
	}

	return &resources, nil

}

func fetchDevicesWithRetry(ctx context.Context, cmd okta.ApiListDevicesRequest, sleep time.Duration) ([]okta.DeviceList, error) {
	log := svc1log.FromContext(ctx)

	var allDevices []okta.DeviceList
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		devices, resp, err := cmd.After(cursor).Execute()
		if err != nil {
			log.Info("Devices", svc1log.SafeParam("sleep", sleepExp))
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
		allDevices = append(allDevices, devices...)
		log.Info("Devices", svc1log.SafeParam("count", len(allDevices)))
	}
	return allDevices, nil
}

func fetchDeviceUsersWithRetry(ctx context.Context, cmd okta.ApiListDeviceUsersRequest, sleep time.Duration) ([]okta.DeviceUser, error) {
	log := svc1log.FromContext(ctx)

	sleepExp := sleep
	users, _, err := cmd.Execute()
	for err != nil {
		users, _, err = cmd.Execute()
		if err != nil {
			log.Info("Users", svc1log.SafeParam("sleep", sleepExp))
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	log.Info("Users", svc1log.SafeParam("count", len(users)))
	return users, nil
}

func retry(sleep time.Duration, err error) bool {
	if err.Error() != "too many requests" {
		return false
	}
	time.Sleep(sleep)
	return true
}
