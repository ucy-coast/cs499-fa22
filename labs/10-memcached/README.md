# Lab: Caching

This lab tutorial will introduce you to the Memcached memory-caching system. 

## Background

## Memcached

Memcached is a general-purpose distributed memory caching system. It is often used to speed up dynamic database-driven websites by caching data and objects in RAM to reduce the number of times an external data source (such as a database or API) must be read.

Memcached's APIs provide a very large hash table distributed across multiple machines. When the table is full, subsequent inserts cause older data to be purged in least recently used order. Applications using Memcached typically layer requests and additions into RAM before falling back on a slower backing store, such as a database.

## Caching Hotel Map Data with Memcached

In the first part of this lab, we stepped through how to store data in a MongoDB database. In this, the second part of this lab, we're going to learn how to use Memcached to alleviate the pressure of the MongoDB database and improve the response speed of the application.

### Using Memcached as Query Cache 

Memcached provides a simple set of operations (set, get, and delete) that makes it attractive as an elemental component in a large-scale distributed system. 

We will rely on Memcached to lighten the read load on our MongoDB database. In particular, we will use Memcached as a demand-filled look-aside cache as shown in the Figure below. 

<figure>
  <p align="center"><img src="figures/memcached-look-aside-cache.png" width="80%"></p>
  <figcaption><p align="center">Figure. Memcached as a demand-filled look-aside cache. The left half illustrates the read path for a web server on a cache miss. The right half illustrates the write path.</p></figcaption>
</figure>

When a web application needs data, it first requests the value from Memcached by providing a string key. If the item addressed by that key is not cached, the web application retrieves the data from the database or other backend service and populates the cache with the key-value pair. 

For write requests, the web server issues query statements to the database and then sends a delete request to Memcached that invalidates any stale data. Web applications typically choose to delete cached data instead of updating it because deletes are idempotent. Memcached is not the authoritative source of the data and is therefore allowed to evict cached data.

### Memcached Client

As a first step, you will need to extend the `profile` service to interact with Memcached. 

Package [`memcache`](https://pkg.go.dev/github.com/bradfitz/gomemcache/memcache) provides a client for the memcached cache server.

You will need to extend `cmd/profile/mongodb.go` to add a new command-line flag that receives the address of the memcached server:

```go
memcache_addr = flag.String("memc_addr", "memcached-profile:11211", "Address of the memcache server")
```

and pass the address to the `NewDatabaseSession` function:

```go
dbsession := profile.NewDatabaseSession(*database_addr, *memcache_addr)
```

Then, you will need to make the following changes in `internal/profile/mongodb.go`.

Extend the type `DatabaseSession` with a new `MemcClient` field.

```
type DatabaseSession struct {
	MongoSession *mgo.Session
	MemcClient   *memcache.Client
}
```

Second, extend the function `NewDatabaseSession` as follows:

```go
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
```

In this code, you extend to function to:
- Receive an additional parameter `memc_addr` that corresponds to the address of the memcached service.
- Establish a connection to the memcached service.

Extend the function `GetProfiles` as follows:

```go
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
```

In this code, you:
- Query the memcached server for each hotel by calling `Get`. 
- In case of a memcached miss, retrieve the item from MongoDB and write it into Memcached by calling `Set`. 

You will also need to extend the import statement to include the additional packages `encoding/json`, `time`, and `github.com/bradfitz/gomemcache/memcache`.

### Memcached Service

Having completed the implementation of our `profile` service, we will now set up a Memcached service that the `profile` service can use. Depending on how you deploy the new Hotel Map, you can create the MongoDB service either through Docker Compose or Kubernetes.

#### Docker Compose 

We provide you with a `docker-compose.yml` file that defines a new service `memcached-profile`:

```yaml
services:
  frontend:
  ...
  memcached-profile:
    image: memcached
    container_name: 'hotel_app_profile_memc'
    restart: always
    environment:
      - MEMCACHED_CACHE_SIZE=128
      - MEMCACHED_THREADS=2
    logging:
      options:
        max-size: 50m
```

The `memcached-profile` service spins a Memcached server inside a container. 

The Memcached server uses 2 threads and sets aside `MEMCACHED_CACHE_SIZE` MB of RAM for the cache.

To ask Compose to build the `profile` image:

```
docker-compose build profile
```

For the rest of the services, we will use the prebuilt images available through DockerHub.

Define the environment variable `REGISTRY`:

```
export REGISTRY=hvolos01
```

Then, pull the images like so:

```
docker-compose pull frontend search geo rate jaeger mongodb-profile
```

We are now ready to run our app by invoking:

```
docker-compose up
```

To stop containers and remove containers, networks, volumes, and images created by `up`.

```
docker-compose down --volumes
```
