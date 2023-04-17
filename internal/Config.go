package internal

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/mitchellh/mapstructure"
)

// Configuration is the struct which is used to bind the config file to for easier use
type Configuration struct {
	Debug            bool   `config:"Debug"`
	Port             int    `config:"Port"`
	MediaDestination string `config:"MediaDestination"`
	LogDestination   string `config:"LogDestination"`
}

// this will hold the Variable in Runtime.
var configuration Configuration

// ReadConfig This will read the config.yaml in the folder config beside the Application
// After reading it, it will Validate it also.
func ReadConfig(configPath *string) {
	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig = &mapstructure.DecoderConfig{
			TagName: "config",
		}
	}, config.ParseEnv)

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	HandleFatalErrorf(config.LoadFiles(*configPath), "Could not Read configuration")

	configurationValueTemp := Configuration{}
	HandleFatalErrorf(config.BindStruct("", &configurationValueTemp), "Could not bind Struct to Configuration.")
	DebugLogf("Configuration is %s", configurationValueTemp)
	configuration = configurationValueTemp
}
