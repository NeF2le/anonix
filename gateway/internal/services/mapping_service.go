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

func (s *MappingService) UpdateMappingDek(ctx context.Context, req *mapping.UpdateMappingDekRequest) (
	*mapping.UpdateMappingDekResponse, error) {
	resultChan := make(chan *mapping.UpdateMappingDekResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.UpdateMappingDek(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateMappingDek: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) UpdateMappingCrypto(ctx context.Context, req *mapping.UpdateMappingCryptoRequest) (
	*mapping.UpdateMappingCryptoResponse, error) {
	resultChan := make(chan *mapping.UpdateMappingCryptoResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.UpdateMappingCrypto(ctx, req)
		if err != nil {
			return err
		}
		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateMappingCrypto: %w", err)
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

func (s *MappingService) GetMappingByToken(ctx context.Context, req *mapping.GetMappingByTokenRequest) (
	*mapping.GetMappingResponse, error) {
	resultChan := make(chan *mapping.GetMappingResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetMappingByToken(ctx, req)
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

func (s *MappingService) GetKind(ctx context.Context, req *mapping.GetKindRequest) (
	*mapping.GetKindResponse, error) {
	resultChan := make(chan *mapping.GetKindResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetKind(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetKind: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) GetKindByName(ctx context.Context, req *mapping.GetKindByNameRequest) (
	*mapping.GetKindByNameResponse, error) {
	resultChan := make(chan *mapping.GetKindByNameResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetKindByName(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetKindByName: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) ListKinds(ctx context.Context, req *mapping.ListKindsRequest) (
	*mapping.ListKindsResponse, error) {
	resultChan := make(chan *mapping.ListKindsResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.ListKinds(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call ListKinds: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) CreateKind(ctx context.Context, req *mapping.CreateKindRequest) (
	*mapping.CreateKindResponse, error) {
	resultChan := make(chan *mapping.CreateKindResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.CreateKind(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CreateKind: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) UpdateKind(ctx context.Context, req *mapping.UpdateKindRequest) (
	*mapping.UpdateKindResponse, error) {
	resultChan := make(chan *mapping.UpdateKindResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.UpdateKind(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call UpdateKind: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) DeleteKind(ctx context.Context, req *mapping.DeleteKindRequest) (
	*mapping.DeleteKindResponse, error) {
	resultChan := make(chan *mapping.DeleteKindResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.DeleteKind(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call DeleteKind: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) CreateAuditLog(ctx context.Context, req *mapping.CreateAuditLogRequest) (
	*mapping.CreateAuditLogResponse, error) {
	resultChan := make(chan *mapping.CreateAuditLogResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.CreateAuditLog(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call CreateAuditLog: %w", err)
	}

	return <-resultChan, nil
}

func (s *MappingService) GetAuditLogList(ctx context.Context, req *mapping.GetAuditLogListRequest) (
	*mapping.GetAuditLogListResponse, error) {
	resultChan := make(chan *mapping.GetAuditLogListResponse, 1)

	err := callers.Retry(func() error {
		resp, err := s.MappingServiceRepo.GetAuditLogList(ctx, req)
		if err != nil {
			return err
		}

		resultChan <- resp
		return nil
	}, s.MaxRetries, s.BaseDelay)

	if err != nil {
		return nil, fmt.Errorf("couldn't call GetAuditLogList: %w", err)
	}

	return <-resultChan, nil
}
