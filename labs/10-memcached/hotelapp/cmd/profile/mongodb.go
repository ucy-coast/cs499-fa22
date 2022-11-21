//go:build mongodb

package main

import (
	"flag"
	"path/filepath"

	"github.com/ucy-coast/hotel-app/internal/profile"
)

var (
	database_addr = flag.String("db_addr", "mongodb-profile:27017", "Address of the data base server")
)

func initializeProfileDatabase() *profile.DatabaseSession {
	dbsession := profile.NewDatabaseSession(*database_addr)
	dbsession.LoadDataFromJsonFile(filepath.Join(*jsonDataDir, "hotels.json"))
	return dbsession
}
