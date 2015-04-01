# pgskail
PostgreSQL high-availability made easy.

## What is

`pgskail` aims at making it easier to horizontally scale PostgreSQL by automating:

* Failover - in case of failure, perform leader election and redirect living replicas to point to new leader
* Horizontal scaling - add/remove replicas when needed
* Monitoring - constantly gather and store metrics that will tell us how our cluster is doing in real-time, while at the same time allow to determine if we need/can scale up/down

## Disclaimer

This is an idea that's been brewing in my mind for quite some time now. It has been enriched by the following projects:

* [Compose.io Governor](https://github.com/compose/governor)
* [Cybertec pgwatch](http://www.cybertec.at/en/products/pgwatch-cybertec-enterprise-postgresql-monitor)