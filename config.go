package cddb

type config struct {
	Client string
	User   string
}

var cddbConfig *config

func init() {
	cddbConfig = &config{}
	cddbConfig.Client = "XXXXXXX-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	cddbConfig.User = "XXXXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}
