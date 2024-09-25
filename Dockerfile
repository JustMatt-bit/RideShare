FROM mysql:9.0.1

ENV MYSQL_ROOT_PASSWORD=root

ENV MYSQL_DATABASE=rideshare
ENV MYSQL_USER=rideshare
ENV MYSQL_PASSWORD=rideshare

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306