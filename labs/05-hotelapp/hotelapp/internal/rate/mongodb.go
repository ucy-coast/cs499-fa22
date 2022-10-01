//go:build mongodb

package rate

import (
	log "github.com/sirupsen/logrus"

	"github.com/ucy-coast/hotel-app/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseSession struct {
	MongoSession *mgo.Session
}

type RatePlan struct {
	HotelId  string    `bson:"hotelId"`
	Code     string    `bson:"code"`
	InDate   string    `bson:"inDate"`
	OutDate  string    `bson:"outDate"`
	RoomType *RoomType `bson:"roomType"`
}

type RoomType struct {
	BookableRate       float64 `bson:"bookableRate"`
	Code               string  `bson:"code"`
	RoomDescription    string  `bson:"roomDescription"`
	TotalRate          float64 `bson:"totalRate"`
	TotalRateInclusive float64 `bson:"totalRateInclusive"`
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

func (db *DatabaseSession) LoadDataFromJsonFile(rateJsonPath string) {
	util.LoadDataFromJsonFile(db.MongoSession, "rate-db", "inventory", rateJsonPath)
}

// GetRates gets rates for hotels for specific date range.
func (db *DatabaseSession) GetRates(hotelIds []string) (RatePlans, error) {
	ratePlans := make(RatePlans, 0)

	session := db.MongoSession.Copy()
	defer session.Close()
	c := session.DB("rate-db").C("inventory")

	for _, hotelID := range hotelIds {
		tmpRatePlans := make(RatePlans, 0)
		err := c.Find(&bson.M{"hotelId": hotelID}).All(&tmpRatePlans)
		if err != nil {
			log.Fatalf("Tried to find hotelId [%v], but got error", hotelID, err.Error())
		} else {
			for _, r := range tmpRatePlans {
				ratePlans = append(ratePlans, r)
			}
		}
	}

	return ratePlans, nil
}
