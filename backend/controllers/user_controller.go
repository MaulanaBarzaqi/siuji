package controllers

import (
	"siuji-backend/services"
	"siuji-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{userService: s}
}

func (c *UserController) CreateParticipant(ctx fiber.Ctx) error {
	var req services.UserRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	user, err := c.userService.CreateParticipant(req)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create participant", err.Error())
	}

	return utils.Created(ctx, "Participant created successfully", user)
}

func (c *UserController) ImportParticipants(ctx fiber.Ctx) error {
	// Ambil file dari form-data dengan key "file"
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return utils.BadRequest(ctx, "File is required", err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to open uploaded file", err.Error())
	}
	defer file.Close()

	count, err := c.userService.ImportParticipantsFromExcel(file)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to import participants", err.Error())
	}

	return utils.Success(ctx, "Participants imported successfully", fiber.Map{
		"total_imported": count,
	})
}

func (c *UserController) UpdateParticipant(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid participant ID", err.Error())
	}

	var req services.UserRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	user, err := c.userService.UpdateParticipant(uint(id), req)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to update participant", err.Error())
	}

	return utils.Success(ctx, "Participant updated successfully", user)
}

func (c *UserController) GetAllParticipants(ctx fiber.Ctx) error {
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "-created_at")
	
	limitStr := ctx.Query("limit", "10")
	pageStr := ctx.Query("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	users, total, err := c.userService.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch participants", err.Error())
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	meta := utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalDatas: int(total),
		TotalPages: int(totalPages),
		Filter:     filter,
		Sort:       sort,
	}
	return utils.SuccessPagination(ctx, "Participants fetched successfully", users, meta)
}

func (c *UserController) DeleteParticipant(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid participant ID", err.Error())
	}

	err = c.userService.DeleteParticipant(uint(id))
	if err != nil {
		return utils.BadRequest(ctx, "Failed to delete participant", err.Error())
	}

	return utils.Success(ctx, "Participant deleted successfully", nil)
}