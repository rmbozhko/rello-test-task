package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rello-test-task/util"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDatabase := SetupTestDatabase()
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Errorf("failed to get path")
	}
	configPath := filepath.Dir(path) + "../.."
	config, err := util.LoadConfig(configPath)
	testDB, err = sql.Open(config.DBDriver, testDatabase.DbSource)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	defer testDatabase.TearDown()
	os.Exit(m.Run())
}
