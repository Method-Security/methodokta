package config

type OktaData struct {
	Domain   string
	APIToken string
}

type RootFlags struct {
	Quiet    bool
	Verbose  bool
	OktaData OktaData
}
