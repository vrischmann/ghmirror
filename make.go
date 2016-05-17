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
	"strings"
	"time"

	"git.vrischmann.me/bstats/pkg"
)

var (
	flLinux = flag.Bool("linux", false, "Build with GOOS=linux")
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

func getEnv() []string {
	env := os.Environ()

	// TODO(vincent): make this generic somehow
	if !(*flLinux) {
		return env
	}

	goosPos := 0
	for i, el := range env {
		if strings.HasPrefix(el, "GOOS=") {
			goosPos = i
		}
	}

	env = append(env[:goosPos], env[goosPos+1:]...)
	env = append(env, "GOOS=linux")

	return env
}

func goBuild(output, ldflags string) error {
	cmd := exec.Command("go", "build", "--ldflags", ldflags, "-o", output)
	cmd.Dir = "cmd"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = getEnv()

	return cmd.Run()
}

func check(err error) {
	if err != nil {
		log.Printf("check failed: %v", err)
		failed = true
	}
}

var failed bool

const bstatsFile = "ghmirror.bst"

func main() {
	now := time.Now()
	check(bstats.Begin(bstatsFile))

	flag.Parse()

	musts := func(s string, err error) string {
		if err != nil {
			log.Fatal(err)
		}

		return s
	}

	commit := musts(getCommit())
	version := musts(getVersion())
	ldflags := fmt.Sprintf("-X main.commit=%s -X main.version=%s", commit, version)

	check(goBuild("ghmirror", ldflags))

	var statusCode int
	if failed {
		statusCode = 1
	}

	check(bstats.End(bstatsFile, statusCode))

	elapsed := time.Since(now)
	log.Printf("build time: %s", elapsed)

	os.Exit(statusCode)
}
