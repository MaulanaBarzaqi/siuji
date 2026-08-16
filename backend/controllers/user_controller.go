package controllers

import (
	"siuji-backend/services"
	"siuji-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func(c *UserController) GetAllUsers(ctx fiber.Ctx) error {
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
	users, totalData, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to fetch users", err.Error())
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
	return utils.SuccessPagination(ctx, "List user retrieved successfully.", users, meta)
}

func (c *UserController) GetUserByPublicID(ctx fiber.Ctx) error {
	publicID := ctx.Params("user_public_id")

	user, err := c.service.GetUserByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "user not found", err.Error())
	}
	return utils.Success(ctx, "User detail retrieved successfully.", user)
}

func(c *UserController) DeleteUser(ctx fiber.Ctx) error {
	publicID := ctx.Params("user_public_id")

	err := c.service.DeleteUser(publicID)
	if err != nil {
		return utils.NotFound(ctx, "User not found", err.Error())
	}
	return utils.SuccessNoData(ctx, "delete user retrieved successfully.")
}