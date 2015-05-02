package main

import (
	"io"
	"os/exec"
)

func gitClone(url, dest string) error {
	return runGitCommand(nil, nil, "", "clone", url, dest)
}

func gitPull(dir string) error {
	return runGitCommand(nil, nil, dir, "pull")
}

func runGitCommand(input io.Reader, output io.Writer, cwd string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = cwd
	c.Stdin = input
	c.Stdout = output
	c.Stderr = output

	return c.Run()
}
