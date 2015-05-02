package main

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

func gitClone(url, dest string) error {
	var buf bytes.Buffer
	err := runGitCommand(nil, &buf, "", "clone", "-q", url, dest)
	if err != nil {
		return errors.New(buf.String())
	}

	return nil
}

func gitPull(dir string) error {
	var buf bytes.Buffer
	err := runGitCommand(nil, &buf, dir, "pull")
	if err != nil {
		return errors.New(buf.String())
	}

	return nil
}

func runGitCommand(input io.Reader, output io.Writer, cwd string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = cwd
	c.Stdin = input
	c.Stdout = output
	c.Stderr = output

	return c.Run()
}
