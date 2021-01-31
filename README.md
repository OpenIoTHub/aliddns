# aliddns
[![Build Status](https://travis-ci.com/OpenIoTHub/aliddns.svg?branch=master)](https://travis-ci.com/OpenIoTHub/aliddns)

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/aliddns)

```sh
aliddns -c /path/to/config/file/aliddns.yaml
```
or just:
```
aliddns
```
(use default config file: ./aliddns.yaml)

or:
```sh
 aliddns run -i myid -k mykey -m iothub.cloud -s www -c 60
```
-i {AccessId} -k {AccessKey} -m {MainDomain} -s {SubDomainName} -c {CheckUpdateInterval}

You can install the pre-compiled binary (in several different ways),
use Docker.

Here are the steps for each of them:

## Install the pre-compiled binary

**homebrew tap** :

```sh
$ brew install OpenIoTHub/tap/aliddns
```

**homebrew** (may not be the latest version):

```sh
$ brew install aliddns
```

**snapcraft**:

```sh
$ sudo snap install aliddns
```
config file path: /root/snap/aliddns/current/aliddns.yaml

edit config file then:
```sh
sudo snap restart aliddns
```

**scoop**:

```sh
$ scoop bucket add OpenIoTHub https://github.com/OpenIoTHub/scoop-bucket.git
$ scoop install aliddns
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

config file path: /etc/aliddns/aliddns.yaml

edit config file then:
```sh
sudo systemctl restart aliddns
```

**Shell script**:

```sh
$ curl -sfL https://install.goreleaser.com/github.com/OpenIoTHub/aliddns.sh | sh
```

**manually**:

Download the pre-compiled binaries from the [releases page][releases] and
copy to the desired location.

## Running with Docker

You can also use it within a Docker container. To do that, you'll need to
execute something more-or-less like the following:

```sh
$ docker run openiothub/aliddns:latest run -i {AccessId} -k {AccessKey} -m {MainDomain} -s {SubDomainName} -c {CheckUpdateInterval}
```
example:
```sh
$ docker run openiothub/aliddns:latest run -i myid -k mykey -m iothub.cloud -s www -c 60
```
Note that the image will almost always have the last stable Go version.

[releases]: https://github.com/OpenIoTHub/aliddns/releases

