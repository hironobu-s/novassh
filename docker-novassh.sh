#!/bin/sh

ARGS=$@
DOTSSH=~/.ssh
NAME=novassh

docker run \
       -ti \
       --rm \
       --name $NAME \
       -e OS_TENANT_NAME=$OS_TENANT_NAME \
       -e OS_PASSWORD=$OS_PASSWORD \
       -e OS_AUTH_URL=$OS_AUTH_URL \
       -e OS_USERNAME=$OS_USERNAME \
       -e OS_REGION_NAME=$OS_REGION_NAME \
       -v $DOTSSH:/root/.ssh \
       hironobu/novassh \
       /bin/novassh $@
