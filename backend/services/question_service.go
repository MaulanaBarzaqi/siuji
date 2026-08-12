package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"

	"github.com/google/uuid"
)

type QuestionService interface {
	CreateQuestion(sectionPublicID string, req *models.QuestionRequest) (*models.Question, error)
	UpdateQuestion(questionPublicID string, req *models.QuestionRequest) (*models.Question, error)
	DeleteQuestion(questionPublicID string) error
}

type questionService struct {
	questionRepo repositories.QuestionRepository
	sectionRepo repositories.SectionRepository
}

func NewQuestionService(
	questionRepo repositories.QuestionRepository, 
	sectionRepo repositories.SectionRepository,
	) QuestionService {
		return &questionService{
			questionRepo: questionRepo,
			sectionRepo: sectionRepo,
		}
}

func (s *questionService) CreateQuestion(sectionPublicID string, req *models.QuestionRequest) (*models.Question, error) {
	section, err := s.sectionRepo.FindByPublicID(sectionPublicID)
	if err != nil {
		return nil, errors.New("section not found")
	}
	maxPos, err := s.questionRepo.GetMaxPositionInSection(section.ID)
	if err != nil {
		return nil, err
	}
	question := &models.Question{
		PublicID: uuid.New().String(),
		SectionID: section.ID,
		Question: req.Question,
		AudioURL: req.AudioURL,
		ImageURL: req.ImageURL,
		Passage: req.Passage,
		Position: maxPos + 1,
	}
	if err := s.questionRepo.Create(question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionService) UpdateQuestion(questionPublicID string, req *models.QuestionRequest) (*models.Question, error) {
	question, err := s.questionRepo.FindByPublicID(questionPublicID)
	if err != nil {
		return nil, errors.New("question not found")
	}
	question.Question = req.Question
	question.AudioURL = req.AudioURL
	question.ImageURL = req.ImageURL
	question.Passage = req.Passage
	if err := s.questionRepo.Update(question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionService) DeleteQuestion(questionPublicID string) error {
	_, err := s.questionRepo.FindByPublicID(questionPublicID)
	if err != nil {
		return errors.New("question not found")
	}
	return s.questionRepo.Delete(questionPublicID)
}