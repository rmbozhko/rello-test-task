package db

import (
	"context"
	"database/sql"
	"rello-test-task/db/models"
	"rello-test-task/util"
	"strconv"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestCreateBuilding(t *testing.T) {
	populateDBWithValidRandomBuilding(t)
}

func TestGetAllBuildings(t *testing.T) {
	buildingsTotalNumber := 5
	for i := 0; i < buildingsTotalNumber; i++ {
		populateDBWithValidRandomBuilding(t)
	}

	buildings, err := models.Buildings().All(context.Background(), testDB)

	require.NoError(t, err)
	require.NotEmpty(t, buildings)
	require.True(t, buildingsTotalNumber <= len(buildings))
}

func TestGetBuildingById(t *testing.T) {
	createdBuilding := populateDBWithValidRandomBuilding(t)

	building, err := models.Buildings(qm.Where("id=?", strconv.Itoa(createdBuilding.ID))).One(context.Background(), testDB)

	checkFetchedBuildingIsValid(t, err, *building, createdBuilding)
}

func TestGetBuildingById_NotFound(t *testing.T) {
	building, err := models.Buildings(qm.Where("id=?", strconv.Itoa(-1))).One(context.Background(), testDB)

	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, building)
}

func TestUpdateBuildingById(t *testing.T) {
	createdBuilding := populateDBWithValidRandomBuilding(t)

	createdBuilding.Name = faker.Sentence()
	createdBuilding.Address = null.String{String: faker.Paragraph(), Valid: true}

	rows, err := createdBuilding.Update(context.Background(), testDB, boil.Infer())

	require.Equal(t, 1, int(rows))

	building, err := models.Buildings(qm.Where("id=?", createdBuilding.ID)).One(context.Background(), testDB)

	checkFetchedBuildingIsValid(t, err, *building, createdBuilding)
}

func TestUpdateBuildingById_NotFound(t *testing.T) {
	var building models.Building

	building.ID = -1
	building.Name = faker.Sentence()
	building.Address = null.String{String: faker.Paragraph(), Valid: true}

	rows, err := building.Update(context.Background(), testDB, boil.Infer())
	require.NoError(t, err)
	require.Zero(t, rows)
}

func TestDeleteBuildingById(t *testing.T) {
	createdBuilding := populateDBWithValidRandomBuilding(t)

	fetchedBuilding, fetchedErr := models.FindBuilding(context.Background(), testDB, createdBuilding.ID)

	deleteRows, deleteErr := fetchedBuilding.Delete(context.Background(), testDB)

	checkFetchedBuildingIsValid(t, fetchedErr, *fetchedBuilding, createdBuilding)
	require.NoError(t, deleteErr)
	require.Equal(t, 1, int(deleteRows))
}

func populateDBWithValidRandomBuilding(t *testing.T) models.Building {

	randomBuilding := util.GenerateRandomBuilding()
	var building models.Building
	building.ID = int(randomBuilding.ID)
	building.Address = null.String{String: randomBuilding.Address, Valid: true}
	building.Name = randomBuilding.Name

	err := building.Insert(context.Background(), testDB, boil.Infer())

	require.NoError(t, err)

	fetchedBuilding, err := models.Buildings(qm.Where("id=?", building.ID)).One(context.Background(), testDB)

	checkInsertedBuildingIsValid(t, err, *fetchedBuilding, building)

	return building
}

func checkInsertedBuildingIsValid(t *testing.T, err error, actual models.Building, expected models.Building) {
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.NotZero(t, actual.ID)

	require.Equal(t, actual.Name, expected.Name)
	require.Equal(t, actual.Address, expected.Address)
}

func checkFetchedBuildingIsValid(t *testing.T, err error, actual models.Building, expected models.Building) {
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.NotZero(t, actual.ID)

	require.Equal(t, actual.Name, expected.Name)
	require.Equal(t, actual.Address, expected.Address)
}
