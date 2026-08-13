package controllers

import (
	"siuji-backend/models"
	"siuji-backend/services"
	"siuji-backend/utils"

	"github.com/gofiber/fiber/v3"
)

type OptionController struct {
	optionService services.OptionService
}

func NewOptionController(optionService services.OptionService) *OptionController {
	return &OptionController{
		optionService: optionService,
	}
}

func (c *OptionController) CreateOption(ctx fiber.Ctx) error {
	questionPublicID := ctx.Params("question_public_id")

	var req models.OptionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}
	option, err := c.optionService.CreateOption(questionPublicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to create option", err.Error())
	}
	return utils.Success(ctx, "Option created successfully", option)
}

func (c *OptionController) UpdateOption(ctx fiber.Ctx) error {
	optionPublicID := ctx.Params("option_public_id")

	var req models.OptionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body.", err.Error())
	}

	option, err := c.optionService.UpdateOption(optionPublicID, &req)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to update option", err.Error())
	}

	return utils.Success(ctx, "Option updated successfully.", option)
}

func (c *OptionController) DeleteOption(ctx fiber.Ctx) error {
	optionPublicID := ctx.Params("option_public_id")

	err := c.optionService.DeleteOption(optionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to delete option", err.Error())
	}

	return utils.SuccessNoData(ctx, "Option deleted successfully.")
}