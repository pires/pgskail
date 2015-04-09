FROM phusion/baseimage
MAINTAINER pjpires@gmail.com

EXPOSE 5432

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]

VOLUME ["/var/lib/data"]

ENV ETCD localhost
ENV PG_VERSION 9.4
ENV PG_CLUSTER pgskail

RUN \
  locale-gen en_US.UTF-8 && \
  echo "deb http://archive.ubuntu.com/ubuntu trusty main universe" > /etc/apt/sources.list && \
  echo "deb http://apt.postgresql.org/pub/repos/apt trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
  apt-get update && \
  apt-get upgrade &&\
  apt-get -y install  wget ca-certificates && \
  wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
  apt-get update && \
  apt-get -y install postgresql-$PG_VERSION postgresql-client-$PG_VERSION postgresql-contrib-$PG_VERSION && \
  pg_dropcluster $PG_VERSION main && pg_createcluster --locale en_US.UTF-8 $PG_VERSION $PG_CLUSTER && \
  echo "host    all             all             0.0.0.0/0 trust" >> /etc/postgresql/$PG_VERSION/$PG_CLUSTER/pg_hba.conf && \
  echo "listen_addresses='*'" >> /etc/postgresql/$PG_VERSION/$PG_CLUSTER/postgresql.conf

ADD pgskail /usr/bin/pgskail
ADD run_pgskail.sh /etc/service/pgskail/run

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*