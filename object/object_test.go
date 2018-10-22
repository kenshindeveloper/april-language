package object

import "testing"

func TestStringHash(t *testing.T) {
	hello0 := &String{Value: "hola mundo"}
	hello1 := &String{Value: "hola mundo"}
	diff0 := &String{Value: "la vida es bella"}
	diff1 := &String{Value: "la vida es bella"}

	if hello0.HashKey() != hello1.HashKey() {
		t.Errorf("strings with same context have different hash keys")
	}

	if diff0.HashKey() != diff1.HashKey() {
		t.Errorf("strings with same context have different hash keys")
	}

	if hello0.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different context have same hash keys")
	}
}
