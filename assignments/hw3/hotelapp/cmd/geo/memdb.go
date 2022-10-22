//go:build memdb

package main

import (
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/geo"
)

func initializeGeoDatabase() *geo.DatabaseSession {
	dbsession := geo.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "geo.json"))
	return dbsession
}
