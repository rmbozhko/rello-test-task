package api

import (
	"database/sql"
	"errors"
	"net/http"
	"rello-test-task/db/models"

	"github.com/gofiber/fiber/v2"
)

type ApartmentResponse struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Floor    int    `json:"floor"`
	SQMeters int    `json:"sq_meters"`
}

// @Summary Get apartments
// @Tags Apartment
// @Description Get all apartments
// @ID get-apartments
// @Produce json
// @Success 200 {object} []ApartmentResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /apartments [get]
func (s *Server) GetAllApartments(ctx *fiber.Ctx) error {
	apartments, err := models.Apartments().All(ctx.Context(), s.store.GetDB())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return err
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	apartmentsResponse := mapToApartmentsResponse(apartments)
	ctx.Status(http.StatusOK).JSON(apartmentsResponse)

	return nil
}

func mapToApartmentsResponse(apartments models.ApartmentSlice) []ApartmentResponse {
	responseapartments := make([]ApartmentResponse, 0, len(apartments))

	for _, apartment := range apartments {
		responseapartments = append(responseapartments, mapToApartmentResponse(apartment))
	}

	return responseapartments
}

func mapToApartmentResponse(apartment *models.Apartment) ApartmentResponse {

	ApartmentResponse := ApartmentResponse{
		ID:       int(apartment.ID),
		Number:   apartment.Number.String,
		Floor:    apartment.Floor.Int,
		SQMeters: apartment.SQMeters.Int,
	}
	return ApartmentResponse
}
