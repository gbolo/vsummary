vSummary [![](https://images.microbadger.com/badges/image/gbolo/vsummary.svg)](http://microbadger.com/images/gbolo/vsummary "Image Badge") [![Build Status](https://travis-ci.org/gbolo/vsummary.svg?branch=master)](https://travis-ci.org/gbolo/vsummary) [![Go Report Card](https://goreportcard.com/badge/github.com/gbolo/vsummary)](https://goreportcard.com/report/github.com/gbolo/vsummary)
================

vSummary is an open source tool for collecting and displaying a summary of your vSphere Environment(s). Visit the [demo site](https://vsummary.linuxctl.com/) or see the screenshots below for a better understanding.

![Alt text](https://raw.githubusercontent.com/gbolo/vSummary/php/screenshots/screenshot_1.png "Screenshot 1")

## Quick Start
To quickly give `vsummary-server` a test drive, you can spin up the included
`docker-compose` test environment.
You will need both `docker` and `docker-compose` for this :)

```
# download the docker-compose.yml (if you haven't already cloned this repo)
wget https://raw.githubusercontent.com/gbolo/vsummary/master/docker-compose.yml

# start it up
docker-compose up

# go to http://127.0.0.1:8080
```

## Requirements

### vCenter Credentials
vSummary requires a credential per vCenter server that it polls.
Ideally This credential should be a service user with **READ-ONLY** privileges to the top level vCenter object

### Database Backend
Currently, `vsummary-server` requires a MYSQL `5.x` backend. It can use MySQL, Percona, or MariaDB.

Create a new database and an associated user user for this database.
The database should be empty, as `vsummary-server` will create the required tables and schema.
The user should be given sufficient privileges to accomplish this:
```
+--------------------------------------------------------+
| Grants for vsummary@%                                  |
+--------------------------------------------------------+
| GRANT USAGE ON *.* TO 'vsummary'@'%'                   |
| GRANT ALL PRIVILEGES ON `vsummary`.* TO 'vsummary'@'%' |
+--------------------------------------------------------+
```

Once the database and credentials are prepared,
you can configure `vsummary-server` to consume it like:
```
# data source name (DSN)
# format: <username>:<password>@<host>:<port>/<database>
VSUMMARY_BACKEND_DB_DSN="vsummary:secret@(localhost:3306)/vsummary"
```

## External Polling
`vsummary-server` has a built in vCenter poller which is configurable from the web UI. However, if the `vsummary-server` does not have direct access to all vCenter servers, then you may deploy one or more `vsummary-poller` to collect from other vCenter(s) and send that information back to your `vsummary-server` (via the API) for a centralized view of your entire vCenter infrastructure. The docker image has both `vsummary-server` and `vsummary-poller` pre-installed. You can also build `vsummary-poller` for your OS of choice if you would like to avoid using docker for this.

### Usage of vsummary-poller
`vsummary-poller` has two modes:

- `deamonize`: in this mode, vCenter destinations are defined in the [configuration file](https://github.com/gbolo/vsummary/blob/master/testdata/sampleconfig/vsummary-poller.yaml) and the `vsummary-poller` will keep polling until manually stopped.

- `pollendpoint` in this mode, you must specify a vCenter via command line flags defined below:

```
Usage:
  vsummary-poller pollendpoint [flags]

Flags:
  -e, --environment string   environment/name of vcenter (friendly name)
  -h, --help                 help for pollendpoint
  -p, --password string      password for user (will prompt if not specified)
  -u, --username string      username for vcenter (readonly privilege needed)
  -s, --vcenter string       fqdn/ip of vcenter server

Global Flags:
      --config string         config file
      --log-level string      supported levels: INFO, WARNING, CRITICAL, DEBUG (default "INFO")
      --vsummary-url string   vsummary-server URL
```

## Make Targets (Development)
Contributions are welcomed and appreciated!
Ensure that your changes pass the `make lint` target.

Full list of make targets:
```
$ make help

all                        Build server and poller binaries
server                     Build server binary
poller                     Build poller binary
docker                     Build docker image
fmt                        Run gofmt on all source files
goimports                  Run goimports on all source files
lint                       Run golangci-lint for code issues
test                       Run go unit tests
integration-test           Run integration tests
setup-integration-prereqs  Setup integration prerequisites
down-integration-prereqs   Shutdown integration prerequisites
vcsim                      Start local vCenter simulator
clean                      Cleanup everything
```

binaries will be placed in `./bin` directory.

### Integration Test
In addition to passing the `make lint` target, please also ensure your changes pass
the integration test. This will bring up a local vCenter Simulator and MySQL server
(docker containers) and run tests against them:

```
# this brings up the required containers
make setup-integration-prereqs

# this runs the actual testing
make integration-test

# this stops and removes the required containers
make down-integration-prereqs
```

## Docker Deployment
A docker image is automatically created for every commit with the tag `latest`
and also for every tagged release with the image tag `x.x`.
**It is recommended to use releases**,
however the master branch should always be in a working state (hopefully).

To deploy the `vsummary-server` run the following:
```
docker run -d --name vsummary-server \
  -p 8080:8080 \
  -e VSUMMARY_BACKEND_DB_DSN="<CHANGE_ME>" \
  gbolo/vsummary:1.0
```

then open your browser to port `8080` on the deployed machine.

## Configuration
This project uses [viper](https://github.com/spf13/viper) which means all values in the configuration file
can be overwritten with environment variables. This makes it ideal for
docker deployments. For a full list of configuration options see:
[testdata/testdata/sampleconfig/vsummary-config.yaml](https://github.com/gbolo/vsummary/blob/master/testdata/sampleconfig/vsummary-config.yaml)

### Overriding configuration options with environment variables
All environment variables **must be upper case.** Each variable **must begin with** `VSUMMARY_`.
Follow the yaml indentation and for each indentation include a `_`.

For example:
```yaml
---
# http server settings -------------------------------------------------------------------------------------------------
server:
  # port to listen on
  bind_port: 8080
  # enable access log on stdout
  access_log: true
```

In the above sample yaml file, to override the port you would set: `VSUMMARY_SERVER_BIND_PORT=443`.
If you would like to use only configuration file, then mount the desired config file to `/opt/vsummary/vsummary-config.yaml`
in the docker image.

## Go Rewrite

`v1.0` of `vsummary` has been completely rewritten in `golang`. I have been working on this on and off for some time now in small bursts.
This project was meant to be a learning experience for writing a big go app with multiple components (REST API server and client, polling deamon, golang CI best practices, exc). Keep in mind that I am a DevOps guy (not a software engineer). I do this for fun.

### License
MIT

**Free Software, Hell Yeah!**
