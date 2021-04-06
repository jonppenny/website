# Notes for server setup

[DO server setup tutorial](https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04)

[DO server setup tutorial alternative](https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker-and-nginx-on-ubuntu-18-04)

## Example service file
```shell
[Unit]
Description=goweb

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/home/jon/go-web
ExecStart=/home/jon/go-web/web

[Install]
WantedBy=multi-user.target
```

Note the working directory line, this is required for static files.