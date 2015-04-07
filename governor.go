package main

import (
	"flag"
	"net"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
	//	"github.com/jackc/pgx"
)

const (
	ETCD_PREFIX      = "pgskail"
	ETCD_TTL_FOREVER = 0
	KEY_LEADER       = "leader"
)

type Options struct {
	CleanKeystore  bool
	EtcdHost       string
	HealthCheckTTL uint64
	PgHost         string
	PgPort         int
	Electable      bool
}

var (
	options  = &Options{}
	client   *etcd.Client
	pgServer string
	amLeader bool = false
)

func initializeFlags() {
	flag.Set("logtostderr", "true")
	flag.BoolVar(&options.CleanKeystore, "clean", false, "Clean-up keystore and start over?")
	flag.StringVar(&options.EtcdHost, "etcd_host", "127.0.0.1", "Hostname or IP address where Etcd is listening on")
	flag.Uint64Var(&options.HealthCheckTTL, "ttl", 10, "Health-check interval in seconds")
	flag.StringVar(&options.PgHost, "pg_host", "127.0.0.1", "Hostname or IP address where PostgreSQL server is listening on")
	flag.IntVar(&options.PgPort, "pg_port", 5432, "TCP port where PostgreSQL server is listening on")
	flag.BoolVar(&options.Electable, "electable", true, "Is leader electable?")
	flag.Parse()
}

func GetOptions() Options {
	return *options
}

/**
 * Governor function responsible for handling failover
 */
func govern() {
	// get current leader
	if leaderResponse, _ := client.Get(KEY_LEADER, false, false); leaderResponse == nil {
		// no leader. run for it?
		if options.Electable {
			glog.Info("No leader was found. Running for leader...")
			// some competitor may have won already, so fail if it is so
			if _, err := client.Create(KEY_LEADER, pgServer, options.HealthCheckTTL+5); err == nil {
				amLeader = true
				glog.Info("Won leader race")
			} else {
				glog.Info("Lost leader race")
				// look again for the leader
				govern()
			}
		} else {
			glog.Warning("No leader was found. Retrying in", options.HealthCheckTTL)
		}
	} else {
		if leader := leaderResponse.Node.Value; leader == pgServer {
			client.Update(KEY_LEADER, pgServer, options.HealthCheckTTL+5)
			// log only if we
			amPromoted := !amLeader && true
			if amLeader = true; amPromoted {
				glog.Info("I'm leader")
			}
		} else {
			amLeader = false
			glog.Info("Leader is", leader)
		}
	}
}

func main() {
	initializeFlags()
	defer glog.Flush()

	glog.Info("pgskail governor")

	// validate PostgreSQL
	if options.PgPort < 1025 || options.PgPort > 65535 {
		glog.Fatal("Bad --pg_port value: ", options.PgPort)
	}
	pgServer = options.PgHost + ":" + strconv.Itoa(options.PgPort)
	glog.Info("Connecting to PostgreSQL at ", pgServer, "...")
	if _, err := net.Dial("tcp", pgServer); err != nil {
		glog.Error("There was an error while connecting to PostgreSQL")
		glog.Fatal(err)
	}

	// validate Etcd
	etcdServer := options.EtcdHost + ":2379"
	glog.Info("Connecting to Etcd at ", etcdServer, "...")
	machines := []string{"http://" + etcdServer}
	client = etcd.NewClient(machines)
	if _, err := net.Dial("tcp", etcdServer); err != nil {
		glog.Error("There was an error while connecting to Etcd")
		glog.Fatal(err)
	}

	// need to clean-up?
	if options.CleanKeystore {
		glog.Info("Cleaning-up keystore...")
		// clean-up directory
		client.Delete(ETCD_PREFIX, true)
		// create directory
		if _, err := client.CreateDir(ETCD_PREFIX, ETCD_TTL_FOREVER); err != nil {
			glog.Fatal("There was an error while creating", ETCD_PREFIX, " directory in etcd. Code:", getErrorCode(err))
		}
		glog.Info("Keystore cleaned-up")
	}

	// govern once right now
	govern()
	// schedule future govern executions
	stop := schedule(govern, options.HealthCheckTTL)
	defer close(stop)

	select {}
}

/**
 * Return int code from etcd.EtcdError
 */
func getErrorCode(err error) int {
	etcdErr, _ := err.(*etcd.EtcdError)
	return etcdErr.ErrorCode
}

/**
 * Schedule a function to run every _interval_ seconds
 */
func schedule(what func(), interval uint64) chan struct{} {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				what()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}
