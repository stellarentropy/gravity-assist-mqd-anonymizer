package components

import (
	"context"
	"sync"
)

func StartComponents(ctx context.Context, wg *sync.WaitGroup, components ...Component) {
	logger.Info().Msg("starting components")

	for _, component := range components {
		wg.Add(1)
		go func(component Component) {
			defer wg.Done()
			logger.Info().
				Str("component", component.Name()).
				Msg("starting component")
			component.Start(ctx, wg)
		}(component)
	}
}

func StartComponentsAndWait(ctx context.Context, components ...Component) {
	var wg sync.WaitGroup
	StartComponents(ctx, &wg, components...)
	wg.Wait()
}
