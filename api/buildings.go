package api

import (
	"database/sql"
	"errors"
	"net/http"
	"rello-test-task/db/models"

	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type BuildingResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// @Summary Get buildings
// @Tags Buildings
// @Description Get all buildings
// @ID get-buildings
// @Produce json
// @Success 200 {object} []BuildingResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /buildings [get]
func (s *Server) GetAllBuildings(ctx *fiber.Ctx) error {
	buildings, err := models.Buildings().All(ctx.Context(), s.store.GetDB())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return err
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	buildingsResponse := mapToBuildingsResponse(buildings)
	ctx.Status(http.StatusOK).JSON(buildingsResponse)

	return nil
}

// @Summary Get building by id
// @Tags building
// @Description Get a specific building by the specified id
// @ID get-building-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific building id"
// @Success 200 {object} BuildingResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /buildings/{id} [get]
func (s *Server) GetBuildingById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	building, err := models.Buildings(qm.Where("id=?", id)).One(ctx.Context(), s.store.GetDB())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return nil
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	buildingResponse := mapToBuildingResponse(building)
	ctx.Status(http.StatusOK).JSON(buildingResponse)

	return nil
}

func mapToBuildingsResponse(buildings models.BuildingSlice) []BuildingResponse {
	responsebuildings := make([]BuildingResponse, 0, len(buildings))

	for _, building := range buildings {
		responsebuildings = append(responsebuildings, mapToBuildingResponse(building))
	}

	return responsebuildings
}

func mapToBuildingResponse(building *models.Building) BuildingResponse {

	BuildingResponse := BuildingResponse{
		ID:      int(building.ID),
		Name:    building.Name,
		Address: building.Address.String,
	}
	return BuildingResponse
}
