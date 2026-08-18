package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"
)

type ParticipantExamService interface {
	GetMyPeriods(userID uint) ([]models.ParticipantPeriodListResponse, error)
	GetMyPeriodDetail(periodPublicID string, userID uint) (*models.ParticipantPeriodDetailResponse, error)
}

type participantExamService struct {
	participantRepo repositories.ParticipantRepository
	periodRepo repositories.PeriodRepository
}

func NewPeriodExamService(
	participantRepo repositories.ParticipantRepository, 
	periodRepo repositories.PeriodRepository,
) ParticipantExamService {
	return &participantExamService{
		participantRepo: participantRepo,
		periodRepo: periodRepo,
	}
}

func (s *participantExamService) GetMyPeriods(userID uint) ([]models.ParticipantPeriodListResponse, error) {
	participantPeriods, err := s.participantRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.ParticipantPeriodListResponse
	for _, pp := range participantPeriods {
		responses = append(responses, models.ParticipantPeriodListResponse{
			PublicID: pp.PublicID,
			Period: models.PeriodItemResponse{
				PublicID:  pp.Period.PublicID,
				Title:     pp.Period.Title,
				StartTime: pp.Period.StartTime,
				EndTime:   pp.Period.EndTime,
				Status:    pp.Period.Status,
			},
			Status: pp.Status,
			Score:  pp.Score,
		})
	}

	return responses, nil
}

func (s *participantExamService) GetMyPeriodDetail(periodPublicID string, userID uint) (*models.ParticipantPeriodDetailResponse, error) {
	// Validasi apakah periode dengan public_id tersebut ada
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}
	// Cari relasi participant period berdasarkan internal period ID dan user ID
	participantPeriod, err := s.participantRepo.FindByPeriodIDAndUserID(period.ID, userID)
	if err != nil {
		return nil, err
	}
	if participantPeriod == nil {
		return nil, errors.New("you are not registered in this period")
	}
	// Mapping ke response DTO yang diharapkan
	response := &models.ParticipantPeriodDetailResponse{
		PublicID:          period.PublicID,
		Title:             period.Title,
		Month:             period.Month,
		Year:              period.Year,
		StartTime:         period.StartTime,
		EndTime:           period.EndTime,
		ParticipantStatus: participantPeriod.Status,
	}

	return response, nil
}