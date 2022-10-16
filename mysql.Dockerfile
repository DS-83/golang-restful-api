# syntax=docker/dockerfile:1

## Build
FROM mysql:latest

WORKDIR /usr/src/database

ENV MYSQL_ROOT_PASSWORD=123456
ENV MYSQL_DATA_DIR=/data/db
ENV MYSQL_LOG_DIR=/data/log


USER mysql:mysql

COPY ./schema.sql .
