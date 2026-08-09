package controllers

import (
	"siuji-backend/models"
	"siuji-backend/services"
	"siuji-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type PeriodController struct {
	periodService services.PeriodService
}

func NewPeriodController(periodService services.PeriodService) *PeriodController {
	return &PeriodController{periodService: periodService}
}


func (c *PeriodController) CreatePeriod(ctx fiber.Ctx) error {
	var req models.PeriodRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "failed to parse request body", err.Error())
	}
	period, err := c.periodService.CreatePeriod(&req)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to create period", err.Error())
	}
	return utils.Created(ctx, "Create period retrieved successfully.", period)
}

func (c *PeriodController) GetAllPeriods(ctx fiber.Ctx) error {
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
	periods, totalData, err := c.periodService.GetAllPeriods(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch periods", err.Error())
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
	return utils.SuccessPagination(ctx, "List period retrieved successfully.", periods, meta)
}

func (c *PeriodController) GetPeriodByPublicID(ctx fiber.Ctx) error {
	publicID := ctx.Params("period_public_id")

	period, err := c.periodService.GetPeriodByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Period not found", err.Error())
	}

	return utils.Success(ctx, "Period detail retrieved successfully", period)
}

func (c *PeriodController) UpdatePeriod(ctx fiber.Ctx) error {
	publicID := ctx.Params("period_public_id")

	var req models.PeriodRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	period, err := c.periodService.UpdatePeriod(publicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to update period", err.Error())
	}

	return utils.Success(ctx, "Period updated successfully", period)
}

func (c *PeriodController) DeletePeriod(ctx fiber.Ctx) error {
	publicID := ctx.Params("period_public_id")

	err := c.periodService.DeletePeriod(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Period not found", err.Error())
	}

	return utils.SuccessNoData(ctx, "Delete period retrieved successfully.")
}

func (c *PeriodController) AddSectionToPeriod(ctx fiber.Ctx) error {
	periodPublicID := ctx.Params("period_public_id")

	var req struct {
		SectionPublicID string `json:"section_public_id"`
		Position        int    `json:"position"`
	}

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	periodSection, err := c.periodService.AddSectionToPeriod(periodPublicID, req.SectionPublicID, req.Position)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to assign section to period", err.Error())
	}

	responseData := fiber.Map{
		"public_id":         periodSection.PublicID,
		"period_public_id":  periodPublicID,
		"section_public_id": req.SectionPublicID,
		"title":             periodSection.Section.Title,
		"position":          periodSection.Position,
		"created_at":        periodSection.CreatedAt,
	}

	return utils.Created(ctx, "Section assigned to period successfully", responseData)
}

func (c *PeriodController) RemoveSectionFromPeriod(ctx fiber.Ctx) error {
	periodPublicID := ctx.Params("period_public_id")
	sectionPublicID := ctx.Params("section_public_id")

	err := c.periodService.RemoveSectionFromPeriod(periodPublicID, sectionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to remove section from period", err.Error())
	}

	return utils.SuccessNoData(ctx, "Section removed from period successfully")
}

func (c *PeriodController) ReorderSections(ctx fiber.Ctx) error {
	periodPublicID := ctx.Params("period_public_id")

	var req struct {
		SectionPublicIDs []string `json:"section_public_ids"`
	}

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	err := c.periodService.ReorderSections(periodPublicID, req.SectionPublicIDs)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to reorder sections", err.Error())
	}

	return utils.SuccessNoData(ctx, "Section positions updated successfully")
}