package service

import (
	"context"

	rpmonitoring "gitlab-payment.intra.rakuten-it.com/unbreakable/rpay-golang-common/monitoring"
)

type HealthService interface {
	Up(context.Context) error
}

func NewHealthService() HealthService {
	return &healthService{}
}

type healthService struct {
}

func (s *healthService) Up(ctx context.Context) error {
	logger := rpmonitoring.GetTracingLogger(ctx)
	logger.Info().Msg("Test")
	return nil
}
