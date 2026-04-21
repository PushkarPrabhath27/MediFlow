package jobs

import (
	"context"
	"time"

	"github.com/mediflow/backend/internal/equipment"
	"github.com/mediflow/backend/internal/request"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	requestService   request.Service
	equipmentService equipment.Service
}

func NewWorker(rs request.Service, es equipment.Service) *Worker {
	return &Worker{
		requestService:   rs,
		equipmentService: es,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Info().Msg("Background worker started")

	// 1. Auto-match ticker
	matchTicker := time.NewTicker(5 * time.Minute)
	// 2. Stock check ticker
	stockTicker := time.NewTicker(1 * time.Hour)
	// 3. Timeout check ticker
	timeoutTicker := time.NewTicker(15 * time.Minute)

	defer matchTicker.Stop()
	defer stockTicker.Stop()
	defer timeoutTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Background worker stopping")
			return
		case <-matchTicker.C:
			w.runAutoMatch()
		case <-stockTicker.C:
			w.runStockCheck()
		case <-timeoutTicker.C:
			w.runTimeoutCheck()
		}
	}
}

func (w *Worker) runAutoMatch() {
	log.Debug().Msg("Running auto-match background job")
	// Logic to find pending requests and run matching
}

func (w *Worker) runStockCheck() {
	log.Debug().Msg("Running min-stock background check")
	// Logic to compare current availability vs department_min_stock
}

func (w *Worker) runTimeoutCheck() {
	log.Debug().Msg("Running transit-timeout background check")
	// Logic to find requests in_transit for > 2 hours and alert
}
