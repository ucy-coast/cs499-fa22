# Lab: Essential Tools

The purpose of this lab tutorial is to give you hands-on experience with a few essential tools that you will find helpful for this course.

## Prerequisites

### Setting up the Experiment Environment in Cloudlab

For this tutorial, you will be using a small CloudLab cluster. 

Start a new experiment on CloudLab using the `multi-node-cluster` profile in the `UCY-CS499-DC` project, configured with two (2) physical machine nodes. 

## Version Control: Git and GitHub

[Git](https://git-scm.com/) is a distributed version control system for managing source code. Version control is a system for tracking changes to files. As you modify files, the version control system records and saves each change. This allows you to restore a previous version of your code at any time.

[GitHub](https://github.com/) is a code hosting platform for version control and collaboration. It lets you and others work together on projects from anywhere. Go [there](https://github.com/signup) and create an account if you don’t have one. 

A typical workflow involves editing and working on your content in your local repository on your computer, and then sending your changes to the remote repository on GitHub.

### Getting started

Open a remote SSH terminal session to `node0`. Make sure you enable forwarding of the authentication agent connection.

Make sure you have already installed Git on your working machine (`node0`). To check if git is already installed, open the Terminal and hit:

```
git --version
```

To install the latest version of git:

```
sudo apt -y install git
```

### Working locally 

To use Git you will have to setup a repository. You can take an existing directory to make a Git repository, or create an empty directory.

```bash
# create a new directory, and initialize it with git-specific functions
git init my-repo

# change into the `my-repo` directory
cd my-repo

# create the first file in the project
touch README.md

# git isn't aware of the file, stage it
git add README.md

# take a snapshot of the staging area
git commit -m "add README to initial commit"
```

### Hosting your source code on GitHub

You need to [create](https://github.com/new) a repository for your project. **Do not** initialize the repository with a README, .gitignore or License file. This empty repository will await your code.

<figure>
  <p align="center"><img src="figures/github-create-repo-public.png" width="60%"></p>
  <figcaption><p align="center">Figure. Creating a new GitHub repository </p></figcaption>
</figure>

Before you can push commits made on your local branch to a remote repository, you will need to provide the path for the repository you created on GitHub and rename your local branch.

#### Providing the path to the remote repository 

You can push commits made on your local branch to a remote repository identified by a remote URL. A remote URL is Git's fancy way of saying "the place where your code is stored". That URL could be your repository on GitHub, or another user's fork, or even on a completely different server.

You can only push to two types of URL addresses:
- An HTTPS URL like `https://github.com/user/repo.git`. The `https://` URLs are available on all repositories, regardless of visibility. `https://` URLs work even if you are behind a firewall or proxy.
- An SSH URL, like `git@github.com:user/repo.git`. SSH URLs provide access to a Git repository via SSH, a secure protocol. To use these URLs, you must generate an SSH keypair on your computer and add the public key to your account on GitHub.com. For more information, see [Connecting to GitHub with SSH](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh).

You can use the git remote add command to match a remote URL with a name. For example, you'd type the following in the command line:

```
git remote add origin  <REMOTE_URL> 
```

This associates the name origin with the `REMOTE_URL`.

You can use the command `git remote set-url` to [change a remote's URL](https://docs.github.com/en/github/getting-started-with-github/managing-remote-repositories).

To provide the path for the repository you created on GitHub using the SSH URL:

```
git remote add origin git@github.com:YOUR-USERNAME/YOUR-REPOSITORY-NAME.git
```

#### Renaming the default branch

Every Git repository has an initial branch, which is the first branch to be created when a new repository is generated. Historically, the default name for this initial branch was `master`, but `main` is an increasingly popular choice. 

Since the default branch name for new repositories created on GitHub is now `main`, you will have to to rename your local `master` branch to `main`:

```
git branch -M main
```

#### Pushing and pulling changes

You're now setup to push changes to the remote repository:

```
git push -u origin main
```

If you're already setup for push as above, then the following will bring changes down and merge them in:

```
git pull
```

### Cloning an existing repository

If you want to get a copy of an existing Git repository, the command you need is `git clone`.

For example, to get a copy of the course's repository using an HTTPS URL:

```
git clone https://github.com/ucy-coast/cs499-fa22.git
```

## Configuration Management: `parallel-ssh` and Ansible

### Run a command over SSH on a single machine

When the only thing you need to do over the SSH connection is execute a single quick command, you might not want to take the separate actions of connecting and authenticating, running the command, and then disconnecting.

SSH allows users to append the desired command directly to the connection attempt. The command executes, and the connection is closed.

The basic syntax is `ssh user@host "command"`.

For example, you could check the installation status of a package:

```
$ ssh alice@amd198.utah.cloudlab.us "dpkg -l | grep nano"
```

You can also run commands from a local file. The `-s` option of `bash` helps to read the executable command from the standard input:

```
ssh alice@amd198.utah.cloudlab.us 'bash -s' < src/get_host_info.sh
```

If you need to elevate your privileges on the far side of the SSH connection with `sudo`, then force the use of a pseudo-terminal with `-t`. Use this if `sudo` will challenge you for a password. The command looks like this:

```
$ ssh -t alice@amd198.utah.cloudlab.us "sudo apt -y install nano"
```

### Run a command over SSH on multiple machines 

For cluster sizes larger than 2, `parallel-ssh` and `ansible` are very useful tools for running the same commands on multiple machines.

Getting started with these tools is simple. 
1. Choose a machine as your management system and install the tool. Your management node can also be a managed node.
2. Ensure you have passwordless ssh access from your management system to each managed node.
3. Create a hosts file containing an inventory of your nodes.
4. Start using the tool 

We suggest the following:
- For step 1, you use `node0` as the management system.
- For step 2, you use a CloudLab profile that automatically sets up ssh keys for passwordless access. All the CloudLab profiles under `UCY-CS499-DC` meet this requirement.

#### Using `parallel-ssh`

[`parallel-ssh`](https://manpages.org/parallel-ssh) is a program for executing ssh in parallel on a number of hosts. 

On Debian based distributions, you can install `parallel-ssh` using `apt`:

```
$ sudo apt -y install pssh
```

To use `parallel-ssh`, you need to create a text file called hosts file that contains a list of all the hosts that you want to have the command executed on:

```
node0
node1
```

Once you have a hosts file, you can run `parallel-ssh` for all the hosts from this file to execute a command. You can use the `-i` option to display standard output and standard error as each host completes. 

For example, you can run `date` on each host:

```
$ parallel-ssh -i -h pssh_hosts date
```

Sample output:

```
[1] 06:48:49 [SUCCESS] node0
Wed Jul  6 06:48:49 CDT 2022
[2] 06:48:49 [SUCCESS] node1
Wed Jul  6 06:48:49 CDT 2022
```

As another example, you can run `apt` to install `nano` on each host:

```
$ parallel-ssh -i -h pssh_hosts -- sudo apt -y install nano
```

#### Using `ansible`

[Ansible](https://www.ansible.com/) is a modern configuration management tool that facilitates the task of setting up and maintaining remote servers. [Configuration management](https://en.wikipedia.org/wiki/Configuration_management) is an automated method for maintaining computer systems and software in a known, consistent state. 

<figure>
  <p align="center"><img src="figures/ansible-architecture.png" width="60%"></p>
  <figcaption><p align="center">Figure. Ansible Architecture</p></figcaption>
</figure>

On Debian based distributions, you can install `ansible` using `apt`:

```
sudo apt -y install ansible
```

Similarly to `parallel-ssh`, Ansible uses an inventory file, called hosts to determine which nodes you want to manage. This is a plain-text file which lists individual nodes or groups of nodes (e.g. Web servers, Database servers, etc.). 
The default location for the inventory file is /etc/ansible/hosts, but it’s possible to create inventory files in any location that better suits your needs. In this case, you’ll need to provide the path to your custom inventory file with the `-i` parameter when running Ansible commands. Using per-project inventory files is a good practice to minimize the risk of running a command on the wrong group of servers. A simple hosts file you can start off with looks like this:

```
[webservers]
node1
```

This defines a group of nodes we call `webservers` with two specified hosts in it.

We can now immediately launch Ansible to see if our setup works:

```
$ ansible webservers -i ./hosts -m ping
node0 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
node1 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
```

Ansible connects to each individual node in the group of webservers, transmits the required module (here: `ping`), launches the module, and returns the module’s output. Instead of addressing a group of nodes, we can also specify an individual node or a set of nodes with wildcards.

Modules use the available context to determine what actions if any needed to bring the managed host to the desired state and are [idempotent](https://en.wikipedia.org/wiki/Idempotence#Computer_science_meaning), that means if you run the same task again and again, the state of the machine will not change. To find the list of available modules, use `ansible-doc -l`

Instead of shooting off individual ansible commands, we can group these together into so-called playbooks which declare a specific configuration we want to apply to a node. Tasks specified are processed in the order we specify. A playbook is expressed in YAML:

```
---
- name: Configure webserver with nginx
  hosts: webservers
  vars:
    web_root: "{{ ansible_env.HOME }}/static-site"
  tasks:
    - name: install nginx
      apt: 
        name: nginx 
        update_cache: yes
      become: yes
    
    - name: copy the nginx config file
      template:
        src: static_site.cfg.j2
        dest: /etc/nginx/sites-available/static_site.cfg
      become: yes
    
    - name: create symlink
      file:
        src: /etc/nginx/sites-available/static_site.cfg
        dest: /etc/nginx/sites-enabled/default
        state: link
      become: yes

    - name: ensure {{ web_root }} dir exists
      file:
        path: "{{ web_root }}"
        state: directory

    - name: copy index.html
      copy: 
        src: index.html
        dest: "{{ web_root }}/index.html"

    - name: restart nginx
      service: 
        name: nginx 
        state: restarted
      become: yes
```

The format is very readable. 
With `hosts`, we choose specific hosts and/or groups in our inventory to execute against. 
With `vars`, we define variables we want to pass to tasks. These variables can be used within the playbook and within templates as `{{ var }}`. 
With `tasks`, we define the list of tasks we want to execute. Each task hask a name which helps us track playbook progress and a module we want Ansible to invoke.
Finally, for tasks that require root privileges such as installing packages, we use `become` to ask Ansible to activate privilege escalation and run corresponding tasks as `root` user. 

The playbook uses several types of modules, many of which are self explanatory. The template module construct a file’s content using variables. Ansible uses the [Jinja2](http://jinja.pocoo.org/docs/) templating language.

The `ansible-playbook` utility processes the playbook and instructs the nodes to perform the tasks, starting with an implicit invocation of the setup module, which collects system information for Ansible. Tasks are performed top-down and an error causes Ansible to stop processing tasks for that particular node.

```
ansible-playbook -i ./hosts nginx.yml
```

Once the playbook is finished, if you go to your browser and access `node1`'s public hostname or IP address you should see the following page:

```
nginx, configured by Ansible
If you can see this, Ansible successfully installed nginx.
```

Alternatively, you can use a `curl` command to GET a remote resource and have it displayed in the terminal:

```
curl http://node1/index.html
```

**Exercise**: Extend the nginx setup to display the node that the nginx service is running on. For example, when visiting `node1` you should get the following page: 

```
nginx, configured by Ansible
If you can see this, Ansible successfully installed nginx.
Running on node1.
```

## HTTP Benchmarking: `wrk`

Wrk is a modern HTTP benchmarking tool, which measures the latency of your HTTP services at high loads.

Latency refers to the time interval between the moment the request was made (by wrk) and the moment the response was received (from the service). This can be used to simulate the latency a visitor would experience on your site when visiting it using a browser or any other method that sends HTTP requests.

`wrk` is useful for testing any website or application that relies on HTTP, such as:
- Rails and other Ruby applications
- Express and other JavaScript applications
- PHP applications
- Static websites running on web servers
- Sites and applications behind load balancers like Nginx
- Your caching layer

Tests can’t be compared to real users, but they should give you a good estimate of expected latency so you can better plan your infrastructure. Tests can also give you insight into your performance bottlenecks.

You will use `wrk` to benchmark the static website running on the Ngnix web server configured in the previous section. You will run `wrk` on `node0`, which we refer to as the benchmarking machine. You should already be running the nginx web werver on `node1`, which we refer to as the application machine.

<figure>
  <p align="center"><img src="figures/wrk-application-overview.png" width="80%"></p>
  <figcaption><p align="center">Figure. wrk benchmarking</p></figcaption>
</figure>

### Install wrk on benchmarking machine

On Debian based distributions, you can install `wrk` using `apt`:

```
sudo apt update
sudo apt -y install wrk
```

Alternatively, you can build `wrk` from source:

```
git clone https://github.com/wg/wrk.git
cd wrk
make -j
```

Passing flag `-j` to make starts make in parallel mode, which can significantly improve the performance of a build on modern multicore systems.

### Run a wrk benchmark test

The simplest case we could run with wrk is:

```
wrk -t2 -c5 -d5s --timeout 2s http://node1/
```

Which means:

- `-t2`: Use two separate threads
- `-c5`: Open six connections (the first client is zero)
- `-d5s`: Run the test for five seconds
- `--timeout 2s`: Define a two-second timeout
- `--latency`: Print latency statistics  
- `http://node1/`: The target application is listening on `node1`
- Benchmark the `/` path of our application

This can also be described as six users that request our home page repeatedly for five seconds.

The illustration below shows this situation:

<figure>
  <p align="center"><img src="figures/wrk-architecture-structure.png" width="80%"></p>
  <figcaption><p align="center">Figure. wrk benchmarking</p></figcaption>
</figure>

Wait a few seconds for the test to run, and look at the results, which we’ll analyze in the next step.

### Evaluate the output

Output:

```
Running 5s test @ http://node1/
  2 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   193.80us   54.14us   3.30ms   86.37%
    Req/Sec     9.91k   425.67    10.86k    72.55%
  Latency Distribution
     50%  175.00us
     75%  214.00us
     90%  257.00us
     99%  351.00us
  100509 requests in 5.10s, 43.03MB read
Requests/sec:  19707.66
Transfer/sec:      8.44MB
```

- Current configuration summary:

  ```
  Running 5s test @ http://node1/
    2 threads and 5 connections
  ```

  Here we can see a brief summary of our benchmark configuration. The benchmark took 5 seconds, the benchmarked machine hostname is `node1`, and the test used two threads.

- Normal distribution parameters for the latency and req/sec statistics:

  ```
    Thread Stats   Avg      Stdev     Max   +/- Stdev
      Latency   193.80us   54.14us   3.30ms   86.37%
      Req/Sec     9.91k   425.67    10.86k    72.55%
  ```

  This part shows us the normal distribution details for our benchmark - what parameters a Gaussian function would have.

  Benchmarks don’t always have normal distributions, and that’s why these results might be misleading. Therefore always look at the Max and +/- Stdev values. If those values are high, then you can expect that your distribution might have a heavy tail.

- Percentile latency statistics:

  ```
    Latency Distribution
      50%  175.00us
      75%  214.00us
      90%  257.00us
      99%  351.00us
  ```

  Here we can see the 50%, 75%, 90% and 99% latency percentiles. For example, the 99th percentile latency represents the maximum latency, in seconds, for the fastest 99% of requests. Here, nginx processed 99% of requests in less than 351.00 microseconds. 
  
  Since averages can be misleading, latency percentiles can be very useful when interpreting system performance and tail latency. 

- Statistics about the request numbers, transferred data, and throughput:

  ```
    100509 requests in 5.10s, 43.03MB read
  Requests/sec:  19707.66
  Transfer/sec:      8.44MB
  ```

  Here we see that during the time of 5.1 seconds, wrk could do 100509 requests and transfer 43.03MB of data. Combined with simple math (total number of requrests/benchmark duration) we get the result of 19707.66 requests per second.

## Performance Monitoring: `top`, `perf`

### `top`

The `top` (table of processes) command shows a real-time view of running processes in Linux and displays kernel-managed tasks. The command also provides a system information summary that shows resource utilization, including CPU and memory usage.

To run the top command, type top in the command line and press Enter. The command starts in interactive command mode, showing the active processes and other system information. To quit top, press `q`.

```
top - 06:33:06 up  6:33,  1 user,  load average: 0.00, 0.02, 0.01
Tasks: 474 total,   1 running, 261 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
KiB Mem : 19670142+total, 19374380+free,   988212 used,  1969400 buff/cache
KiB Swap:  3145724 total,  3145724 free,        0 used. 19446206+avail Mem 

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND   
 7891 alice     20   0   42240   4316   3212 R   1.0  0.0   0:00.58 top      
 2545 root      20   0 3564444  53976  29628 S   0.7  0.0   2:47.30 containerd
 6447 root      20   0       0      0      0 I   0.3  0.0   0:08.81 kworker/0:4
 7880 alice     20   0  107992   3756   2744 S   0.3  0.0   0:00.01 sshd     
    1 root      20   0   77832   8884   6528 S   0.0  0.0   0:07.75 systemd            
    2 root      20   0       0      0      0 S   0.0  0.0   0:00.01 kthreadd    
    4 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 kworker/0:0H
    7 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 mm_percpu_wq
    8 root      20   0       0      0      0 S   0.0  0.0   0:00.00 ksoftirqd/0 
    9 root      20   0       0      0      0 I   0.0  0.0   0:00.92 rcu_sched                  
```

Let's now use `top` to observe `nginx` activity. We will run the `wrk` workload generator on `node0` for sixty seconds, which should give us enough time to observe process activity on `node1` using top.

Run `wrk` on `node0`:

```
wrk -t2 -c5 -d60s --timeout 2s http://node1/
```

Run `top` on `node1`:

```
top - 06:38:14 up  6:38,  1 user,  load average: 0.35, 0.08, 0.03
Tasks: 474 total,   3 running, 258 sleeping,   0 stopped,   0 zombie
%Cpu(s):  1.7 us,  2.7 sy,  0.0 ni, 94.7 id,  0.0 wa,  0.0 hi,  0.9 si,  0.0 st
KiB Mem : 19670142+total, 19372916+free,   969588 used,  2002672 buff/cache
KiB Swap:  3145724 total,  3145724 free,        0 used. 19447926+avail Mem 

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND   
 6789 www-data  20   0  143800   7124   5296 R  91.7  0.0   0:54.95 nginx    
 6791 www-data  20   0  143800   7072   5248 S  73.9  0.0   0:44.65 nginx    
 6794 www-data  20   0  143800   7072   5248 R  20.8  0.0   0:13.66 nginx    
 6793 www-data  20   0  143800   6452   4732 S   1.3  0.0   0:00.55 nginx    
 2545 root      20   0 3564444  53976  29628 S   1.0  0.0   2:49.56 containerd
 7930 alice     20   0   42240   4224   3120 R   1.0  0.0   0:00.20 top      
    9 root      20   0       0      0      0 I   0.3  0.0   0:00.99 rcu_sched   
 7880 alice     20   0  107992   3756   2744 S   0.3  0.0   0:00.07 sshd     
    1 root      20   0   77832   8884   6528 S   0.0  0.0   0:07.77 systemd     
    2 root      20   0       0      0      0 S   0.0  0.0   0:00.01 kthreadd    
    4 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 kworker/0:0H
    7 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 mm_percpu_wq
    8 root      20   0       0      0      0 S   0.0  0.0   0:00.00 ksoftirqd/0         
```

Here, we can see that Nginx uses four worker processes simultaneously to handle incoming requests. In Nginx, “auto” is the default number of worker processes, which lets Nginx to automatically adjust the number of worker processes based on available cores. 

### `perf`

The Linux perf command provides support for sampling applications and reading performance counters. perf consists of two parts: the kernel space implementation and the userland tools. 

On Debian based distributions, you can install the `perf` userland tools using `apt`:

```
sudo apt update
sudo apt -y install linux-tools-common linux-tools-generic linux-tools-`uname -r`
```

`perf stat` provides a general performance statistic for a program. You can attach to a running (own) process, monitor a new process or monitor the whole system. The latter is only available for root user, as the performance data can provide hints on the internals of the application.

We can use the following one-liner to find the ID of the process with the highest CPU utilization.

```
ps -eo pid --sort=-%cpu | awk 'NR==2{print $1}'
```

which we can then use to attach perf to the running process:

```
sudo perf stat --pid=`ps -eo pid --sort=-%cpu | awk 'NR==2{print $1}'` sleep 5
```

```
 Performance counter stats for process id '6789':

       4422.324322      task-clock (msec)         #    0.884 CPUs utilized          
              9136      context-switches          #    0.002 M/sec                  
              3553      cpu-migrations            #    0.803 K/sec                  
                56      page-faults               #    0.013 K/sec                  
        3536087548      cycles                    #    0.800 GHz                    
        2498325845      instructions              #    0.71  insn per cycle         
         488483203      branches                  #  110.458 M/sec                  
          10398953      branch-misses             #    2.13% of all branches        

       5.002104997 seconds time elapsed
```

## Text Processing: AWK

AWK is an interpreted programming language. It is very powerful and specially designed for text processing. Its name is derived from the family names of its authors − **Alfred Aho** (2020 ACM A.M. Turing Award Laureate), **Peter Weinberger**, and **Brian Kernighan**.

GAWK is the GNU implementation of the AWK programming language. On Debian based distributions, you can install `gawk` using `apt`:

```
sudo apt update
sudo apt -y install gawk
```

An AWK program consists of a sequence of pattern-action statements and optional function definitions. It processes text files. AWK is a line oriented language. It divides a file into lines called records. Each line is broken up into a sequence of fields. The fields are accessed by special variables: `$1` reads the first field, `$2` the second and so on. The `$0` variable refers to the whole record.

The structure of an AWK program has the following form:

```
pattern { action }
```

The pattern is a test that is performed on each of the records. If the condition is met then the action is performed. Either pattern or action can be omitted, but not both. The default pattern matches each line and the default action is to print the record.

An AWK program can be run in two basic ways: 

1. the program is read from a separate file; the name of the program follows the -f option

    ```
    gawk -f progfile file ...
    ```

2. the program is specified on the command line enclosed by quote characters

    ```
    gawk 'program' file ...
    ```

### Sample commands 

Consider the following text file as the input file for all cases below: 

```
$ cat > employee.txt 
```

```
ajay manager account 45000
sunil clerk account 25000
varun manager sales 50000
amit manager account 47000
tarun peon sales 15000
deepak clerk sales 23000
sunil peon sales 13000
satvik director purchase 80000 
```

#### Default behavior of Awk

By default Awk prints every line of data from the specified file.  

```
$ gawk '{print}' employee.txt
```

Output:  

```
ajay manager account 45000
sunil clerk account 25000
varun manager sales 50000
amit manager account 47000
tarun peon sales 15000
deepak clerk sales 23000
sunil peon sales 13000
satvik director purchase 80000 
```

In the above example, no pattern is given. So the actions are applicable to all the lines. Action print without any argument prints the whole line by default, so it prints all the lines of the file without failure. 

#### Print the lines which match the given pattern

```
$ gawk '/manager/ {print}' employee.txt 
```

Output:  

```
ajay manager account 45000
varun manager sales 50000
amit manager account 47000 
```

In the above example, the awk command prints all the line which matches with the ‘manager’. 

#### Splitting a Line Into Fields

For each record i.e line, the awk command splits the record delimited by whitespace character by default and stores it in the `$n` variables. If the line has 4 words, it will be stored in `$1`, `$2`, `$3` and `$4` respectively. Also, `$0` represents the whole line.  

```
$ gawk '{print $1,$4}' employee.txt 
```

Output:  

```
ajay 45000
sunil 25000
varun 50000
amit 47000
tarun 15000
deepak 23000
sunil 13000
satvik 80000 
```

In the above example, `$1` and `$4` represents Name and Salary fields respectively. 

### Processing `wrk` output

Here we provide you with a simple GAWK script for extracting the latency percentiles from the `wrk` output:

```
BEGIN {
	print "Percentile (%)", "Latency (us)"

	lat2us["us"] = 1
	lat2us["ms"] = 1000
	lat2us["s"] = 10000000
}

/^[[:space:]]+[0-9]+\.?[0-9]*\%/ {
	match ($1, /([0-9]+)\%/, perc)
	match ($2, /([0-9]+\.[0-9]*)([a-z]+)/, lat)
	lat_val = lat[1]
	lat_unit = lat[2]
	print perc[1],  lat_val*lat2us[lat_unit]
}
```

The script uses the `match` function that is specific to GAWK to extract a substring matched by a regular expression.

Running the script on the `wrk` output

```
wrk -t2 -c5 -d5s --timeout 2s http://node1/ | gawk -f wrk_latency.awk
```

produces the following output

```
Percentile (%) Latency (us)
50 175.1
75 214
90 257
99 351
```

## Plotting: Gnuplot

Sometimes it is really nice to just take a quick look at some data. However, when working on remote computers, it is a bit of a burden to move data files to a local computer to create a plot in something like `R`. One solution is to use gnuplot and make a quick plot that is rendered in the terminal.

On Debian based distributions, you can install `gnuplot` using `apt`:

```
sudo apt update
sudo apt -y install gnuplot
```

Here is a very simplified gnuplot code we can use to plot the latency distribution produced by `wrk` and `awk`:

set terminal dumb size 120, 30; set autoscale; plot '-' using 1:2 with lines notitle

Let's break this down:

- `set terminal dumb size 120, 30`: gnuplot has 'terminals', which is essentially the output format for the plot. Here we are using `dumb` which renders the plot in ASCII characters in the terminal. You can also specify size parameters for the plot, in this case we're using `size 120, 30` to make it a bit larger than default (you can play around with this).
- `set autoscale`: This just makes it so that the axes are automatically scaled, which is normally desireable.
- `plot '-' using 1:2 with lines notitle`: Performs the plotting magic. The `'-'` is the file from which to take the data and plot, which is being read in from STDIN (`'-'`), but you could specify a file name instead. Next, using 1:2 tells gnuplot to use columns 1 and 2 from the data file for plotting (x = 1, y = 2). Change accordingly to plot any column combination you desire. Finally, with lines notitle just makes the plot a line plot with no title.

This should allow for many basic plots (but note the lack of axis labels!). This plotting script is called as follows:

```
gnuplot -p -e "set terminal dumb size 120, 30; set autoscale; plot '-' using 1:2 with lines notitle"
```

This basically feeds the script from above into the gnuplot command-line call. The `-p` flag just allows the plot to persist beyond the command call (otherwise it disappears) and `-e` tells gnuplot to expect the following script, which is surrounded by quotes.

```
360 +--------------------------------------------------------------------+   
    |      +      +      +      +      +     +      +      +      +     *|   
340 |-+                                                               **-|   
    |                                                                *   |   
320 |-+                                                            **  +-|   
    |                                                             *      |   
300 |-+                                                          *     +-|   
    |                                                          **        |   
280 |-+                                                       *        +-|   
    |                                                       **           |   
260 |-+                                                  ***           +-|   
    |                                                ****                |   
240 |-+                                          ****                  +-|   
    |                                        ****                        |   
220 |-+                                  ****                          +-|   
    |                             *******                                |   
200 |-+               ************                                     +-|   
    |     ************                                                   |   
180 |*****                                                             +-|   
    |      +      +      +      +      +     +      +      +      +      |   
160 +--------------------------------------------------------------------+   
    50     55     60     65     70     75    80     85     90     95    100  
```