# Lab: Essential Tools

The purpose of this lab tutorial is to give you hands on experience with a few essential tools.


TODO
Git/Github
Ansible
HTTP benchmarking, wrk
perf stat
top
mpstat
Awk
Gnuplot

Homework: Loadbalancing

### Run a command over SSH on a single machine

What if the only thing you need to do over the SSH connection is execute a single quick command? You might not want to take the separate actions of connecting and authenticating, running the command, and then disconnecting.

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
$ sudo apt -y install ansible
```

Similarly to `parallel-ssh`, Ansible uses an inventory file, called hosts to determine which nodes you want to manage. This is a plain-text file which lists individual nodes or groups of nodes (e.g. Web servers, Database servers, etc.). 
The default location for the inventory file is /etc/ansible/hosts, but it’s possible to create inventory files in any location that better suits your needs. In this case, you’ll need to provide the path to your custom inventory file with the `-i` parameter when running Ansible commands. Using per-project inventory files is a good practice to minimize the risk of running a command on the wrong group of servers. A simple hosts file you can start off with looks like this:

```
[webservers]
node0
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
$ ansible-playbook -i ./hosts.ini nginx.yml
```

Once the playbook is finished, if you go to your browser and access `node0` or `node1`'s public hostname or IP address you should see the following page:

```
nginx, configured by Ansible
If you can see this, Ansible successfully installed nginx.
```

**Exercise**: Extend the nginx setup to display the node that the nginx service is running on. For example, when visiting `node1` you should get the following page: 

```
nginx, configured by Ansible
If you can see this, Ansible successfully installed nginx.
Running on node1.
```