# Homework Assignment 1 (Individual)

Submission Deadline: October 4, 2022

This assignment comprises two parts. This is an individual assignment to be conducted individually by each student.

## Part 1: Intro to Go

In this part you will solve two short problems as a way to familiarize yourself with the Go programming language. We expect you to already have a basic knowledge of the language. If you're starting from nothing, we highly recommend going through the [Golang tour](http://tour.golang.org/list) before you begin this assignment. Get started by installing Go on your machine or setting up an experiment in CloudLab using the `multi-node-cluster`. In this class we will use the latest stable version 1.18.

### Software 

You will find the code for this assignment in the directory `wordutil`. 
The two problems that you need to solve are in `wordcount.go` and `wordindex.go`. 
The problem in `workindexall.go` is optional and worth a small number of extra credit points. 
You should only add code to places that say `TODO: implement me`. 
Do not change any of the function signatures as our testing framework uses them.

**Q1 - Word Count**: The task is to count occurrences of each word in a given string.

**Q2 - Word Index**: The task is to find the first occurrence of each word in a given string.

**Q3 - Word Index All (optional)**: The task is to find all occurrences of each word in a given string.

### Testing

Our grading uses the tests in `wordcount_test.go` and `wordindex_test.go` provided to you. To test the correctness of your code, run the following:

```
$ cd wordutil
$ go test
```

If all tests pass, you should see the following output:

```
$ go test
PASS
ok      /path/to/wordutil   0.009s
```

## Part 2: Storage Hierarchy Measurements

In this part you will study the memory and storage hierarchy in the data center. You will measure and report the latency, bandwidth, and capacity characteristics of a cluster of machines of your choice in CloudLab. It is recommended that you use one of the following machine types: `c220g1`, `c240g1`, `c220g2`, `c240g2`, `c220g5`, `c240g5` (from CloudLab Wisconsin). 

Present your results in a graph similar to the one we examined in class and comment. You will submit your answers for this assignment as files in the directory `hierarchy`. Add your comments and observations in the text file `README.md` under `hierarchy`. Add your graph as a new pdf or image file. 

#### 1. Local memory 

To measure latency and bandwidth to local memory (DRAM), you can use the `mlc` (Intel Memory Latency Checker) tool available [here](https://www.intel.com/content/www/us/en/developer/articles/tool/intelr-memory-latency-checker.html). 

To find out how much physical memory the machine has, you can use the command `cat /proc/meminfo`. The total amount of memory will be displayed as `MemTotal`.

#### 2. Local disk 

To measure latency and bandwidth to the local disk, you can use the following tools:

- `ioping` for measuring disk latency.
  - Installation:
    ```
    $ sudo apt−get update
    $ sudo apt−get install ioping −y – 
    ```
  - Example usage:
    ```
    $ ioping −c 20 /tmp/
    ```

- `fio` for measuring disk bandwidth.
  - Installation: 
    ```
    $ sudo apt−get update
    $ sudo apt−get install fio −y
    ```
  - Example usage: To measure read bandwidth:
    ```
    $ fio --randrepeat=1 --ioengine=libaio --direct=1 --gtod_reduce=1 --name=fiotest --filename=testfio --bs=4k --iodepth=64 --size=1G --readwrite=read
    ```

To find disk drive specification, including disk capacity, you can use the `hdparm` tool.
  - Installation:

    ```
    $ sudo apt-get update
    $ sudo apt−get install hdparm −y
    ```

  - Example usage: To retrieve information for the disk device `/dev/sda` 

    ```
    $ sudo hdparm -I /dev/sda
    ```

#### 3. Remote memory and disk 

To measure network latency and bandwidth between two machines in the cluster, you can use the following tools:

- `ping` for measuring network latency.
  - Installation:
    ```
    $ sudo apt-get update
    $ sudo apt−get install iputils−ping −y
    ```
  - Example usage: To measure the latency between `node0` and `node1`, run the following on `node0`:
    ```
    $ ping node1
    ```
- `iperf` for measuring network bandwidth. 
  - Installation: 
    ```
    $ sudo apt−get update
    $ sudo apt−get install iperf−y
    ```
  - Example usage: To measure bandwidth between `node0` and `node1`, run iperf in server mode on `node0`:
    ```
    $ iperf −s −i 1 −w 4M -u
    ```
    and run iperf in client mode on `node1`:
    ```
    $ iperf −c node0 −e −i 1 −u −b 10m
    ```
    Note: In the above example, the client is using a targeted bandwidth equal to 10MB/s via the `−b` parameter which may not reach the network's bandwidth limit. You should try various incremental values ​​for the target bandwidth until you reach the limit.

## Point Distribution

| Problem    | Points |
|------------|--------|
| Repo setup | 20     |
| Q1.1       | 20     |
| Q1.2       | 30     |
| Q1.3       | 10     |
| Q2         | 30     |

### Submission

Now you need to submit your assignment. Commit your change and push it to the remote repository by doing the following:

```
$ git commit -am "[you fill me in]"
$ git push -u origin main
```

You may push you code as many times you like, grading and submission time will be based on your last push.

You should receive full credit as long as the tests pass and you did not modify the test files or hardcode the results. 

## How do I git?

Please read this [Git Tutorial](https://docs.github.com/en/get-started/quickstart/hello-world).

The basic git workflow in the shell (assuming you already have a repo set up):

- git pull
- do some work
- git status (shows what has changed)
- git add all files you want to commit
- git commit -m "brief message on your update"
- git push

All programming assignments, require Git for submission.

We are using Github for distributing and collecting your assignments. At the time of seeing this, you should have already [forked](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the [class repository](https://github.com/ucy-coast/cs499-fa22/) into your private repository. Normally, you only need to clone the repository once, and you will have everything you need for this assignment.

```
$ git clone git@github.com:user/repo.git
$ cd hw1
```

Now, you have everything you need for doing the assignment, i.e., instructions and starter code. Git allows you to keep track of the changes you make to the code. For example, if you want to checkpoint your progress, you can commit your changes by running:

```
$ git commit -am 'partial solution to assignment 1'
```

You should do this early and often! You can push your changes to Github after you commit with:

```
$ git push origin main
``` 
