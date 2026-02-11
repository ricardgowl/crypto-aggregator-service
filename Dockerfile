FROM ubuntu:latest
LABEL authors="ricardgo"

ENTRYPOINT ["top", "-b"]