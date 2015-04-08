FROM phusion/baseimage
MAINTAINER pjpires@gmail.com

EXPOSE 5432

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]

VOLUME ["/var/lib/data"]

ENV PG_VERSION 9.4
RUN \
  locale-gen en_US.UTF-8 && \
  echo "deb http://archive.ubuntu.com/ubuntu trusty main universe" > /etc/apt/sources.list && \
  echo "deb http://apt.postgresql.org/pub/repos/apt trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
  apt-get update && \
  apt-get upgrade &&\
  apt-get -y install  wget ca-certificates && \
  wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
  apt-get -y --force-yes install postgresql-$PG_VERSION postgresql-client-$PG_VERSION postgresql-contrib-$PG_VERSION && \
  pg_dropcluster $PG_VERSION main && pg_createcluster --locale en_US.UTF-8 $PG_VERSION main && \
  echo "host    all             all             0.0.0.0/0 trust" >> /etc/postgresql/$PG_VERSION/main/pg_hba.conf && \
  echo "listen_addresses='*'" >> /etc/postgresql/$PG_VERSION/main/postgresql.conf

ADD pgskail /usr/bin/pgskail
ADD run_pgskail.sh /etc/service/pgskail/run

## TODO
#RUN service postgresql start && \
# su postgres sh -c "createuser -d -r -s docker" && \
# su postgres sh -c "createdb -O docker docker" && \
# su postgres sh -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE docker to docker;\""
# CMD ["su", "postgres", "-c", "/usr/lib/postgresql/$PG_VERSION/bin/postgres -D /var/lib/postgresql/$PG_VERSION/main/ -c config_file=/etc/postgresql/$PG_VERSION/main/postgresql.conf"]

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*