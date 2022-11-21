//go:build memdb

package main

import (
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/profile"
)

func initializeProfileDatabase() *profile.DatabaseSession {
	dbsession := profile.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "hotels.json"))
	return dbsession
}
