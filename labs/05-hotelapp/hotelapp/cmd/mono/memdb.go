//go:build memdb

package main

import (
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/geo"
	"github.com/ucy-coast/hotel-app/internal/profile"
	"github.com/ucy-coast/hotel-app/internal/rate"
)

func initializeProfileDatabase() *profile.DatabaseSession {
	dbsession := profile.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "hotels.json"))
	return dbsession
}

func initializeGeoDatabase() *geo.DatabaseSession {
	dbsession := geo.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "geo.json"))
	return dbsession
}

func initializeRateDatabase() *rate.DatabaseSession {
	dbsession := rate.NewDatabaseSession()
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "inventory.json"))
	return dbsession
}
