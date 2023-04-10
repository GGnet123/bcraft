FROM postgres:11

ADD db/schema/init_up.sql /docker-entrypoint-initdb.d