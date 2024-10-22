package cmd

import (
	"github.com/method-security/methodokta/internal/application"
	"github.com/spf13/cobra"
)

func (a *MethodOkta) InitApplicationCommand() {
	applicationCmd := &cobra.Command{
		Use:   "application",
		Short: "Audit and command Applications",
		Long:  `Audit and command Applications`,
	}

	enumerateCmd := &cobra.Command{
		Use:   "enumerate",
		Short: "Enumerate Applications",
		Long:  `Enumerate Applications`,
		Run: func(cmd *cobra.Command, args []string) {
			report, err := application.EnumerateApplication(cmd.Context(), a.RootFlags.Limit, a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	applicationCmd.AddCommand(enumerateCmd)
	a.RootCmd.AddCommand(applicationCmd)
}
