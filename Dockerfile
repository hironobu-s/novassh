FROM alpine
MAINTAINER Hironobu(hiro@hironobu.org)

RUN apk update && apk add ca-certificates openssh && rm -rf /var/cache/apk/*
RUN wget https://github.com/hironobu-s/novassh/releases/download/current/novassh-linux.amd64.gz
RUN gunzip -c novassh-linux.amd64.gz > novassh
RUN chmod +x novassh
RUN mv novassh /bin/
