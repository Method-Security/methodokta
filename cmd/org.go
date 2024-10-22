package cmd

import (
	"github.com/method-security/methodokta/internal/org"
	"github.com/spf13/cobra"
)

func (a *MethodOkta) InitOrgCommand() {
	orgCmd := &cobra.Command{
		Use:   "org",
		Short: "Audit and command Orgs",
		Long:  `Audit and command Orgs`,
	}

	enumerateCmd := &cobra.Command{
		Use:   "enumerate",
		Short: "Enumerate Orgs",
		Long:  `Enumerate Orgs`,
		Run: func(cmd *cobra.Command, args []string) {
			report, err := org.EnumerateOrg(cmd.Context(), a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	orgCmd.AddCommand(enumerateCmd)
	a.RootCmd.AddCommand(orgCmd)
}
