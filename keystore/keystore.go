package keystore

import (
	"encoding/json"
	"net"

	"github.com/coreos/go-etcd/etcd"

	"github.com/pires/pgskail/util"
)

const (
	PREFIX      = "pgskail"
	TTL_FOREVER = 0
	KEY_LEADER  = "leader"
	DIR_MEMBERS = "members"
)

var (
	client *etcd.Client
)

// begin member stuff
type Member struct {
	Name, URL string
	TTL       uint64
	LastCheck uint64
}

func (member Member) Register() bool {
	json, err := json.Marshal(member)
	if err == nil {
		path := []string{PREFIX, DIR_MEMBERS}
		return writeKey(path, member.Name, string(json), member.TTL, true) == nil
	}

	return false
}

func (member Member) UnRegister() bool {
	path := []string{PREFIX, DIR_MEMBERS}
	return remove(path, member.Name) == nil
}

// end member stuff

// begin keystore stuff
func Initialize() error {
	var err error = nil
	if err = remove([]string{}, PREFIX); err != nil {
		if err = writeDir([]string{}, PREFIX, TTL_FOREVER, false); err != nil {
			err = writeDir([]string{PREFIX}, DIR_MEMBERS, TTL_FOREVER, false)
		}
	}
	return err
}

func GetLeader() (*Member, error) {
	response, err := get([]string{PREFIX}, KEY_LEADER)
	if err == nil && response != "" {
		var leader Member
		b := []byte(response)
		err = json.Unmarshal(b, &leader)
		return &leader, err
	}
	return nil, err
}

func SetLeader(member Member, ttl uint64, overwrite bool) error {
	json, err := json.Marshal(member)
	if err == nil {
		err = writeKey([]string{PREFIX}, KEY_LEADER, string(json), ttl, overwrite)
	}
	return err
}

// end keystore stuff

// etcd stuff
/*
 * TODO add support for multiple etcd servers
 */
func Connect(etcdServer string) error {
	machines := []string{"http://" + etcdServer}
	client = etcd.NewClient(machines)
	_, err := net.Dial("tcp", etcdServer)
	return err
}

func get(pathNodes []string, key string) (string, error) {
	path := util.MakePath(append(pathNodes, key))
	response, err := client.Get(path, true, true)
	var value string
	if response != nil {
		value = response.Node.Value
	}
	return value, err
}

func remove(pathNodes []string, key string) error {
	path := util.MakePath(append(pathNodes, key))
	_, err := client.Delete(path, true)
	return err
}

/**
 * Creates or Updates a directory depending on whether _override_ is false
 * or true, on a path, with a certain time-to-live.
 */
func writeDir(pathNodes []string, dir string, ttl uint64, override bool) error {
	var err error = nil
	path := util.MakePath(append(pathNodes, dir))
	if override {
		_, err = client.UpdateDir(path, ttl)
	} else {
		_, err = client.CreateDir(path, ttl)
	}
	return err
}

/**
 * Creates or Updates a key/value pair depending on whether _override_ is false
 * or true, on a path, with a certain time-to-live.
 */
func writeKey(pathNodes []string, key string, value string, ttl uint64, override bool) error {
	var err error = nil
	path := util.MakePath(append(pathNodes, key))
	if override {
		_, err = client.Set(path, value, ttl)
	} else {
		_, err = client.Create(path, value, ttl)
	}
	return err
}

// end etcd stuff
