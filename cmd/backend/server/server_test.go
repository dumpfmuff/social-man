package server

import "testing"

func TestInit(t *testing.T) {
	Init("https", "localhost", 443)
}
