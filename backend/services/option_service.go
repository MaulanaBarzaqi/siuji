package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"

	"github.com/google/uuid"
)

type OptionService interface {
	CreateOption(questionPublicID string, req *models.OptionRequest) (*models.Option, error)
	UpdateOption(optionPublicID string, req *models.OptionRequest) (*models.Option, error)
	DeleteOption(optionPublicID string) error
}

type optionService struct {
	optionRepo repositories.OptionRepository
	questionRepo repositories.QuestionRepository
}

func NewOptionService(
	optionRepo repositories.OptionRepository, 
	questionRepo repositories.QuestionRepository,
	) OptionService {
		return &optionService{
			optionRepo: optionRepo,
			questionRepo: questionRepo,
		}
}

func(s *optionService) CreateOption(questionPublicID string, req *models.OptionRequest) (*models.Option, error) {
	question, err := s.questionRepo.FindByPublicID(questionPublicID)
	if err != nil {
		return nil, errors.New("question not found")
	}
	maxPos, err := s.optionRepo.GetMaxPositionInQuestion(question.ID)
	if err != nil {
		return nil, err
	}
	option := &models.Option{
		PublicID: uuid.New().String(),
		QuestionID: question.ID,
		Label: req.Label,
		OptionText: req.OptionText,
		Position: maxPos + 1,
	}
	if err := s.optionRepo.Create(option); err != nil {
		return nil, err
	}
	return option, nil
}

func(s *optionService) UpdateOption(optionPublicID string, req *models.OptionRequest) (*models.Option, error) {
	option, err := s.optionRepo.FindByPublicID(optionPublicID)
	if err != nil {
		return nil, errors.New("option not found")
	}
	option.Label = req.Label
	option.OptionText = req.OptionText
	if err := s.optionRepo.Update(option); err != nil {
		return nil, err
	}
	return option, nil
}

func(s *optionService) DeleteOption(optionPublicID string) error {
	_, err := s.optionRepo.FindByPublicID(optionPublicID)
	if err != nil {
		return errors.New("option not found")
	}

	return s.optionRepo.Delete(optionPublicID)
}