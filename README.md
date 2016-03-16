# novassh [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Build Status](https://travis-ci.org/hironobu-s/novassh.svg?branch=master)](https://travis-ci.org/hironobu-s/novassh) [![codebeat badge](https://codebeat.co/badges/97e0e868-2796-41d9-82a1-d1740acdc4d3)](https://codebeat.co/projects/github-com-hironobu-s-novassh)


This is an SSH wrapper program to connect OpenStack(Nova) instance with the instance name.

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

### 1. Authentication

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

Use ``--novassh-list`` option.

```
novassh --novassh-list
```

### 3. Connect to the instance.

You can connect to it with the instance name instead of Hostname or IP Address.

```shell
novassh username@instance-name
```

And you can also use novassh with some options for SSH command.

```shell
novassh -L 8080:internal-host:8080 username@instance-name
```

## Options

```
OPTIONS:
	--novassh-command: Specify SSH command (default: "ssh").
	--novassh-deauth:  Remove credential cache.
	--novassh-debug:   Output some debug messages.
	--novassh-list:    Display instances.
	--novassh-help:            Print this message.

    Any other options from novassh will pass to the SSH command.

ENVIRONMENTS:
	NOVASSH_COMMAND: Specify SSH command (default: "ssh").
```

## Credential cache

novassh saves your authentication information such as username, password, tenant-id to the cache file(~/.novassh) in order to reduce the connection for the Identity service(Keystone). You can use ```--novassh-deauth``` option to remove it. 

## Author

Hironobu Saitoh - hiro@hironobu.org

## License

MIT
