package fuzz

import (
	"bytes"
	"os"

	"github.com/vrischmann/envconfig"
)

func Fuzz(data []byte) int {
	var conf struct {
		String string
		Bytes  []byte
		Int    int
	}

	parts := bytes.Split(data, []byte{0x0, 0x0, 0x0})

	os.Setenv("STRING", string(parts[0]))
	os.Setenv("BYTES", string(parts[1]))
	os.Setenv("INT", string(parts[2]))

	if err := envconfig.Init(&conf); err != nil {
		panic(err)
	}

	return 1
}
