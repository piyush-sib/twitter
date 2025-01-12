package closer

import (
	"log"

	"go.uber.org/dig"
)

type CloserResult struct {
	dig.Out

	Close func() error `group:"closers"`
}

type CloserParams struct {
	dig.In

	CloseFunctions []func() error `group:"closers"`
}

type Closer struct {
	closeFunctions []func() error
}

func NewCloser(c CloserParams) *Closer {
	return &Closer{
		c.CloseFunctions,
	}
}

func (c *Closer) Close() {
	log.Println("closing everything")
	for _, f := range c.closeFunctions {
		if err := f(); err != nil {
			log.Printf("error while closing: %e \n", err)
		}
	}
}
