# Lab: Data Services

This lab tutorial will introduce you to the MongoDB document-oriented database.

## Background

### MongoDB

[MongoDB](https://www.mongodb.com/basics) is an open source cross-platform document-oriented database system. Its scale-out architecture allows you to meet the increasing demand for your system by adding more nodes to share the load.

MongoDB stores data records as [BSON](https://www.mongodb.com/docs/manual/reference/bson-types/) documents. BSON is a binary representation of JSON documents. Documents are gathered together in collections. A database stores one or more collections of documents.

The document data model maps naturally to objects in application code, making it simple for developers to learn and use.

The document model provides flexibility to work with complex, fast-changing, messy data from numerous sources. It enables developers to quickly deliver new application functionality.

The fields in a JSON document can vary from document to document. Compare that to a traditional relational database table, where adding a field means adding a column to the database table itself and therefore to every record in the database.

Documents can be nested to express hierarchical relationships and to store structures such as arrays.

<figure>
  <p align="center"><img src="figures/mongodb-document.png" width="80%"></p>
  <figcaption><p align="center">Figure. MongoDB Document</p></figcaption>
</figure>

## Storing Hotel Map Data with MongoDB

So far, data for each of the Hotel Map services is stored in JSON flat files under the `data/` directory. In reality, each of the services would store its data in a persistent datastore service, such as MongoDB.

In this section, you will extend the `profile` service in Hotel Map to store its hotel profile data in a MongoDB database backend. 

### MongoDB Client

As a first step, you will need to extend the `profile` service to interact with MongoDB. To help you, we provide you with a partial implementation in `internal/profile/mongodb.go` that you can fill in to accomplish this task.

#### Data type definitions

The profile data for each hotel is stored in the `hotels` collection in the `profile-db` database using the following composite types:

```go
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
```

Since we use the Protocol Buffers language specification to define the data types, we can re-use the `.proto` defined code for the MongoDB calls also.

From the application home directory, run the following command:

```
protoc --go_out=plugins=grpc:. internal/profile/proto/profile.proto
```

Running this command generates the `internal/profile/proto/profile.pb.go` file, which contains all the protocol buffer code to populate, serialize, and retrieve request and response message types.

When you look at the `profile.pb.go` file, you will notice that the structs the protoc generates have json tags for each field, but there aren't bson tags needed by MongoDB.

If you want to use the same data type, you'll have to modify the code generation. You can use something like [gogoprotobuf](https://github.com/gogo/protobuf) which has an extension to add tags. This should give you bson tags in your structs. You could also post-process your generated files, either with regular expressions or something more complicated involving the go abstract syntax tree. For this lab, we will simply post-process the file manually ourselves. Open `profile.pb.go` in your favorite editor, and replace the json tags with bson tags in the `type Hotel struct`, `type Address struct`, and `type Image struct`. For example:

```
type Hotel struct {
	Id          string   `protobuf:"bytes,1,opt,name=id" bson:"id,omitempty"`
	Name        string   `protobuf:"bytes,2,opt,name=name" bson:"name,omitempty"`
	PhoneNumber string   `protobuf:"bytes,3,opt,name=phoneNumber" bson:"phoneNumber,omitempty"`
	Description string   `protobuf:"bytes,4,opt,name=description" bson:"description,omitempty"`
	Address     *Address `protobuf:"bytes,5,opt,name=address" bson:"address,omitempty"`
	Images      []*Image `protobuf:"bytes,6,rep,name=images" bson:"images,omitempty"`
}
```

should become

```
type Hotel struct {
	Id          string   `protobuf:"bytes,1,opt,name=id" bson:"id,omitempty"`
	Name        string   `protobuf:"bytes,2,opt,name=name" bson:"name,omitempty"`
	PhoneNumber string   `protobuf:"bytes,3,opt,name=phoneNumber" bson:"phoneNumber,omitempty"`
	Description string   `protobuf:"bytes,4,opt,name=description" bson:"description,omitempty"`
	Address     *Address `protobuf:"bytes,5,opt,name=address" bson:"address,omitempty"`
	Images      []*Image `protobuf:"bytes,6,rep,name=images" bson:"images,omitempty"`
}
```

#### MongoDB driver

Package [`mgo`](https://pkg.go.dev/gopkg.in/mgo.v2) offers a rich MongoDB driver for Go that you can use to establish a connection to a MongoDB backend server. 

Usage of the driver revolves around the concept of sessions. To get started, add the following code to the `NewDatabaseSession` function in `internal/profile/mongodb.go` to obtain a session. 

```go
session, err := mgo.Dial(db_addr)
if err != nil {
	log.Fatal(err)
}
log.Info("New session successfull...")

return &DatabaseSession{
	MongoSession: session,
}
```

The `Dial` function will establish one or more connections with the cluster of servers defined by the `db_addr` parameter. From then on, the cluster may be queried to retrieve documents.

You will then extend the `GetProfiles` function to return hotel profiles given a list of hotel IDs. 

Add the following code at the beginning of the `GetProfiles` function:

```go
session := db.MongoSession.Copy()
defer session.Close()
c := session.DB("profile-db").C("hotels")
```

In this code, you:
- Create a new session by calling `session.Copy` on the initial session obtained at dial time. This new session will share the same cluster information and connection pool, and may be easily handed into other methods and functions for organizing logic. 
- Use `defer` to close the session when the function exits and put the connection back into the pool.
- Get a collection to execute the query against.

As a last step, add the following code to retrieve the profile data for each hotel.

```go
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
```

In this code, you:
- Query the collection for each hotel by calling `Find`. `Find` prepares a query using the provided document as a filter. Here, we construct a query filter using the `bson.M` type to match documents with a target hotel `id`. The `bson.M` type defines an unordered representation of a BSON document.
- Retrieve the item from the result set into the result parameter by calling `One`. 

### MongoDB Service

Having completed the implementation of our `profile` service, we will now set up a MongoDB service that the `profile` service can use. Depending on how you deploy the new Hotel Map, you can create the MongoDB service either through Docker Compose or Kubernetes.

#### Docker Compose 

We provide you with a `docker-compose.yml` file that defines a new service `mongodb-profile` and a named volume `profile`:

```yaml
services:
  frontend:
  ...
  profile:
    build:
      context: .
      args:
        - DB=mongodb
    image: hotel_app_profile_single_node_mongodb
    container_name: 'hotel_app_profile'
    entrypoint: profile
    ports:
      - "8081:8081"
    depends_on:
      - mongodb-profile
    restart: always
  ...
  mongodb-profile:
    image: mongo:4.4.6
    container_name: 'hotel_app_profile_mongo'
    ports:
      - "27018:27017"
    restart: always
    volumes:
      - profile:/data/db
    ...
volumes:
  profile:
```

The `mongodb-profile` service spins a MongoDB server inside a container. 

The container uses a named volume `profile` to persist the container directory `/data/db`. Named volumes can persist data after we restart or remove a container. Also, they're accessible by other containers. These volumes are created inside the `/var/lib/docker/volume` local host directory.

Having completed the implementation of the `profile` service, the final step is to build the `profile` image:

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

#### Kubernetes (Optional)

We provide you with YAML manifests that define a MongoDB deployment, expose the deployment through a service, and create a persistent volume for use by the MongoDB deployment.

[Persistent volumes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/) exist independently of the pod lifecycle. They supply your Kubernetes pods with storage that persists across multiple deployments, or that can be attached to different pods over time. They do introduce a lot of additional complexity, but we will skip over most of that here (refer to the official documentation if open questions remain).

There are two relevant Kubernetes API objects, [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) and [PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims). 

A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a resource in the cluster just like a node is a cluster resource. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV. This API object captures the details of the implementation of the storage, be that NFS, iSCSI, or a cloud-provider-specific storage system.

A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources and PVCs consume PV resources. Pods can request specific levels of resources (CPU and Memory). Claims can request specific size and access modes (e.g., they can be mounted ReadWriteOnce, ReadOnlyMany or ReadWriteMany, see AccessModes).

A persistentVolumeClaim volume is used to mount a PersistentVolume into a Pod. PersistentVolumeClaims are a way for users to "claim" durable storage (such as a GCE PersistentDisk or an iSCSI volume) without knowing the details of the particular cloud environment.


```
cat profile-persistent-volume.yaml
```

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: profile-pv
spec:
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 1Gi
  storageClassName: profile-storage
  hostPath:
    path: /local/volumes/profile-pv
    type: Directory
```

```
cat profile-persistent-pvc.yaml
```

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: profile-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: profile-storage
  resources:
    requests:
      storage: 1Gi
```

### The `mongo` shell

The `mongo` shell is an interactive JavaScript interface to MongoDB. You can use the `mongo` shell to query and update data as well as perform administrative operations.

For example, to show all documents in collection `hotels` use below listed code:

```
mongo --eval "db.hotels.find().pretty()" profile-db
```
