package services

import (
	"context"
	"fmt"
	"github.com/NeF2le/anonix/common/callers"
	"github.com/NeF2le/anonix/common/gen/mapping"
	"github.com/NeF2le/anonix/gateway/internal/ports"
	"time"
)

type MappingService struct {
	MappingServiceRepo ports.MappingServiceRepository
	BaseDelay          time.Duration
	MaxRetries         uint
}

func NewMappingService(
	mappingServiceRepo ports.MappingServiceRepository,
	maxRetries uint,
	baseDelay time.Duration) *MappingService {
	return &MappingService{
		MappingServiceRepo: mappingServiceRepo,
		BaseDelay:          baseDelay,
		MaxRetries:         maxRetries,
	}
}

func (s *MappingService) CreateMapping(ctx context.Context, req *mapping.CreateMappingRequest) (
	*mapping.CreateMappingResponse, error) {
	resultChan := make(chan *mapping.CreateMappingResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.CreateMapping(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CreateMapping: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) DeleteMapping(ctx context.Context, req *mapping.DeleteMappingRequest) (
	*mapping.DeleteMappingResponse, error) {
	resultChan := make(chan *mapping.DeleteMappingResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.DeleteMapping(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call DeleteMapping: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) UpdateMapping(ctx context.Context, req *mapping.UpdateMappingRequest) (
	*mapping.UpdateMappingResponse, error) {
	resultChan := make(chan *mapping.UpdateMappingResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.UpdateMapping(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateMapping: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) GetMapping(ctx context.Context, req *mapping.GetMappingRequest) (
	*mapping.GetMappingResponse, error) {
	resultChan := make(chan *mapping.GetMappingResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetMapping(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetMapping: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) GetMappingList(ctx context.Context, req *mapping.GetMappingListRequest) (
	*mapping.GetMappingListResponse, error) {
	resultChan := make(chan *mapping.GetMappingListResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetMappingList(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetMappingByToken: %w", err)
	}

	return <-resultChan, nil
}
