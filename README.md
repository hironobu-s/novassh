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

### 3-1. SSH Connection

You can use novassh in the same way as SSH does.

```shell
novassh username@instance-name
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
```

## Credential Cache

**novassh** always sends an authentication request to Identity Service(Keystone). To reduce the connections, you may use ```--authcache``` option that save your credentials such as username, password, tenant-id, etc., in the cache file(~/.novassh). It will connect to your instance more quickly.

If you need to connect to other OpenStack environment, you may use ```--deauth``` option to remove the cache file.

## Author

Hironobu Saitoh - hiro@hironobu.org

## License

MIT
