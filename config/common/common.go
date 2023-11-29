package common_config

import (
	"time"

	"github.com/stellarentropy/gravity-assist-common/config"
)

type CommonConfig struct {
	HealthListenAddress string
	HealthListenPort    int
	HealthReadTimeout   time.Duration
	HealthWriteTimeout  time.Duration

	GracefulShutdownTimeout time.Duration
}

var Common = &CommonConfig{
	// region Health
	HealthListenAddress: config.NewEnv("GRAVITY_ASSIST_HEALTH_LISTEN_ADDRESS").
		WithDefault("0.0.0.0").
		WithRequired().
		GetAddress(),

	HealthListenPort: config.NewEnv("GRAVITY_ASSIST_HEALTH_LISTEN_PORT").
		WithDefault("1234").
		WithRequired().
		GetPort(),

	HealthReadTimeout: config.NewEnv("GRAVITY_ASSIST_HEALTH_READ_TIMEOUT").
		WithDefault("60s").
		WithRequired().
		GetDuration(),

	HealthWriteTimeout: config.NewEnv("GRAVITY_ASSIST_HEALTH_WRITE_TIMEOUT").
		WithDefault("60s").
		WithRequired().
		GetDuration(),
	// endregion

	GracefulShutdownTimeout: config.NewEnv("GRAVITY_ASSIST_GRACEFUL_SHUTDOWN_TIMEOUT").
		WithDefault("60s").
		WithRequired().
		GetDuration(),
}
