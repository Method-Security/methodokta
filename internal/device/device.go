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
	devices, resp, err := client.DeviceAPI.ListDevices(ctx).Expand("").Execute()
	if err != nil {
		errors = append(errors, err.Error())
		time.Sleep(sleep)
		devices, resp, err = client.DeviceAPI.ListDevices(ctx).Expand("").Execute()
		if err != nil {
			return &methodokta.DeviceReport{}, err
		}
	}
	var allDevices []okta.DeviceList
	allDevices = append(allDevices, devices...)
	for resp.HasNextPage() {
		parsedURL, _ := url.Parse(resp.NextPage())
		cursor := parsedURL.Query().Get("after")
		devices, resp, err = client.DeviceAPI.ListDevices(ctx).After(cursor).Execute()
		if err != nil {
			errors = append(errors, err.Error())
			time.Sleep(sleep)
			devices, resp, err = client.DeviceAPI.ListDevices(ctx).After(cursor).Execute()
			if err != nil {
				return &methodokta.DeviceReport{}, err
			}
		}
		allDevices = append(allDevices, devices...)
	}

	// Loop through Applications
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

		//Device User
		users, _, err := client.DeviceAPI.ListDeviceUsers(ctx, *d.Id).Execute()
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
