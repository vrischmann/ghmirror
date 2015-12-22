package main

import (
	"bytes"
	"errors"
	"fmt"
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

	err := runGitCommand(nil, &buf, dir, "checkout", "master")
	if err != nil {
		return fmt.Errorf(`running command "git checkout master", err=%v`, buf.String())
	}

	buf.Reset()

	// When force pushing it will mess up the local directory sometimes, so reset everytime.
	err = runGitCommand(nil, &buf, dir, "reset", "--hard")
	if err != nil {
		return fmt.Errorf(`running command "git reset --hard", err=%v`, buf.String())
	}

	buf.Reset()

	err = runGitCommand(nil, &buf, dir, "pull", "--rebase")
	if err != nil {
		return fmt.Errorf(`running command "git pull --rebase", err=%v`, buf.String())
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
