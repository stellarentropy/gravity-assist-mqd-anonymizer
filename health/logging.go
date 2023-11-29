package health

import (
	"github.com/stellarentropy/gravity-assist-common/logging"

	"github.com/rs/zerolog"
)

// logger records events with a preset 'component' field set to 'metrics',
// categorizing them specifically for the metrics component of the application
// using [zerolog.Logger].
var logger zerolog.Logger

// init prepares the [logger] with a pre-configured [zerolog.Logger] and sets
// 'metrics' as the component for event categorization, executed automatically
// upon package import.
func init() {
	logger = logging.GetLogger().With().Str("component", "health").Logger()
}
