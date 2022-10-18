//go:build mongodb

package geo

import (
	log "github.com/sirupsen/logrus"

	"github.com/hailocab/go-geoindex"
	"github.com/ucy-coast/hotel-app/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseSession struct {
	MongoSession *mgo.Session
}

// point represents a hotels's geo location on map
type point struct {
	Pid  string  `bson:"hotelId"`
	Plat float64 `bson:"lat"`
	Plon float64 `bson:"lon"`
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

func (db *DatabaseSession) LoadDataFromJsonFile(geoJsonPath string) {
	util.LoadDataFromJsonFile(db.MongoSession, "geo-db", "geo", geoJsonPath)
}

// newGeoIndex returns a geo index with points loaded
func (db *DatabaseSession) newGeoIndex() *geoindex.ClusteringIndex {
	s := db.MongoSession.Copy()
	defer s.Close()
	c := s.DB("geo-db").C("geo")

	var points []*point
	err := c.Find(bson.M{}).All(&points)
	if err != nil {
		log.Fatalf("Failed get geo data: ", err)
	}

	// add points to index
	index := geoindex.NewClusteringIndex()
	for _, point := range points {
		index.Add(point)
	}

	return index
}
