# Configuration and Deployment

* `/github-contributions` is the production code, including process scripts (on master)
* `/srv` is where the application is deployed via the deployment script in `util`

Prerequisites:
* install nginx
* have systemd (not upstart!)
* have virtualenv

as root:

```sh
# code:
mkdir /github-contributions
cd /github-contributions
git clone git@github.com:hut8/github-contributions .

# arch only
rm -rf /srv/http /srv/ftp

# perms:
chown -R www-data:www-data . # or http:http in arch

# nginx - debian:
rm /etc/nginx/sites-available/default
ln -s /github-contributions/conf/github-contributions.nginx.conf /etc/nginx/conf.d/

# nginx - arch
# edit /etc/nginx/nginx.conf to include conf.d/* and delete default site
mkdir /etc/nginx/conf.d
ln -s /github-contributions/conf/github-contributions.nginx.conf /etc/nginx/conf.d/

# uwsgi / emperor:
ln -s /github-contributions/conf/emperor.ini /etc/uwsgi
ln -s /github-contributions/conf/emperor.uwsgi.service /etc/systemd/system # debian
ln -s /github-contributions/conf/emperor.uwsgi.service /etc/systemd/system/multi-user.target.wants # arch
mkdir -p /etc/uwsgi/vassals
ln -s /github-contributions/conf/github-contributions.uwsgi.ini /etc/uwsgi/vassals

# profit
systemctl enable emperor.uwsgi
systemctl start emperor.uwsgi
systemctl restart nginx

# that probably didn't work so debug it
```
