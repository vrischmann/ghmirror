// +build makefile

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/vrischmann/flagutil"
	"github.com/vrischmann/gomaker"
)

func getVersion() (string, error) {
	data, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		return "", fmt.Errorf("unable to read VERSION. err=%v", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func getCommit() (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("unable to run git rev-parse. out=%s err=%v", strings.TrimSpace(buf.String()), err)
	}

	return strings.TrimSpace(buf.String()), nil
}

func musts(s string, err error) string {
	if err != nil {
		log.Fatal(err)
	}

	return s
}

var (
	failed bool

	flOS   flagutil.Strings
	flArch flagutil.Strings
)

func init() {
	flag.Var(&flOS, "os", "List of GOOS values")
	flag.Var(&flArch, "arch", "List of GOARCH values")
}

func main() {
	flag.Parse()

	now := time.Now()

	if len(flOS) == 0 {
		flOS = flagutil.Strings{runtime.GOOS}
	}
	if len(flArch) == 0 {
		flArch = flagutil.Strings{runtime.GOARCH}
	}

	commit := musts(getCommit())
	version := musts(getVersion())
	ldflags := fmt.Sprintf("-X main.commit=%s -X main.version=%s", commit, version)

	for _, os := range flOS {
		for _, arch := range flArch {
			now2 := time.Now()

			err := gomaker.Build(gomaker.BuildParams{
				OS:      os,
				Arch:    arch,
				Output:  "ghmirror_" + os + "_" + arch,
				Dir:     "cmd",
				LDFlags: ldflags,
			})
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("built for %s_%s in %s", os, arch, time.Since(now2))
		}
	}

	elapsed := time.Since(now)
	log.Printf("build time: %s", elapsed)
}
