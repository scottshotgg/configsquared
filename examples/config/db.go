package config

type Db struct {
	porterino int
	addr      string
}

func (c *Db) Porterino() int {
	return c.porterino
}
func (c *Db) Addr() string {
	return c.addr
}
