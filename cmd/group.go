package cmd

import (
	"github.com/method-security/methodokta/internal/group"
	"github.com/spf13/cobra"
)

func (a *MethodOkta) InitGroupCommand() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Audit and command Groups",
		Long:  `Audit and command Groups`,
	}

	enumerateCmd := &cobra.Command{
		Use:   "enumerate",
		Short: "Enumerate Groups",
		Long:  `Enumerate Groups`,
		Run: func(cmd *cobra.Command, args []string) {
			report, err := group.EnumerateGroup(cmd.Context(), a.RootFlags.Limit, a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	groupCmd.AddCommand(enumerateCmd)
	a.RootCmd.AddCommand(groupCmd)
}
