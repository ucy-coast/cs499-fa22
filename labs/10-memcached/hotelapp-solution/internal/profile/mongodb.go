//go:build mongodb

package profile

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"

	pb "github.com/ucy-coast/hotel-app/internal/profile/proto"
	"github.com/ucy-coast/hotel-app/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type DatabaseSession struct {
	MongoSession *mgo.Session
	MemcClient   *memcache.Client
}

func NewDatabaseSession(db_addr string, memc_addr string) *DatabaseSession {
	session, err := mgo.Dial(db_addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("New session successfull...")

	memc_client := memcache.New(memc_addr)
	memc_client.Timeout = time.Second * 2
	memc_client.MaxIdleConns = 512

	return &DatabaseSession{
		MongoSession: session,
		MemcClient: memc_client,
	}
}

func (db *DatabaseSession) LoadDataFromJsonFile(profilesJsonPath string) {
	util.LoadDataFromJsonFile(db.MongoSession, "profile-db", "hotels", profilesJsonPath)
}

// GetProfiles returns hotel profiles for requested IDs
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*pb.Hotel, error) {
	hotels := make([]*pb.Hotel, 0)
	for _, id := range hotelIds {
		// first check memcached
		item, err := db.MemcClient.Get(id)
		if err == nil {
			// memcached hit
			log.Infof("Memcached hit: hotel_id == %v\n", id)
			hotel_prof := new(pb.Hotel)
			if err = json.Unmarshal(item.Value, hotel_prof); err != nil {
				log.Warn(err)
			}
			hotels = append(hotels, hotel_prof)
		} else if err == memcache.ErrCacheMiss {
			// memcached miss, set up mongo connection
			log.Infof("Memcached miss: hotel_id == %v\n", id)
			session := db.MongoSession.Copy()
			defer session.Close()
			c := session.DB("profile-db").C("hotels")
		
			hotel_prof := new(pb.Hotel)
			err := c.Find(bson.M{"id": id}).One(&hotel_prof)
			if err != nil {
				log.Fatalf("Failed get hotels data: ", err)
			}
			hotels = append(hotels, hotel_prof)

			prof_json, err := json.Marshal(hotel_prof)
			if err != nil {
				log.Errorf("Failed to marshal hotel [id: %v] with err:", hotel_prof.Id, err)
			}
			memc_str := string(prof_json)

			// write to memcached
			err = db.MemcClient.Set(&memcache.Item{Key: id, Value: []byte(memc_str)})
			if err != nil {
				log.Warn("MMC error: ", err)
			}
		} else {
			log.Errorf("Memcached error = %s\n", err)
			panic(err)
		}	
	}
	return hotels, nil	
}
