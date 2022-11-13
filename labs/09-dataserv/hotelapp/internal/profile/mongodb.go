//go:build mongodb

package profile

import (
	log "github.com/sirupsen/logrus"

	pb "github.com/ucy-coast/hotel-app/internal/profile/proto"
	"github.com/ucy-coast/hotel-app/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseSession struct {
	MongoSession *mgo.Session
}

func NewDatabaseSession(db_addr string) *DatabaseSession {
	// TODO: Implement me
}

func (db *DatabaseSession) LoadDataFromJsonFile(profilesJsonPath string) {
	util.LoadDataFromJsonFile(db.MongoSession, "profile-db", "hotels", profilesJsonPath)
}

// GetProfiles returns hotel profiles for requested IDs
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*pb.Hotel, error) {
	// TODO: Implement me
}