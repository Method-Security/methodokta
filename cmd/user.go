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
			report, err := user.EnumerateUser(cmd.Context(), a.OktaConfig)
			if err != nil {
				errorMessage := err.Error()
				a.OutputSignal.ErrorMessage = &errorMessage
				a.OutputSignal.Status = 1
			}
			a.OutputSignal.Content = report
		},
	}

	userCmd.AddCommand(enumerateCmd)
	a.RootCmd.AddCommand(userCmd)
}
