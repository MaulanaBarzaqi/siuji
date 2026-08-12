package controllers

import (
	"siuji-backend/models"
	"siuji-backend/services"
	"siuji-backend/utils"

	"github.com/gofiber/fiber/v3"
)

type QuestionController struct {
	questionService services.QuestionService
}

func NewQuestionController(questionService services.QuestionService) *QuestionController {
	return &QuestionController{
		questionService: questionService,
	}
}

func (c *QuestionController) CreateQuestion(ctx fiber.Ctx) error {
	sectionPublicID := ctx.Params("section_public_id")

	var req models.QuestionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "invalid request body.", err.Error())
	}
	question, err := c.questionService.CreateQuestion(sectionPublicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to create question", err.Error())
	}
	return utils.Success(ctx, "Question created successfully.", question)
}

func (c *QuestionController) UpdateQuestion(ctx fiber.Ctx) error {
	questionPublicID := ctx.Params("question_public_id")

	var req models.QuestionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "invalid request body.", err.Error())
	}
	question, err := c.questionService.UpdateQuestion(questionPublicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to update question.", err.Error())
	}
	return utils.Success(ctx, "Question update successfully.", question)
}

func (c *QuestionController) DeleteQuestion(ctx fiber.Ctx) error {
	questionPublicID := ctx.Params("question_public_id")

	err := c.questionService.DeleteQuestion(questionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to delete question", err.Error())
	}

	return utils.SuccessNoData(ctx, "Question deleted successfully.")
}