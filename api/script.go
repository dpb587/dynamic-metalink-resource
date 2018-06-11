package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func ExecuteScript(script string, env map[string]string) ([]byte, error) {
	tmpfile, err := ioutil.TempFile("", "live-metalink")
	if err != nil {
		return nil, errors.Wrap(err, "creating script")
	}

	defer os.RemoveAll(tmpfile.Name())

	err = tmpfile.Chmod(0755)
	if err != nil {
		return nil, errors.Wrap(err, "chmoding script")
	}

	if len(script) < 2 || script[0:2] != "#!" {
		script = fmt.Sprintf("#!/bin/bash -eu\n\n%s", script)
	}

	_, err = tmpfile.WriteString(script)
	if err != nil {
		return nil, errors.Wrap(err, "writing script")
	}

	err = tmpfile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "closing script")
	}

	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command(tmpfile.Name())
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Env = mergedEnv(env)

	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "running script")
	}

	return stdout.Bytes(), nil
}

func mergedEnv(env map[string]string) []string {
	var merged []string

	for k, v := range env {
		merged = append(merged, fmt.Sprintf("%s=%s", k, v))
	}

	for _, s := range os.Environ() {
		if _, found := env[strings.SplitN(s, "=", 2)[0]]; found {
			continue
		}

		merged = append(merged, s)
	}

	return merged
}
