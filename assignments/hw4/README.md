# Homework Assignment 4 (Group)

Submission Deadline: November 29, 2022

This assignment comprises two parts. This is a group assignment to be conducted in teams of two students.

## Part 1: Scale Hotel Map Microservices

In this part, you will scale out the Hotel Map application by running multiple microservice instances on a cluster. You can use deploy Hotel Map using either Docker Swarm or Kubernetes. You can use either your microservices container images (from homework 3) or the ones available from DockerHub. You should evaluate how throughput and latency performance scales with increasing number of microservice instances. You can focus on a single microservice of your choice (e.g. search). 

A microservice deployment can have many identical back-end instances serving many client requests. Each backend server has a certain capacity. Load balancing is key for distributing the load from clients across available instances. Load balancing had many benefits and some of them are: (i) Tolerance of failures: if one of your replicas fails, then other servers can serve the request, (ii) Increased Scalability: you can distribute user traffic across many servers increasing the scalability, (iii) Improved throughput: you can improve the throughput of the application by distributing traffic across various backend servers, and (iv) No downside deployment: you can achieve no downtime deployment using rolling deployment techniques.

Hotel Map microservices use gRPC for communication with other microservices. gRPC is one of the most popular modern RPC frameworks for inter-process communication. It's a great choice for microservice architecture. Docker and Kubernetes services provide load-balanced IP Addresses. But, this default load balancing doesn't work out of the box with gRPC. gRPC works on HTTP/2. The TCP connection on the HHTP/2 is long-lived. A single connection can multiplex many requests. This reduces the overhead associated with connection management. But it also means that connection-level load balancing is not very useful. The default load balancing in Docker Swarm and Kubernetes is based on connection-level load balancing. For that reason, default load balancing does not work with gRPC.

Below we provide hints for how to use the two types of load balancing options available in gRPC, proxy and client-side.

### Proxy load balancing

In Proxy load balancing, the client issues RPCs to a Load Balancer (LB) proxy. The LB distributes the RPC call to one of the available backend servers that implement the actual logic for serving the call. The LB keeps track of load on each backend and implements algorithms for distributing load fairly. The clients themselves do not know about the backend servers. Clients can be untrusted. This architecture is typically used for user-facing services where clients from open internet can connect to the servers.

The pros and cons of proxy load balancing include:
- Pros: Clients are not aware of backend instances
- Pros: Helps you work with clients where incoming load cannot be trusted
- Cons: Since the load balancer is in the data path, higher latency is incurred
- Cons: Load balancer throughput may limit scalability

You can use NGINX as a gRPC proxy and load-balancer. For example, consider a gRPC microservice `profile` that listens to port `8081` deployed using Docker. You can create the NGINX service through `docker-compose.yml`:

```
services:

  profile:
    build: .
    image: ${REGISTRY-127.0.0.1:5000}/hotel_app_profile_single_node_memdb
    entrypoint: profile
    ports:
      - "8081"
    restart: always

  nginx:
    image: nginx:1.20.0
    container_name: nginx
    ports:
      - "8581:8581"
    volumes:
      - ./conf/nginx.conf:/etc/nginx/nginx.conf:ro

  ...
```

The NGINX proxy config looks like:

```
upstream profile_server {
  server profile:8081;
}

server {
  listen 8581 http2;
  location / {
    grpc_pass grpc://profile_server;
  }
}
```

The main things that are happening here are that we are defining NGINX to listen on port `8581` and proxy this HTTP2 traffic to our gRPC server defined as `profile_server`. NGINX figures out that this `serviceName:port` combo resolves to more than one instance through Docker DNS. By default, NGINX will round robin over these servers as the requests come in. There is a way to set the load-balancing behavior to do other things, which you can learn about more [here](https://www.nginx.com/faq/what-are-the-load-balancing-algorithms-supported/).

There are a few notable things that need to happen in your `docker-compose.yml` file.

#### Let your containers grow

Make sure you remove any `container_name` from a service you want to scale, otherwise you will get a warning.

This is important because docker will need to name your containers individually when you want to have more than one of them running.

#### Don’t port clash

We need to make sure that if you are mapping ports, you use the correct format. The standard host port mapping in short syntax is `HOST:CONTAINER` which will lead to port clashes when you attempt to spin up more than one container. We will use ephemeral host ports instead.

Instead of:

```
   ports:
     - "8081:8081"
```

Do this:

```
   ports:
     - "8581"
```     

Doing it this way, Docker will auto-”magic”-ly grab unused ports from the host to map to the container and you won’t know what these are ahead of time. You can see what they ended up being after you bring your service up.

#### Get the proxy hooked up

Using the `nginx` service in `docker-compose.yml` plus the `nginx.conf` should be all you need here. Just make sure that you replace the `profile:8081` with your service’s name and port if it is different from the example.

#### Bring it up

After working through the things outlined above, you start your proxy and service up with a certain number of instances.

To scale the service with [Docker Swarm](https://docs.docker.com/engine/swarm/swarm-tutorial/scale-service/), you run the following command:

```
docker service scale profile=3
```

To scale the service with Docker Compose, you need to pass an additional argument `--scale <serviceName>:<number of instances>`.

```
docker-compose up --scale profile=3
```

#### Inspecting containers

You can use the handy command `docker stats` to get a view in your terminal of your containers. This is a nice and quick way to see the running containers’ CPU, memory, and network utilization, but it shows you these live with no history view.

### Client-side load balancing

In Client-side load balancing, the client is aware of many backend servers and chooses one to use for each RPC. If the client wishes it can implement the load balancing algorithms based on load report from the server. For simple deployment, clients can round-robin requests between available servers.

You can do client-side round-robin load-balancing using Kubernetes headless service. This simple load balancing works out of the box with gRPC. The downside is that it does not take into account the load on the server.

#### Configuring headless service

Luckily, Kubernetes allows clients to discover pod IPs through DNS lookups. Usually, when you perform a DNS lookup for a service, the DNS server returns a single IP — the service’s cluster IP. But if you tell Kubernetes you don’t need a cluster IP for your service (you do this by setting the clusterIP field to None in the service specification ), the DNS server will return the pod IPs instead of the single service IP. Instead of returning a single DNS A record, the DNS server will return multiple A records for the service, each pointing to the IP of an individual pod backing the service at that moment. Clients can therefore do a simple DNS A record lookup and get the IPs of all the pods that are part of the service. The client can then use that information to connect to one, many, or all of them.Setting the clusterIP field in a service spec to None makes the service headless, as Kubernetes won’t assign it a cluster IP through which clients could connect to the pods backing it. 

To make a service a headless service, the only field you need to change is to set `.spec.clusterIP` field as `None`.

#### Verifying DNS 

To confirm the DNS of headless service, create a pod with image [tutum/dnsutils](https://hub.docker.com/r/tutum/dnsutils)as:

```
kubectl run dnsutils --image=tutum/dnsutils --command -- sleep infinity
```

and then run the command

```
kubectl exec dnsutils --  nslookup grpc-server-service
```

As you can see headless service resolves into the IP address of all pods connected through service. Contrast this with the output returned for non-headless service.

#### Configuring client

The only change left is on the client application is to point to the headless service with the port of server pods.

## Part 2: Data Storage Services

Extend the `rate` microservice to store data in MongoDB and cache frequently-accessed data in Memcached. 
You will need to fill in code in `docker-compose.yml`, `internal/rate/mongodb.go`, and also generate (and post-process) the corresponding proto buffer.

Evaluate the throughput and latency performance of the MongoDB-based implementation when deployed on single and multiple machines. 

Evaluate the throughput and latency performance of the Memcached-based implementation when deployed on single and multiple machines. 

Compare the performance of the MongoDB implementation and the Memcached-based implementations.


## Point Distribution

| Problem    | Points |
|------------|--------|
| Q1         | 50     |
| Q2         | 60     |

### Submission

Now you need to submit your assignment. Commit your change and push it to the remote repository by doing the following:

```
$ git commit -am "[you fill me in]"
$ git push -u origin main
```

You may push you code as many times you like, grading and submission time will be based on your last push.
