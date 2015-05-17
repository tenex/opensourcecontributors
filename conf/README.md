# Configuration and Deployment

Prerequisites:
* install nginx
* have systemd (not upstart!)
* have virtualenv

as root:

```sh

# code:
mkdir /srv
cd /srv
git clone git@github.com:hut8/github-contributions .

# venv:
virtualenv venv
. /srv/venv/bin/activate
pip install -r requirements.txt

# perms:
chown -R www-data:www-data .

# nginx:
rm /etc/nginx/sites-available/default
ln -s /srv/conf/github-contributions.nginx.conf /etc/nginx/conf.d/

# uwsgi:
mkdir -p /etc/uwsgi/vassals
ln -s /srv/conf/github-contributions.uwsgi.ini /etc/uwsgi/vassals

# emperor:
ln -s /srv/conf/emperor.ini /etc/uwsgi
ln -s /srv/conf/emperor.uwsgi.service /etc/systemd/system

# profit
systemctl enable emperor.uwsgi
systemctl start emperor.uwsgi
systemctl restart nginx

# that probably didn't work so debug it

```
