package governor

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"

	"github.com/pires/pgskail/keystore"
	"github.com/pires/pgskail/service"
	"github.com/pires/pgskail/util"
)

const ()

var (
	options  service.Options
	member   keystore.Member
	amLeader bool = false
)

// TODO schedule member ttl updates

/**
 * Governor function responsible for handling failover
 */
func govern() {
	// get current leader
	if leader, _ := keystore.GetLeader(); leader == nil {
		// no leader. run for it?
		if options.MemberElectable {
			log.Println("No leader was found. Running for leader...")
			// some competitor may have won already, so fail if it is so
			if err := keystore.SetLeader(member, options.LeaderTTL+5, false); err == nil {
				amLeader = true
				log.Println("Won leader race")
			} else {
				log.Println("Lost leader race")
				// look again for the leader
			}
		} else {
			log.Println("No leader was found. Retrying in", options.LeaderTTL)
		}
	} else {
		if leader.URL == member.URL {
			if err := keystore.SetLeader(member, options.LeaderTTL+5, true); err == nil {
				amPromoted := !amLeader && true
				if amLeader = true; amPromoted {
					log.Println("I'm leader")
				}
			} else {
				log.Fatal("Failed to update leader status", err)
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
	pgServer := options.PgHost + ":" + strconv.Itoa(options.PgPort)
	log.Println("Connecting to PostgreSQL at", pgServer, "...")
	if _, err := net.Dial("tcp", pgServer); err != nil {
		log.Fatal("There was an error while connecting to PostgreSQL -> ", err)
	}

	// validate Etcd
	// TODO retry connection to etcd a pre-defined number of times before failing
	etcdServer := options.EtcdHost + ":2379"
	log.Println("Connecting to Etcd at", etcdServer, "...")
	if err := keystore.Connect(etcdServer); err != nil {
		log.Fatal("There was an error while connecting to Etcd -> ", err)
	}

	// need to clean-up?
	if options.CleanKeystore {
		log.Println("Cleaning-up keystore...")
		if err := keystore.Initialize(); err != nil {
			log.Fatal("There was an error while initializing keystore, code:", getErrorCode(err))
		}
		log.Println("Keystore cleaned-up")
	}

	// register as member
	url := "postgresql://" + pgServer // TODO add authentication support
	member = keystore.Member{Name: options.PgHost, URL: url, TTL: options.MemberTTL, LastCheck: uint64(time.Now().UnixNano() / 1000000)}
	if registered := member.Register(); !registered {
		log.Fatal("Couldn't register as a member")
	}
	// govern once right now
	govern()
	// schedule future govern executions
	stop := util.Schedule(govern, options.LeaderTTL)

	return stop
}

/**
 * Return int code from etcd.EtcdError
 */
func getErrorCode(err error) int {
	etcdErr, _ := err.(*etcd.EtcdError)
	return etcdErr.ErrorCode
}
