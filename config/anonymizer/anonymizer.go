package anonymizer_config

import (
	"github.com/stellarentropy/gravity-assist-common/config"
)

type AnonymizerConfig struct {
	ListenAddress string
	ListenPort    int

	AgentAddress string
	AgentPort    int
	AgentSchema  string
}

var Anonymizer = &AnonymizerConfig{
	ListenAddress: config.NewEnv("GRAVITY_ASSIST_ANONYMIZER_LISTEN_ADDRESS").
		WithDefault("0.0.0.0").
		GetAddress(),

	ListenPort: config.NewEnv("GRAVITY_ASSIST_ANONYMIZER_LISTEN_PORT").
		WithDefault("7070").
		GetPort(),

	AgentAddress: config.NewEnv("GRAVITY_ASSIST_AGENT_LISTEN_ADDRESS").
		WithDefault("127.0.0.1").
		GetAddress(),

	AgentPort: config.NewEnv("GRAVITY_ASSIST_AGENT_LISTEN_PORT").
		WithDefault("7071").
		GetPort(),

	AgentSchema: config.NewEnv("GRAVITY_ASSIST_AGENT_SCHEMA").
		WithDefault("http").
		WithOptions("http", "https").
		GetString(),
}
