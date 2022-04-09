package main

import (
	"os"
	"testing"
)

var oldArgs = os.Args

func TestFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd"}
}

func TestConfigFileParsing(t *testing.T) {

}

func TestFlagsOverwriteConfigFile(t *testing.T) {

}
