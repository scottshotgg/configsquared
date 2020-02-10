package config

type Db struct {
	port string
	addr string
}

func (db *Db) Port() string {
	return db.port
}

func (db *Db) Addr() string {
	return db.addr
}
