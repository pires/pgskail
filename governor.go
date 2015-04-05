package main

import (
	"flag"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type Options struct {
	EtcdHost       string
	HealthCheckTTL time.Duration
	RunForLeader   bool
}

var options = &Options{}

func initFlags() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&options.EtcdHost, "etcd_server", "http://127.0.0.1:2379", "The etcd server IP:port pair to connect to")
	flag.DurationVar(&options.HealthCheckTTL, "health_check_ttl", time.Duration(10)*time.Second, "Health-check time-to-live, e.g. 10s. Must confirm to http://golang.org/pkg/time/#ParseDuration")
	flag.BoolVar(&options.RunForLeader, "run_for_leader", false, "Run for leader?")
	flag.Parse()
}

const ETCD_PREFIX = "pgskail/governor"

func main() {
	initFlags()

	defer glog.Flush()

	glog.Infoln("pgskail governor")
	glog.Info("Connecting to ", options.EtcdHost)
	machines := []string{options.EtcdHost}
	client := etcd.NewClient(machines)
	if _, err := client.Set(ETCD_PREFIX, "test", 0); err != nil {
		glog.Fatalln("There was an error while connecting to etcd at ", options.EtcdHost, "\n  ", err)
	}
}
