package controllers

import (
	"siuji-backend/models"
	"siuji-backend/services"
	"siuji-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ParticipantController struct {
	service services.ParticipantService
}

func NewParticipantController(s services.ParticipantService) *ParticipantController {
	return &ParticipantController{service: s}
}

func (c *ParticipantController) AddParticipantToPeriod(ctx fiber.Ctx) error {
	// Ambil period_public_id dari parameter URL
	periodPublicID := ctx.Params("period_public_id")
	// Bind JSON request body ke struct ParticipantRequest
	var req models.ParticipantRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}
	// Panggil service untuk memproses penambahan participant
	participantPeriod, err := c.service.AddParticipantToPeriod(periodPublicID, &req)
	if err != nil {
		// Tangani jika periode tidak ditemukan
		if err.Error() == "period not found" {
			return utils.NotFound(ctx, "Period not found", err.Error())
		}
		// Tangani error lainnya (misal: sudah terdaftar)
		return utils.BadRequest(ctx, "Failed to add participant", err.Error())
	}
	// Mapping response data agar sesuai dengan format output yang diinginkan
	responseData := map[string]interface{}{
		"public_id":        participantPeriod.PublicID,
		"period_public_id": periodPublicID,
		"user": map[string]interface{}{
			"public_id":  participantPeriod.User.PublicID.String(),
			"name":       participantPeriod.User.Name,
			"email":      participantPeriod.User.Email,
			"nim":        participantPeriod.User.NIM,
			"university": participantPeriod.User.University,
			"role":       participantPeriod.User.Role,
		},
		"status":     participantPeriod.Status,
		"score":      participantPeriod.Score,
		"created_at": participantPeriod.CreatedAt,
	}
	// Kembalikan response menggunakan utils.Created (Status 201)
	return utils.Created(ctx, "Participant added to period successfully.", responseData)
}

func (c *ParticipantController) ImportParticipant(ctx fiber.Ctx) error {
	// Ambil period_public_id dari parameter URL
	periodPublicID := ctx.Params("period_public_id")
	// Ambil file dari form-data dengan key (misal: "file")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return utils.BadRequest(ctx, "Failed to get uploaded file", err.Error())
	}
	// Panggil service untuk memproses import file Excel
	importResult, err := c.service.ImportParticipantsFromExcel(periodPublicID, fileHeader)
	if err != nil {
		// Tangani jika periode tidak ditemukan
		if err.Error() == "period not found" {
			return utils.NotFound(ctx, "Period not found", err.Error())
		}
		return utils.BadRequest(ctx, "Failed to import participants", err.Error())
	}
	// Kembalikan response sukses menggunakan utils.Created (Status 200/201 sesuai kebutuhan)
	return utils.Created(ctx, "Participants imported successfully.", importResult)
}

func (c *ParticipantController) GetParticipantByPeriod(ctx fiber.Ctx) error {
	periodPublicID := ctx.Params("period_public_id")

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
	participant, totalData, err := c.service.GetParticipantsByPeriod(periodPublicID, filter, sort, limit, offset)
		if err != nil {
		return utils.InternalServerError(ctx, "failed to fetch participant", err.Error())
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
	return utils.SuccessPagination(ctx, "List participant retrieved successfully.", participant, meta)
}

func (c *ParticipantController) GetParticipantDetail(ctx fiber.Ctx) error {
	// Ambil period_public_id dan user_public_id dari parameter URL
	periodPublicID := ctx.Params("period_public_id")
	userPublicID := ctx.Params("user_public_id")
	// Panggil service untuk mengambil detail participant
	participantDetail, err := c.service.GetParticipantDetail(periodPublicID, userPublicID)
	if err != nil {
		// Tangani jika periode atau participant tidak ditemukan
		if err.Error() == "period not found" || err.Error() == "participant not found in this period" {
			return utils.NotFound(ctx, "Data not found", err.Error())
		}
		// Tangani error sistem/database lainnya
		return utils.InternalServerError(ctx, "Failed to retrieve participant detail", err.Error())
	}
	// Kembalikan response sukses menggunakan utils.Success (Status 200)
	// Karena service sudah me-return models.ParticipantResponse, kita bisa langsung memasukannya ke Data
	return utils.Success(ctx, "Participant detail retrieved successfully.", participantDetail)
}

func (c *ParticipantController) UpdateParticipant(ctx fiber.Ctx) error {
	// Ambil ID dari parameter URL
	periodPublicID := ctx.Params("period_public_id")
	userPublicID := ctx.Params("user_public_id")
	// Bind JSON request body ke struct UpdateParticipantRequest
	var req models.UpdateParticipantRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}
	// Panggil service untuk melakukan update
	updatedParticipant, err := c.service.UpdateParticipant(periodPublicID, userPublicID, &req)
	if err != nil {
		// Logika handling error yang spesifik
		if err.Error() == "period not found" || err.Error() == "participant not found in this period" {
			return utils.NotFound(ctx, "Data not found", err.Error())
		}
		return utils.BadRequest(ctx, "Failed to update participant", err.Error())
	}
	// Mapping response sesuai kebutuhan JSON Anda
	// Karena service sudah me-return models.ParticipantResponse yang strukturnya 
	// sudah mirip dengan yang diminta, kita bisa langsung mengirimnya.
	// Jika ada perbedaan field, kita bisa mapping manual di sini atau di service.
	return utils.Success(ctx, "Participant updated successfully.", updatedParticipant)
}

func (c *ParticipantController) RemoveParticipantFromPeriod(ctx fiber.Ctx) error {
	// Ambil ID dari parameter URL
	periodPublicID := ctx.Params("period_public_id")
	userPublicID := ctx.Params("user_public_id")
	// Panggil service untuk menghapus relasi participant dari periode
	err := c.service.RemoveParticipantFromPeriod(periodPublicID, userPublicID)
	if err != nil {
		// Tangani jika periode, user, atau relasi tidak ditemukan
		if err.Error() == "period not found" || err.Error() == "user not found" || err.Error() == "record not found" {
			return utils.NotFound(ctx, "Data not found", err.Error())
		}
		// Tangani error lainnya
		return utils.BadRequest(ctx, "Failed to remove participant from period", err.Error())
	}
	// Kembalikan response sukses menggunakan utils.SuccessNoData (karena tidak ada objek data yang dikembalikan)
	return utils.SuccessNoData(ctx, "Participant removed from period successfully.")
}