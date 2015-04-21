package service

type Options struct {
	CleanKeystore   bool
	EtcdHost        string
	LeaderTTL       uint64
	MemberElectable bool
	MemberTTL       uint64
	PgHost          string
	PgPort          int
}
