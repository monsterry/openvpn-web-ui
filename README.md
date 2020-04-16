# OpenVPN-web-ui

## Summary
OpenVPN server web administration interface.

Goal: create quick to deploy and easy to use solution that makes work with small OpenVPN environments a breeze.

If you have docker and docker-compose installed, you can jump directly to [installation](#Prod).

![Status page](docs/images/preview_status.png?raw=true)

Please note this project is in alpha stage. It still needs some work to make it secure and feature complete.

## Motivation

Forked from [adamwalach/openvpn-web-ui](https://github.com/adamwalach/openvpn-web-ui) because I needed to bump EasyRSA to Version3 and disable server configuration.

## Features

* status page that shows server statistics and list of connected clients
* easy creation, revocation and renewal of client certificates based on [EasyRSA Version 3](https://github.com/OpenVPN/easy-rsa)
* ability to download client certificates as a zip package with client configuration inside
* log preview
* modification of OpenVPN configuration file through web interface

## Screenshots

[Screenshots](docs/screenshots.md)

## Usage

After startup web service is visible on port 8080. To login use the following default credentials:

* username: admin
* password: b3secure (this will be soon replaced with random password)

Please change password to your own immediately!

### Dev

Requirements:
* golang environments 1.14
* [beego](https://beego.me/docs/install/)

Execute commands:

    go get github.com/monsterry/openvpn-web-ui
    cd $GOPATH/src/github.com/monsterry/openvpn-web-ui
    bee run -gendoc=true


### Prod

Requirements:
* golang environments 1.14
* [beego](https://beego.me/docs/install/)

Execute commands:

    go get github.com/monsterry/openvpn-web-ui
    cd $GOPATH/src/github.com/monsterry/openvpn-web-ui

    go clean
    go build
    bee pack -ba "-tags prod" -exp=build:docs:.git:conf:data.db:go
    scp openvpn-web-ui.tar.gz mycoolserver:

On server:

    mkdir /opt/openvpn-web-ui
    tar -C /opt/openvpn-web-ui xvfz openvpn-web-ui.tar.gz

Systemd-Service:

```
root@mail:/opt/openvpn-web-ui# cat /etc/systemd/system/openvpnadmin.service
[Unit]
Description=OpenVPN Admin WebUI
After=network.target auditd.service

[Service]
WorkingDirectory=/opt/openvpn-web-ui
ExecStart=/opt/openvpn-web-ui/openvpn-web-ui
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
Alias=openvpnadmin.service
```
#### Prod Upgrade



## Todo

* add unit tests
* add option to modify certificate properties
* generate random admin password at initialization phase
* add versioning
* add automatic ssl/tls (check how [ponzu](https://github.com/ponzu-cms/ponzu) did it)


## License

This project uses [MIT license](LICENSE)

## Remarks

### Vendoring

Go Modules is used for vendoring

### Template
AdminLTE - dashboard & control panel theme. Built on top of Bootstrap 3.

Preview: https://almsaeedstudio.com/themes/AdminLTE/index2.html

