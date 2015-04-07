# pgskail
PostgreSQL high-availability made easy.

## What is

`pgskail` aims at making it easier to horizontally scale PostgreSQL by automating:

* [Failover](https://github.com/pires/pgskail/wiki/Failover) - in case of failure, perform leader election and redirect living replicas to point to new leader
* Horizontal Scaling - add/remove replicas when needed
* [Monitoring](https://github.com/pires/pgskail/wiki/Monitoring) - constantly gather and store metrics that will tell us how our cluster is doing in real-time, while at the same time allow to determine if we need/can scale up/down

This is an idea that's been brewing in my mind for quite some time now. It has been enriched by the following projects:

* [Compose.io Governor](https://github.com/compose/governor)
* [Cybertec pgwatch](http://www.cybertec.at/en/products/pgwatch-cybertec-enterprise-postgresql-monitor)

## Build

```
make
```

## Run

`pgskail` depends on a `etcd` cluster. For development purposes, you can use `docker` like this

```
docker run -d --net host -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 quay.io/coreos/etcd:v2.0.7
```

`pgskail` also depends on a `postgresql` instance. For development purposes, you can use `docker` like this

```
docker run -d --net host --name postgres -e POSTGRES_PASSWORD=teste postgres
``

Now, you can run `pgskail` like  

```
./build/pgskail --etcd_host 192.168.59.104 --pg_host 192.168.59.104
``

## Flags

* `clean` - Clean-up keystore and start over? (default is **false**)
* `etcd_host` - Hostname or IP address where Etcd is listening on (default is **127.0.0.1**)
* `ttl` - Health-check interval in seconds (default is **10**)
* `pg_host` - Hostname or IP address where PostgreSQL server is listening on (default is **127.0.0.1**)
* `pg_port` - TCP port where PostgreSQL server is listening on (default is **5432**)
* `electable` - Is leader electable? (default is **true**)
