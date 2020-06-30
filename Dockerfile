FROM --platform=amd64 debian:buster-backports
RUN apt-get update && apt-get install -y lxc -t buster-backports

ENTRYPOINT ["sleep","2000"]