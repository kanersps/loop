package object

import (
	"github.com/kanersps/loop/models"
	"testing"
)

func TestString_HashKey(t *testing.T) {
	test1 := &models.String{Value: "test"}
	test2 := &models.String{Value: "test"}

	test3 := &models.String{Value: "test2"}
	test4 := &models.String{Value: "test2"}

	if test1.HashKey() != test2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if test3.HashKey() != test4.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if test1.HashKey() == test4.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
