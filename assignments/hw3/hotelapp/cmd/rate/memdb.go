//go:build memdb

package main

import (
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/rate"
)

func initializeRateDatabase() *rate.DatabaseSession {
	dbsession := rate.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "inventory.json"))
	return dbsession
}
