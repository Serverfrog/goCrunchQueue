package internal

import (
	"testing"
)

func Test_main(t *testing.T) {
	path := "config/config.yaml"
	ReadConfig(&path)
	StartServer("Test")
}
