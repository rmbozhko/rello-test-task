package api

import (
	"database/sql"
	"errors"
	"net/http"
	"rello-test-task/db/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

// @Summary Get apartment by id
// @Tags Apartment
// @Description Get a specific apartment by the specified id
// @ID get-apartment-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific apartment id"
// @Success 200 {object} ApartmentResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /apartments/{id} [get]
func (s *Server) GetApartmentById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	apartment, err := models.Apartments(qm.Where("id=?", id)).One(ctx.Context(), s.store.GetDB())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return nil
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	apartmentResponse := mapToApartmentResponse(apartment)
	ctx.Status(http.StatusOK).JSON(apartmentResponse)

	return nil
}

// @Summary Delete apartment by id
// @Tags Apartment
// @Description Delete a specific apartment by the specified id
// @ID delete-apartment-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific apartment id"
// @Success 200
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /apartments/{id} [delete]
func (s *Server) DeleteApartmentById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return nil
	}

	apartment, err := models.FindApartment(ctx.Context(), s.store.GetDB(), i)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return nil
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	_, err = apartment.Delete(ctx.Context(), s.store.GetDB())
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	ctx.Status(http.StatusOK)

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
