package device

import (
	"context"
	"net/url"
	"time"

	methodokta "github.com/method-security/methodokta/generated/go"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func EnumerateDevice(ctx context.Context, sleep time.Duration, oktaConfig *okta.Configuration) (*methodokta.DeviceReport, error) {
	resources := methodokta.DeviceReport{}
	errors := []string{}

	client := okta.NewAPIClient(oktaConfig)

	// Org UID
	org, _, err := client.OrgSettingAPI.GetOrgSettings(ctx).Execute()
	if err != nil {
		return &methodokta.DeviceReport{}, err
	}

	// Fetch all Devices
	getDevicescmd := client.DeviceAPI.ListDevices(ctx).Expand("")
	allDevices, err := fetchDevicesWithRetry(getDevicescmd, sleep)
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

		// Device User
		getDeviceUserCmd := client.DeviceAPI.ListDeviceUsers(ctx, *d.Id)
		users, err := fetchDeviceUserssWithRetry(getDeviceUserCmd, sleep)
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

func fetchDevicesWithRetry(cmd okta.ApiListDevicesRequest, sleep time.Duration) ([]okta.DeviceList, error) {
	var allDevices []okta.DeviceList
	sleepExp := sleep
	cursor := ""
	hasNextPage := true
	for hasNextPage {
		devices, resp, err := cmd.After(cursor).Execute()
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
		allDevices = append(allDevices, devices...)
	}
	return allDevices, nil
}

func fetchDeviceUserssWithRetry(cmd okta.ApiListDeviceUsersRequest, sleep time.Duration) ([]okta.DeviceUser, error) {
	sleepExp := sleep
	devices, _, err := cmd.Execute()
	for err != nil {
		devices, _, err = cmd.Execute()
		if err != nil {
			if !retry(sleepExp, err) {
				return nil, err
			}
			sleepExp *= 2
		}
	}
	return devices, nil
}

func retry(sleep time.Duration, err error) bool {
	if err.Error() != "too many requests" {
		return false
	}
	time.Sleep(sleep)
	return true
}
