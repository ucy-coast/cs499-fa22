#!/bin/bash
  
J2=~/.local/bin/j2
WEB_ROOT=~/static_site

# install nginx
sudo apt-get update
sudo apt-get install -y nginx

# copy the nginx config file
sudo apt-get update
sudo apt-get install -y python-pip
pip install j2cli
export web_root=${WEB_ROOT}
${J2} static_site.cfg.j2 | sudo tee /etc/nginx/sites-available/static_site.cfg

# create symlink
sudo ln -fs /etc/nginx/sites-available/static_site.cfg /etc/nginx/sites-enabled/default

# ensure ${WEB_ROOT} dir exists
mkdir -p ${WEB_ROOT}

# copy index.html
cp index.html ${WEB_ROOT}/index.html

# restart nginx
sudo service nginx restart
