package gomaker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// BuildParams contains parameters for a build.
type BuildParams struct {
	// OS is any value supported by GOOS.
	OS string
	// Arch is any value supported by GOARCH.
	Arch string
	// Output is the name of the output program.
	Output string
	// Dir is the directory where the build is run.
	Dir string
	// Env is any additional environment variables used in the build.
	Env []string
	// LDFlags is passed as-is to the build using --ldflags
	LDFlags string
}

type env []string

func newEnv(orig []string) env {
	e := make(env, len(orig))
	copy(e, orig)

	return e
}

func (e *env) add(key, value string) {
	*e = append(*e, key+"="+value)
}

func (e *env) addFromOS(key string) {
	e.add(key, os.Getenv(key))
}

type buildParams struct {
	output  string
	dir     string
	ldflags string
	env     env
}

func makeBuildParams(bp BuildParams) (res buildParams) {
	res.output = bp.Output
	res.dir = bp.Dir
	res.ldflags = bp.LDFlags
	res.env = newEnv(bp.Env)

	res.env.add("GOOS", bp.OS)
	if bp.OS == "windows" && bp.Output != "" {
		res.output = bp.Output + ".exe"
	}

	// Host dependent stuff
	if runtime.GOOS == "windows" {
		res.env.addFromOS("TMP")
	}

	// Common stuff
	res.env.addFromOS("GOPATH")

	return
}

func makeGoBuildCommand(bp buildParams) *exec.Cmd {
	args := []string{"build"}

	if bp.output != "" {
		args = append(args, "-o", bp.output)
	}

	if bp.ldflags != "" {
		args = append(args, "--ldflags", bp.ldflags)
	}

	return exec.Command("go", args...)
}

// Build builds a Go program using the parameters provided in `params`.
func Build(params BuildParams) error {
	bp := makeBuildParams(params)

	cmd := makeGoBuildCommand(bp)
	cmd.Env = []string(bp.env)
	cmd.Dir = bp.dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("go build failed. err=%v", err)
	}

	return nil
}
