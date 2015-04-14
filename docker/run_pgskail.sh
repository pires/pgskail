#!/bin/sh

su postgres -c "pg_ctlcluster $PG_VERSION $PG_CLUSTER start"
exec /usr/bin/pgskail --etcd_host=$ETCD --pg_host=localhost --clean
