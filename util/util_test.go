package util

import "testing"

func TestRandomID(t *testing.T) {
	id, err := GenerateRandomID()
	if err != nil {
		t.Fatal(err)
	}

	if len(id) == 0 {
		t.Fatal(err)
	}
}
