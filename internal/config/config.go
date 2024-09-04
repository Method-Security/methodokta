package config

import "time"

type OktaData struct {
	Domain   string
	APIToken string
}

type RootFlags struct {
	Quiet        bool
	Verbose      bool
	RequestSleep time.Duration
	OktaData     OktaData
}
