package cmd

import (
	"github.com/method-security/methodokta/internal/user"
	"github.com/spf13/cobra"
)

func (a *MethodOkta) InitUserCommand() {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Audit and command Users",
		Long:  `Audit and command Users`,
	}

	enumerateCmd := &cobra.Command{
		Use:   "enumerate",
		Short: "Enumerate Users",
		Long:  `Enumerate Users`,
		Run: func(cmd *cobra.Command, args []string) {
			report, err := user.EnumerateUser(cmd.Context(), a.RootFlags.Limit, a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Gather Recent Login Data",
		Long:  `Get the most recent Login for each User-Application pair for the last 90 days`,
		Run: func(cmd *cobra.Command, args []string) {
			userFlag, err := cmd.Flags().GetString("user")
			if err != nil {
				a.OutputSignal.AddError(err)
				return
			}
			applicationFlag, err := cmd.Flags().GetString("application")
			if err != nil {
				a.OutputSignal.AddError(err)
				return
			}
			daysFlag, err := cmd.Flags().GetInt("days")
			if err != nil {
				a.OutputSignal.AddError(err)
				return
			}

			report, err := user.EnumerateLogin(cmd.Context(), userFlag, applicationFlag, daysFlag, a.RootFlags.Limit,
				a.RequestSleep, a.OktaConfig)
			if err != nil {
				a.OutputSignal.AddError(err)
			}
			a.OutputSignal.Content = report
		},
	}

	loginCmd.Flags().String("user", "", "The User Account UID to gather Login data for (Defaults to all users).")
	loginCmd.Flags().String("application", "", "The Application UID to gather Login data for (Defaults to all applications).")
	loginCmd.Flags().Int("days", 90, "Number representing how many days to look back in the logs")

	userCmd.AddCommand(enumerateCmd)
	userCmd.AddCommand(loginCmd)
	a.RootCmd.AddCommand(userCmd)
}
