package libfoo

import (
	"rsc.io/quote"
	"fmt"
)

// clientImpl - An implementation of Client interface
//
type clientImpl struct {

}

// NewClient - return an `Client` interface
//
func NewClient() Client {
	return &clientImpl{}
}

// Method1 - the implementation of `Method1` method defined in the `Client` interface
//
func (c *clientImpl) Method1() error {

	// Concrete implementation details are omitted

	return nil
}

// Method2 - the implementation of `Method2` method defined in the `Client` interface
//
func (c *clientImpl) Method2() error {

	// Concrete implementation details are omitted

	return nil
}

// Method3 - the implementation of `Method3` method defined in the `Client` interface
//
func (c *clientImpl) Method3() error {

	// Concrete implementation details are omitted

	return nil
}

// Method4 - the implementation of `Method4` method defined in the `Client` interface
//
func (c *clientImpl) Method4() error {

	// Concrete implementation details are omitted

	fmt.Println("v2 " + quote.Hello())

	return nil
}

// Method5 - the implementation of `Method5` method defined in the `Client` interface
//
func (c *clientImpl) Method5(param1 int, param2 string) error {

	// Concrete implementation details are omitted

	fmt.Println("v2: Method5()")

	return nil
}


// Method6 - the implementation of `Method6` method defined in the `Client` interface
//
func (c *clientImpl) Method6() error {

	// Concrete implementation details are omitted

	return nil
}
