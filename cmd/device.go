package cmd

import (
	"github.com/method-security/methodokta/internal/device"
	"github.com/spf13/cobra"
)

func (a *MethodOkta) InitDeviceCommand() {
	deviceCmd := &cobra.Command{
		Use:   "device",
		Short: "Audit and command Devices",
		Long:  `Audit and command Devices`,
	}

	enumerateCmd := &cobra.Command{
		Use:   "enumerate",
		Short: "Enumerate Devices",
		Long:  `Enumerate Devices`,
		Run: func(cmd *cobra.Command, args []string) {
			report, err := device.EnumerateDevice(cmd.Context(), a.RootFlags.Limit, a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	deviceCmd.AddCommand(enumerateCmd)
	a.RootCmd.AddCommand(deviceCmd)
}
