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
	session, err := mgo.Dial(db_addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("New session successfull...")
	
	return &DatabaseSession{
		MongoSession: session,
	}
}

func (db *DatabaseSession) LoadDataFromJsonFile(profilesJsonPath string) {
	util.LoadDataFromJsonFile(db.MongoSession, "profile-db", "hotels", profilesJsonPath)
}

// GetProfiles returns hotel profiles for requested IDs
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*pb.Hotel, error) {
	session := db.MongoSession.Copy()
	defer session.Close()
	c := session.DB("profile-db").C("hotels")
	
	hotels := make([]*pb.Hotel, 0)

	for _, id := range hotelIds {
		hotel_prof := new(pb.Hotel)
		err := c.Find(bson.M{"id": id}).One(&hotel_prof)
		if err != nil {
			log.Fatalf("Failed get hotels data: ", err)
		}
		hotels = append(hotels, hotel_prof)
	}
	return hotels, nil	
}