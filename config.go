package ms

type IConfig interface {
	Get() interface{}
}

type config struct {
	value interface{}
}

func (c *config) Get() interface{} {
	return c.value
}
