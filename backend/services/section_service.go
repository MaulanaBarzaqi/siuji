package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"

	"github.com/google/uuid"
)

type SectionService interface {
	CreateSection(req *models.SectionRequest) (*models.Section, error)
	GetAllSections(filter, sort string, limit, offset int) ([]models.Section, int64, error)
	GetSectionByPublicID(publicID string) (*models.SectionDetailResponse, error)
	UpdateSection(publicID string, req *models.SectionRequest) (*models.Section, error)
	DeleteSection(publicID string) error

	// answer key & reorder
	UpsertAnswerKey(questionPublicID, optionPublicID string) (*models.AnswerKeyResponse, error)
	ReoderQuestions(questionPublicIDs []string) error
	ReorderOptions(optionPublicIDs []string) error
}

type sectionService struct {
	sectionRepo repositories.SectionRepository
}

func NewSectionService(sectionRepo repositories.SectionRepository) SectionService {
	return &sectionService{sectionRepo: sectionRepo}
}

func(s *sectionService) CreateSection(req *models.SectionRequest) (*models.Section, error) {
	section := &models.Section{
		PublicID: uuid.New().String(),
		Title: req.Title,
	}
	err := s.sectionRepo.Create(section)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func(s *sectionService) GetAllSections(filter, sort string, limit, offset int) ([]models.Section, int64, error) {
	return s.sectionRepo.FindAllPagination(filter, sort, limit, offset)
}

func(s *sectionService) GetSectionByPublicID(publicID string) (*models.SectionDetailResponse, error) {
	section, err := s.sectionRepo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, errors.New("section not found")
	}
	// Mapping ke Response DTO agar struktur questions & options sesuai contract detail section
	var questionResponses []models.QuestionDetailResponse
	for _, q := range section.Questions {
		var optionResponses []models.OptionResponse
		for _, opt := range q.Options {
			optionResponses = append(optionResponses, models.OptionResponse{
				PublicID:   opt.PublicID,
				Label:      opt.Label,
				OptionText: opt.OptionText,
				Position:   opt.Position,
			})
		}

		// Ambil correct_option_public_id dari relation AnswerKey jika ada
		var correctOptPubID *string
		if len(q.AnswerKeys) > 0 && q.AnswerKeys[0].Option.PublicID != "" {
			pubID := q.AnswerKeys[0].Option.PublicID
			correctOptPubID = &pubID
		}

		questionResponses = append(questionResponses, models.QuestionDetailResponse{
			PublicID:              q.PublicID,
			Question:              q.Question,
			AudioURL:              q.AudioURL,
			ImageURL:              q.ImageURL,
			Passage:               q.Passage,
			Position:              q.Position,
			CorrectOptionPublicID: correctOptPubID,
			Options:               optionResponses,
		})
	}

	response := &models.SectionDetailResponse{
		PublicID:  section.PublicID,
		Title:     section.Title,
		Questions: questionResponses,
		CreatedAt: section.CreatedAt,
		UpdatedAt: section.UpdatedAt,
	}

	return response, nil
}

func(s *sectionService) UpdateSection(publicID string, req *models.SectionRequest) (*models.Section, error) {
	section, err := s.sectionRepo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, errors.New("section not found")
	}
	section.Title = req.Title

	err = s.sectionRepo.Update(section)
	if err != nil {
		return nil, err
	}
	return section, nil
} 

func(s *sectionService) DeleteSection(publicID string) error {
	section, err := s.sectionRepo.FindByPublicID(publicID)
	if err != nil {
		return err
	}
	if section == nil {
		return errors.New("section not found")
	}

	return s.sectionRepo.Delete(publicID)
}

func(s *sectionService) UpsertAnswerKey(questionPublicID, optionPublicID string) (*models.AnswerKeyResponse, error) {
	answerKey, err := s.sectionRepo.UpsertAnswerKey(questionPublicID, optionPublicID)
	if err != nil {
		return nil, err
	}

	// Mapping ke response struct API contract nomor 5
	response := &models.AnswerKeyResponse{
		PublicID:              answerKey.PublicID,
		QuestionPublicID:      questionPublicID,
		CorrectOptionPublicID: optionPublicID,
		UpdatedAt:             answerKey.UpdatedAt,
	}

	return response, nil
}

func(s *sectionService) ReoderQuestions(questionPublicIDs []string) error {
	if len(questionPublicIDs) == 0 {
		return errors.New("question public ids cannot be empty")
	}
	return s.sectionRepo.UpdateQuestionPositions(questionPublicIDs)
}

func(s *sectionService) ReorderOptions(optionPublicIDs []string) error {
	if len(optionPublicIDs) == 0 {
		return errors.New("option public ids cannot be empty")
	}
	return s.sectionRepo.UpdateOptionPositions(optionPublicIDs)
}
