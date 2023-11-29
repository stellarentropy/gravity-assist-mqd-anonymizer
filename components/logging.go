package components

import (
	"github.com/stellarentropy/gravity-assist-common/logging"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	logger = logging.GetLogger().With().Str("component", "components").Logger()
}
