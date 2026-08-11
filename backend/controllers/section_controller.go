package controllers

import (
	"siuji-backend/models"
	"siuji-backend/services"
	"siuji-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type SectionController struct {
	sectionService services.SectionService
}

func NewSectionController(sectionService services.SectionService) *SectionController {
	return &SectionController{sectionService: sectionService}
}

func (c *SectionController) CreateSection(ctx fiber.Ctx) error {
	var req models.SectionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "invalid request body",err.Error())
	}
	section, err := c.sectionService.CreateSection(&req)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to create section", err.Error())
	}
	return utils.Success(ctx, "Create section retrieved successfully.", section)
}

func (c *SectionController) GetAllSections(ctx fiber.Ctx) error {
	pageStr := ctx.Query("page", "1")
	limitStr := ctx.Query("limit", "10")
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	sections, totalData, err := c.sectionService.GetAllSections(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch sections", err.Error())
	}
	totalPages := int((totalData + int64(limit) - 1) / int64(limit))
	meta := utils.PaginationMeta{
		Page: page,
		Limit: limit,
		TotalDatas: int(totalData),
		TotalPages: totalPages,
		Filter: filter,
		Sort: sort,
	}
	return utils.SuccessPagination(ctx, "List sections retrieved successfully.", sections, meta)
}

func (c *SectionController) GetSectionByPublicID(ctx fiber.Ctx) error {
	publicID := ctx.Params("section_public_id")

	section, err := c.sectionService.GetSectionByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "section not found", err.Error())
	}
	return utils.Success(ctx, "Section detail retrieved successfully.", section)
}

func (c *SectionController) UpdateSection(ctx fiber.Ctx) error {
	publicID := ctx.Params("section_public_id")

	var req models.SectionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "failed to parse request body", err.Error())
	}
	section, err := c.sectionService.UpdateSection(publicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to update section", err.Error())
	}
	return utils.Success(ctx, "Section updated successfully.", section)
}

func (c *SectionController) DeleteSection(ctx fiber.Ctx) error {
	publicID := ctx.Params("section_public_id")

	err := c.sectionService.DeleteSection(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Section not found", err.Error())
	}
	return utils.SuccessNoData(ctx, "Delete section retrieved successfully.")
}

func (c *SectionController) UpsertAnswerKey(ctx fiber.Ctx) error {
	questionPublicID := ctx.Params("question_public_id")

	var req struct {
		OptionPublicID string `json:"option_public_id" validate:"required"`
	}
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body.", err.Error())
	}
	res, err := c.sectionService.UpsertAnswerKey(questionPublicID, req.OptionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to set answer key", err.Error())
	}
	return utils.Success(ctx, "Answer key saved successfully.", res)
}

func (c *SectionController) ReorderOptions(ctx fiber.Ctx) error {
	var req struct {
		OptionPublicIDs []string `json:"option_public_ids" validate:"required"`
	}
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "invalid request body", err.Error())
	}
	err := c.sectionService.ReorderOptions(req.OptionPublicIDs)
	if err != nil {
		return utils.BadRequest(ctx, "failed to update option position", err.Error())
	}
	return utils.SuccessNoData(ctx, "Option positions updated successfully.")
}

func (c *SectionController) ReorderQuestions(ctx fiber.Ctx) error {
	var req struct {
		QuestionPublicIDs []string `json:"question_public_ids" validate:"required"`
	}
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}
	err := c.sectionService.ReoderQuestions(req.QuestionPublicIDs)
	if err != nil {
		return utils.BadRequest(ctx, "failed to update question position", err.Error())
	}
	return utils.SuccessNoData(ctx, "Question positions updated successfully.")
}