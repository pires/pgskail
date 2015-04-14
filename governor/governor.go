package governor

import (
	"log"
	"net"
	"strconv"

	"github.com/coreos/go-etcd/etcd"

	"github.com/pires/pgskail/service"
	"github.com/pires/pgskail/util"
)

const (
	ETCD_PREFIX      = "pgskail"
	ETCD_TTL_FOREVER = 0
	KEY_LEADER       = "leader"
	KEY_MEMBER_DIR   = "members"
)

var (
	options  service.Options
	client   *etcd.Client
	pgServer string
	amLeader bool = false
)

/**
 * Governor function responsible for handling failover
 */
func govern() {
	// get current leader
	if leaderResponse, _ := client.Get(KEY_LEADER, false, false); leaderResponse == nil {
		// no leader. run for it?
		if options.Electable {
			log.Println("No leader was found. Running for leader...")
			// some competitor may have won already, so fail if it is so
			if _, err := client.Create(KEY_LEADER, pgServer, options.HealthCheckTTL+5); err == nil {
				amLeader = true
				log.Println("Won leader race")
			} else {
				log.Println("Lost leader race")
				// look again for the leader
				govern()
			}
		} else {
			log.Println("No leader was found. Retrying in", options.HealthCheckTTL)
		}
	} else {
		if leader := leaderResponse.Node.Value; leader == pgServer {
			client.Update(KEY_LEADER, pgServer, options.HealthCheckTTL+5)
			// log only if we
			amPromoted := !amLeader && true
			if amLeader = true; amPromoted {
				log.Println("I'm leader")
			}
		} else {
			amLeader = false
			log.Println("Leader is", leader)
		}
	}
}

func Run(_options service.Options) chan struct{} {
	options = _options
	log.Println("Running governor")

	// validate PostgreSQL
	if options.PgPort < 1025 || options.PgPort > 65535 {
		log.Fatal("Bad --pg_port value: ", options.PgPort)
	}
	pgServer = options.PgHost + ":" + strconv.Itoa(options.PgPort)
	log.Println("Connecting to PostgreSQL at", pgServer, "...")
	if _, err := net.Dial("tcp", pgServer); err != nil {
		log.Fatal("There was an error while connecting to PostgreSQL -> ", err)
	}

	// validate Etcd
	// TODO retry connection to etcd a pre-defined number of times before failing
	etcdServer := options.EtcdHost + ":2379"
	log.Println("Connecting to Etcd at", etcdServer, "...")
	machines := []string{"http://" + etcdServer}
	client = etcd.NewClient(machines)
	if _, err := net.Dial("tcp", etcdServer); err != nil {
		log.Fatal("There was an error while connecting to Etcd -> ", err)
	}

	// need to clean-up?
	if options.CleanKeystore {
		log.Println("Cleaning-up keystore...")
		// clean-up directory
		client.Delete(ETCD_PREFIX, true)
		// create directory
		if _, err := client.CreateDir(ETCD_PREFIX, ETCD_TTL_FOREVER); err != nil {
			log.Fatal("There was an error while creating", ETCD_PREFIX, " directory in etcd. Code:", getErrorCode(err))
		}
		log.Println("Keystore cleaned-up")
	}

	// govern once right now
	govern()
	// schedule future govern executions
	stop := util.Schedule(govern, options.HealthCheckTTL)

	return stop
}

/**
 * Return int code from etcd.EtcdError
 */
func getErrorCode(err error) int {
	etcdErr, _ := err.(*etcd.EtcdError)
	return etcdErr.ErrorCode
}
