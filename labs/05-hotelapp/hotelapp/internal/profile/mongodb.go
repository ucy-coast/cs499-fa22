//go:build mongodb

package profile

import (
	log "github.com/sirupsen/logrus"

	"github.com/ucy-coast/hotel-app/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseSession struct {
	MongoSession *mgo.Session
}

type Hotel struct {
	Id          string   `bson:"id"`
	Name        string   `bson:"name"`
	PhoneNumber string   `bson:"phoneNumber"`
	Description string   `bson:"description"`
	Address     *Address `bson:"address"`
}

type Address struct {
	StreetNumber string  `bson:"streetNumber"`
	StreetName   string  `bson:"streetName"`
	City         string  `bson:"city"`
	State        string  `bson:"state"`
	Country      string  `bson:"country"`
	PostalCode   string  `bson:"postalCode"`
	Lat          float32 `bson:"lat"`
	Lon          float32 `bson:"lon"`
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
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*Hotel, error) {
	session := db.MongoSession.Copy()
	defer session.Close()
	c := session.DB("profile-db").C("hotels")

	hotels := make([]*Hotel, 0)

	for _, id := range hotelIds {
		hotel_prof := new(Hotel)
		err := c.Find(bson.M{"id": id}).One(&hotel_prof)
		if err != nil {
			log.Fatalf("Failed get hotels data: ", err)
		}
		hotels = append(hotels, hotel_prof)
	}
	return hotels, nil
}
