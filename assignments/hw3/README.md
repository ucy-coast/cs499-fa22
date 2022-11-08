# Homework Assignment 3 (Group)

Submission Deadline: November 15, 2022

This assignment comprises two parts. This is a group assignment to be conducted in teams of two students.

## Part 1: Implement Hotel Map Microservices

In this part, you will implement the microservices for the Hotel Map application. We provide you with a partial implementation that you can use as a starting point. You can find the code in the directory `hotelapp`. The parts you need to fill in are marked with `TODO` comments. 

You should establish that your web app is running correctly. You can test your application by sending to it search queries through the web interface, for example, using your web browser or the `curl` utility, as described in the [lab notes](https://github.com/ucy-coast/cs499-fa22/blob/main/labs/05-hotelapp/README.md#testing).

## Part 2: Evaluate Hotel Map Microservices

In this part, you will conduct a characterization and analysis of the Hotel Map microservices. 

You will use the [wrk2](https://github.com/giltene/wrk2) HTTP benchmarking tool described in the [lab notes](https://github.com/ucy-coast/cs499-fa22/blob/main/labs/05-hotelapp/README.md#benchmarking). You can find the code for `wrk2` in the directory `hotelapp/wrk2`.

### Single node

Use [Docker Compose](https://docs.docker.com/compose/) to deploy the [monolithic implementation](https://github.com/ucy-coast/cs499-fa22/tree/main/labs/06-docker#deploying-web-applications-with-docker) and the microservices-based implementation on a single machine.

Evaluate the throughput and latency performance of the monolithic and microservices-based implementations when deployed on a single machine.

Using jaeger, [trace the RPC call chains](https://github.com/ucy-coast/cs499-fa22/tree/main/labs/07-compose#tracing-requests) and identify any bottlenecks. 

### Multiple nodes

Use [Docker Swarm](https://docs.docker.com/engine/swarm/swarm-tutorial/) to deploy the microservices-based implementation on multiple machines. Deploy each microservice on a separate machine.

Evaluate the throughput and latency performance of the microservices-based implementation when deployed on multiple machines. Compare its performance to that of the monolithic implementation and the microservices-based implementations when deployed on a single machine.

Using jaeger, [trace the RPC call chains](https://github.com/ucy-coast/cs499-fa22/tree/main/labs/07-compose#tracing-requests) and identify any bottlenecks. 

## Point Distribution

| Problem    | Points |
|------------|--------|
| Q1         | 50     |
| Q2         | 50     |

### Submission

Now you need to submit your assignment. Commit your change and push it to the remote repository by doing the following:

```
$ git commit -am "[you fill me in]"
$ git push -u origin main
```

You may push you code as many times you like, grading and submission time will be based on your last push.
