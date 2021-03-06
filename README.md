# novassh [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Build Status](https://travis-ci.org/hironobu-s/novassh.svg?branch=master)](https://travis-ci.org/hironobu-s/novassh) [![codebeat badge](https://codebeat.co/badges/97e0e868-2796-41d9-82a1-d1740acdc4d3)](https://codebeat.co/projects/github-com-hironobu-s-novassh)

# Overview

**novassh** is a client program for OpenStack(Nova). You can connect to your instance with the instance name instead of Hostname or IP Address via SSH, and also support for a serial console access.

It has been tested on the following environments.

* Rackspace https://www.rackspace.com/
* ConoHa https://www.conoha.jp/
* My OpenStack environment(Liberty)


# Install

Download an executable binary from GitHub release.

**Mac OSX**

```shell
curl -sL https://github.com/hironobu-s/novassh/releases/download/current/novassh-osx.amd64.gz | zcat > novassh && chmod +x ./novassh
```

**Linux(amd64)**

```shell
curl -sL https://github.com/hironobu-s/novassh/releases/download/current/novassh-linux.amd64.gz | zcat > novassh && chmod +x ./novassh
```

**Windows(amd64)**

[ZIP file](https://github.com/hironobu-s/novassh/releases/download/current/novassh.amd64.zip)

# Run in docker

You can run **novassh** inside a container.

```
docker run -ti --rm hironobu/novassh novassh
```

See https://hub.docker.com/r/hironobu/novassh/

# How to use.

### 1. Authentication.

Set the authentication information to environment variables.

```shell
export OS_USERNAME=[username]
export OS_PASSWORD=[password]
export OS_TENANT_NAME=[tenant name]
export OS_AUTH_URL=[identity endpoint]
export OS_REGION_NAME=[region]
```

See also: https://wiki.openstack.org/wiki/OpenStackClient/Authentication

### 2. Show instance list.

Use ``--list`` option.

```
novassh --list
```

Output

```
[Name]       [IP Address]
hironobu-dev 133.130.***.***
go-build     133.130.***.***
test-app1    133.130.***.***
```

### 3-1. SSH Connection

You can use novassh in the same way as SSH does.

```shell
novassh username@instance-name
```

For example.

```shell
novassh hiro@hironobu-dev
```

All options are passed to SSH command.

```shell
novassh -L 8080:internal-host:8080 username@instance-name
```

### 3-2. Serial Console Connection

You can use ```--console``` option to access your instance via serial console. (OpenStack has supported for serial console access to your instance since version Juno.)

```shell
novassh --console username@instance-name
```

Type ```"Ctrl+[ q"``` to disconnect.

### 3-3. Debug output

You can use ```--debug``` option to figure out the problems.

```
DEBU[0000] Command: LIST
DEBU[0000] Send    ==>: POST https://identity.tyo1.conoha.io/v2.0/tokens
DEBU[0000] map[Content-Type:[application/json] Accept:[application/json]]
DEBU[0000] Receive <==: 200 https://identity.tyo1.conoha.io/v2.0/tokens (size=2541)
DEBU[0000] Send    ==>: GET https://compute.tyo1.conoha.io/v2/####################################/servers/detail
DEBU[0000] map[X-Auth-Token:[XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX] Accept:[application/json]]
DEBU[0001] Receive <==: 200 https://compute.tyo1.conoha.io/v2/####################################/servers/detail (size=2302)
DEBU[0001] Machine found: name=example-vm-name, ipaddr=150.95.0.0
DEBU[0001] InterfaceName: ext-150-95-0-0-1
DEBU[0001] InterfaceName: local-gnct47070904-1
[Name]                  [IP Address]
example-vm-name         150.95.0.0
```

## Options

```
OPTIONS:
	--authcache: Store credentials to the cache file ($HOME/.novassh).
	--command:   Specify SSH command (default: "ssh").
	--console:   Use an serial console connection instead of SSH.
	--deauth:    Remove credential cache.
	--debug:     Output some debug messages.
	--list:      Display instances.
	--help:      Print this message.

    Any other options will pass to SSH command.

ENVIRONMENTS:
	NOVASSH_COMMAND: Specify SSH command (default: "ssh").
	NOVASSH_INTERFACE: Specify network interface of instance (default: blank strings which means the auto detection).
```

## Credential Cache

**novassh** always sends an authentication request to Identity Service(Keystone). To reduce the connections, you may use ```--authcache``` option that save your credentials such as username, password, tenant-id, etc., in the cache file(~/.novassh). It will connect to your instance more quickly.

If you need to connect to other OpenStack environment, you may use ```--deauth``` option to remove the cache file.

## Author

Hironobu Saitoh - hiro@hironobu.org

## License

MIT
