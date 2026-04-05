# Docker running Nginx PHP MySQL

### Installation to Dev
- make sure that you have installed php and nginx on your wsl2 system
- make sure that you have file blog.local.conf in /etc/nginx/sites-available/
```
server {
    listen 80;
    server_name blog.local;
    server_tokens off;
    include /etc/nginx/snippets/certbot.conf;
    location / {
        proxy_set_header  Host $host;
        proxy_set_header  X-Real-IP $remote_addr;
        proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header  X-Forwarded-Host $remote_addr;
        proxy_set_header  X-NginX-Proxy true;
        proxy_pass        http://blog-nginx;
        proxy_ssl_session_reuse off;
        proxy_redirect off;
    }
}
```
- make sure that you have file .env in root directory of project with params:

```
# Gateway
GATEWAY_HOTS_PORT=80
GATEWAY_DOCKER_PORT=80
```

- `docker compose up --build` - MUST BE without errors (except nodeJS)
- `make composer-install`
- create database if not exists with name `my.laravel.blog`
- `make migrate`
- `chmod 777 -R app/storage`
- In `hosts` file add:

```
127.0.0.1 blog.local
```

- open url - http://blog.local:80/

## Installation to production by docker-compose

- go to provision directory and prepare target servers
    - cd pro*
    - make site
    - make authorize
    - make generate-deploy-key
    - make authorize-deploy
    - docker-login
- return to main directory and deploy project to target servers by example command:
    - `make deploy`

- you need to have nginx installed on your parent system, but not in running mode so as not to occupy port 80.
- `cd /var/www/phpqa.ru`
- `docker compose up --build --remove-orphans` - for Debug mode!
- after any changes in Nginx config, you need to restart docker container: `docker compose restart gateway`

#### Optional:

- check current DataBase IP by command:
    - Linux variant: `docker inspect mysql_phpqa | grep IPAddress`
    - Windows variant: `docker inspect mysql_phpqa | findstr IPAddress`
- your config for connection from OS: `DATABASE_URL="mysql://root:root@127.0.0.1:8989/my.laravel.blog?serverVersion=8.0`
- your config for connection from docker (and
  Adminer): `DATABASE_URL="mysql://root:root@172.18.0.2:8989/my.laravel.blog?serverVersion=8.0"`
- Adminer example working url for PG:
    - http://localhost:8080/?pgsql=172.19.0.2&username=root
- enter to container `docker ps -a`,
  `docker exec -ti {{ container name }} /bin/sh`
  `docker exec -ti phpqarud_php-fpm_1 /bin/sh`
- `chmod 644 data/db/mysql/my.cnf`

## make certs

- Make sure that your website is accessible on :80 port
- make sure that your nginx working correctly (`sudo nginx -t`/`sudo service nginx start`)
- /etc/letsencrypt/cli.ini
    - https://stackoverflow.com/questions/61770338/too-many-flags-setting-configurators-installers-authenticators-webroot-ngi
- command for make cert: `sudo certbot --nginx -d pets.phpqa.ru`
- if you have any problem with certs then you may remove certbot from system and install again:
  `sudo apt-get remove certbot`
    - `sudo apt-get install certbot`
    - Ubuntu 20+: `sudo apt-get install python3-certbot-nginx` / Earlier
      Versions `sudo apt-get install python-certbot-nginx`
    - or Ansible variant:
        - `cd /var/www/phpqa.ru`
        - `make down`
        - go to `provisioning` directory and run command: `make site`
- Example nginx config for cert installation:

```
server {
    listen 80;
    server_name phpqa.ru www.phpqa.ru;
    root /var/www/phpqa.ru/app/public;
    index index.php;
    try_files $uri $uri/ /index.php;
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

## Renewing certs UPDATE 2024-05-06
- Open dir for ensure needed sites to renew: `/etc/letsencrypt/live/`
- Prepare command to renew certs:
  - Example command for renewing certs:
    - `sudo certbot --nginx -d phpqa.ru -d www.phpqa.ru`
- Stop all docker containers `cd /var/www/phpqa.ru` / `make down`
- start local nginx `sudo service nginx start`/`sudo nginx -t`
- Make sure that your website is accessible on :80 port
- Run command for renew certs, that was prepared previously
- renewed certs will be stored there: `/etc/letsencrypt/live/`
  - Check the certificate validity date:
    - `cd /etc/letsencrypt/live/`
    - `openssl x509 -enddate -noout -in cert.pem`
- copy all of them to your letsencrypt docker directory
  - check that certs really copied
- stop local nginx `sudo service nginx stop`
- Start all docker containers `cd /var/www/phpqa.ru` / `make up`
- PROFIT
##### OR
- go to production server and run command: `certbot renew`
- or `certbot certonly --noninteractive --agree-tos -d phpqa.ru -d www.phpqa.ru`
- or go to provisioning directory and run command:
  - `make renew-certificates`
- ```
  sudo certbot --nginx -d pets.phpqa.ru
  sudo certbot --nginx -d sveltewar.phpqa.ru
  ```

## This project use the following ports :

| Server          | Port | port internal |
|-----------------|------|---------------|
| Nginx           | 8000 |               |
| Nginx SSL       | 3000 |               |
| Optional pg:    |      |               |
| PostgreSQL      | 5432 |   5432        |
| Adminer         | 8080 |               |
| Optional mysql: |      |               |
| MySQL           | 8989 |  3306         |
| PHPMyAdmin      | 8080 |               |

#### Connecting MySQL from [PDO](http://php.net/manual/en/book.pdo.php)

```php
<?php
    try {
        $dsn = 'mysql:host=mysql;dbname=test;charset=utf8;port=3306';
        $pdo = new PDO($dsn, 'dev', 'dev');
    } catch (PDOException $e) {
        echo $e->getMessage();
    }
```

___

### Sources

- https://github.com/nanoninja/docker-nginx-php-mysql
