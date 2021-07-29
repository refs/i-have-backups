package grpc

type counter struct {
	name   string
	latest int32
}

// Increase adds one to the counter
func (c *counter) Increase() {
	c.latest++
}
