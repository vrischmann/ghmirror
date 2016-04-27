package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

func goBuild(output, ldflags string) error {
	cmd := exec.Command("go", "build", "--ldflags", ldflags, "-o", output)
	cmd.Dir = "cmd"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	musts := func(s string, err error) string {
		if err != nil {
			log.Fatal(err)
		}

		return s
	}

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	commit := musts(getCommit())
	version := musts(getVersion())
	ldflags := fmt.Sprintf("-X main.commit=%s -X main.version=%s", commit, version)

	check(goBuild("ghmirror", ldflags))
}
