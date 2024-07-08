package api

import (
	"database/sql"
	"errors"
	"net/http"
	"rello-test-task/db/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type createBuildingRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type BuildingResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// @Summary Create or update a building
// @Tags Building
// @Description Create or update a building if already exists
// @ID create-or-update-building
// @Accept json
// @Produce json
// @Param input body createBuildingRequest true "building entity related data"
// @Success 200 {object} BuildingResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /buildings [post]
func (s *Server) CreateBuilding(ctx *fiber.Ctx) error {
	var requestBody createBuildingRequest

	ctx.BodyParser(&requestBody)

	building := models.Building{
		Name:    requestBody.Name,
		Address: null.String{String: requestBody.Address, Valid: true},
	}
	fetchedBuilding, err := models.Buildings(qm.Where("name=?", building.Name)).One(ctx.Context(), s.store.GetDB())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := building.Insert(ctx.Context(), s.store.GetDB(), boil.Infer())
			if err != nil {
				ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
				return nil
			}
		} else {
			ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
			return err
		}
	} else {
		building.ID = fetchedBuilding.ID
		_, err := building.Update(ctx.Context(), s.store.GetDB(), boil.Infer())
		if err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
			return err
		}
	}

	fetchedBuilding, err = models.Buildings(qm.Where("name=?", building.Name)).One(ctx.Context(), s.store.GetDB())
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	buildingResponse := mapToBuildingResponse(fetchedBuilding)
	ctx.Status(http.StatusOK).JSON(buildingResponse)

	return nil
}

// @Summary Get buildings
// @Tags Building
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
// @Tags Building
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

// @Summary Delete building by id
// @Tags Building
// @Description Delete a specific building by the specified id
// @ID delete-building-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific building id"
// @Success 200
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /buildings/{id} [delete]
func (s *Server) DeleteBuildingById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return nil
	}

	building, err := models.FindBuilding(ctx.Context(), s.store.GetDB(), i)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound).JSON(errorResponse(err))
			return nil
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	_, err = building.Delete(ctx.Context(), s.store.GetDB())
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	ctx.Status(http.StatusOK)

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
