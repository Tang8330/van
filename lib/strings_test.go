package lib

import "testing"

func TestSpaceToDash(t *testing.T) {
	str := "hello there robin!"
	expected := "hello-there-robin!"

	if retVal := SpaceToDash(str); retVal != expected {
		t.Fatal("Failed, the values are not equal", retVal, expected)
	}
}
