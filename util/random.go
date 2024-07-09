package util

import (
	"math/rand"

	"github.com/go-faker/faker/v4"
)

type RandomBuilding struct {
	ID      int32
	Name    string
	Address string
}

func GenerateRandomBuilding() RandomBuilding {
	return RandomBuilding{
		ID:      rand.Int31(),
		Name:    faker.Sentence(),
		Address: faker.Paragraph(),
	}
}

type RandomApartment struct {
	ID       int32
	Number   string
	Floor    int32
	SQMeters int32
}

func GenerateRandomApartment() RandomApartment {
	return RandomApartment{
		ID:       rand.Int31(),
		Number:   faker.Sentence(),
		Floor:    rand.Int31(),
		SQMeters: rand.Int31(),
	}
}
