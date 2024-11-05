package onitama

import (
	"testing"
)

func Test_Onitama(t *testing.T) {
	onitama, err := NewOnitama(123, []string{"red", "blue"})
	if err != nil {
		t.Fatal(err)
	}
	action := onitama.actions()[0]
	if err := onitama.Do(action); err != nil {
		t.Fatal(err)
	}
}
