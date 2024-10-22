package config

type OktaAuth struct {
	APIToken string
	Domain   string
}

type RootFlags struct {
	Quiet    bool
	Verbose  bool
	Limit    int
	OktaAuth OktaAuth
}
