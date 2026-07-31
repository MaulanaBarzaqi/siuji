package utils

import "github.com/gofiber/fiber/v3"

type Response struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        string      `json:"error,omitempty"`
}

type ResponseNoData struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Error        string      `json:"error,omitempty"`
}

type ResponsePaginated struct {
	Status       string         `json:"status"`
	ResponseCode int            `json:"response_code"`
	Message      string         `json:"message,omitempty"`
	Data         interface{}    `json:"data,omitempty"`
	Error        string         `json:"error,omitempty"`
	Meta         PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int    `json:"page" example:"1"`
	Limit      int    `json:"limit" example:"10"`
	TotalDatas int    `json:"total_datas" example:"100"`
	TotalPages int    `json:"total_pages" example:"100"`
	Filter     string `json:"filter" example:"name=maulana"`
	Sort       string `json:"sort" example:"-id"`
}

func Success(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status: "success",
		ResponseCode: fiber.StatusOK,
		Message: message,
		Data: data,
	})
}

func SuccessNoData(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(ResponseNoData{
		Status: "success",
		ResponseCode: fiber.StatusOK,
		Message: message,
	})
}

func SuccessPagination(c fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(ResponsePaginated{
		Status: "success",
		ResponseCode: fiber.StatusOK,
		Message: message,
		Data: data,
		Meta: meta,
	})
}

func NotFoundPagination(c fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusNotFound).JSON(ResponsePaginated{
		Status: "error",
		ResponseCode: fiber.StatusNotFound,
		Message: message,
		Data: data,
		Meta: meta,
	})
}

func Created(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status: "success",
		ResponseCode: fiber.StatusCreated,
		Message: message,
		Data: data,
	})
}

func BadRequest(c fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		ResponseCode: fiber.StatusBadRequest,
		Message: message,
		Error: err,
	})
}

func NotFound(c fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		ResponseCode: fiber.StatusNotFound,
		Message: message,
		Error: err,
	})
}

func Unauthorized(c fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status: "error",
		ResponseCode: fiber.StatusUnauthorized,
		Message: message,
		Error: err,
	})
}

func Forbidden(c fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Status: "error",
		ResponseCode: fiber.StatusForbidden,
		Message: message,
		Error: err,
	})
} 

func InternalServerError(c fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		ResponseCode: fiber.StatusInternalServerError,
		Message: message,
		Error: err,
	})
}


