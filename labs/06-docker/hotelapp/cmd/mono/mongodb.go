//go:build mongodb

package main

import (
	"flag"
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/geo"
	"github.com/ucy-coast/hotel-app/internal/profile"
	"github.com/ucy-coast/hotel-app/internal/rate"
)

var (
	database_addr = flag.String("db_addr", "mongodb-mono:27017", "Address of the data base server")
)

func initializeProfileDatabase() *profile.DatabaseSession {
	dbsession := profile.NewDatabaseSession(*database_addr)
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "hotels.json"))
	return dbsession
}

func initializeGeoDatabase() *geo.DatabaseSession {
	dbsession := geo.NewDatabaseSession(*database_addr)
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "geo.json"))
	return dbsession
}

func initializeRateDatabase() *rate.DatabaseSession {
	dbsession := rate.NewDatabaseSession(*database_addr)
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "inventory.json"))
	return dbsession
}
