package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"github.com/pires/pgskail/governor"
	"github.com/pires/pgskail/service"
)

var (
	options = &service.Options{}
)

func initializeFlags() {
	flag.Set("logtostderr", "true")
	flag.BoolVar(&options.CleanKeystore, "clean", false, "Clean-up keystore and start over?")
	flag.StringVar(&options.EtcdHost, "etcd_host", "127.0.0.1", "Hostname or IP address where Etcd is listening on")
	flag.Uint64Var(&options.LeaderTTL, "ttl", 10, "Leader health-check interval in seconds")
	flag.BoolVar(&options.MemberElectable, "electable", true, "Is member elegible for leader?")
	flag.Uint64Var(&options.MemberTTL, "member_ttl", 30, "Member health-check interval in seconds")
	flag.StringVar(&options.PgHost, "pg_host", "127.0.0.1", "Hostname or IP address where PostgreSQL server is listening on")
	flag.IntVar(&options.PgPort, "pg_port", 5432, "TCP port where PostgreSQL server is listening on")
	
	flag.Parse()
}

func main() {
	initializeFlags()

	log.Println("pgskail")

	g := governor.Run(*options)
	defer close(g)

	select {}
}
