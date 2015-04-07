package service

type Options struct {
	CleanKeystore  bool
	EtcdHost       string
	HealthCheckTTL uint64
	PgHost         string
	PgPort         int
	Electable      bool
}
