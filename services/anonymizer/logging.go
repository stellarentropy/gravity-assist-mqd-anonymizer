package anonymizer

import (
	"os"

	"github.com/stellarentropy/gravity-assist-common/logging"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	logger = logging.GetLogger().With().Str("component", "anonymizer").Logger()

	if os.Getenv("DEBUG") == "true" {
		logger = logger.Level(zerolog.DebugLevel)
	}
}
