package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
const numbers = "0123456789"

func main() {
	os.Mkdir("./corpus", 0755)

	for i := 0; i < 100; i++ {
		f, err := os.Create(fmt.Sprintf("./corpus/%d", i))
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		size := max(rand.Intn(4096), 10)
		buf := make([]byte, size)
		for j := 0; j < len(buf); j++ {
			buf[j] = alphabet[rand.Intn(len(alphabet))]
		}

		encoded := base64.StdEncoding.EncodeToString(buf)

		size = max(rand.Intn(10), 1)
		intsBuf := make([]byte, size)
		for j := 0; j < len(intsBuf); j++ {
			intsBuf[j] = numbers[rand.Intn(len(numbers))]
		}

		_, err = f.Write(buf)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = f.Write([]byte{0x0, 0x0, 0x0})
		if err != nil {
			log.Fatalln(err)
		}

		_, err = f.Write([]byte(encoded))
		if err != nil {
			log.Fatalln(err)
		}

		_, err = f.Write([]byte{0x0, 0x0, 0x0})
		if err != nil {
			log.Fatalln(err)
		}

		_, err = f.Write(intsBuf)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
