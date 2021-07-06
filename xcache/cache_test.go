package xcache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	c := NewDefaultCache()
	c.SetDefault("user", "stonejianbu")
	v, ok := c.Get("user")
	if !ok {
		fmt.Println("=======")
	}
	fmt.Println("v ==== ", v)
}
