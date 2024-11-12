package errors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Errors(t *testing.T) {
	err := test3()
	assert.Equal(t, 7, len(strings.Split(err.Error(), "\n")))
}

func test1() error {
	return Wrap(fmt.Errorf("error"))
}

func test2() error {
	if err := test1(); err != nil {
		return Wrap(err)
	}
	return nil
}

func test3() error {
	if err := test2(); err != nil {
		return Wrap(err)
	}
	return nil
}
