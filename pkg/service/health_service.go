package service

import (
	"context"
	"fmt"
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
	fmt.Println("Test")
	return nil
}
