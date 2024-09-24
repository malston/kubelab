package os

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type ShellExecutor struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func (s *ShellExecutor) Execute(cmd string, args ...string) (string, string, error) {
	c := exec.Command(cmd, args...)
	c.Stdout = &s.Stdout
	c.Stderr = &s.Stderr

	err := c.Run()
	if err != nil {
		return "", "", errors.WithMessage(err, fmt.Sprintf("error executing %s", cmd))
	}

	if len(s.Stderr.String()) > 0 {
		return "", "", errors.WithMessagef(err, "error executing command: %s", s.Stderr.String())
	}

	return s.Stdout.String(), "", nil
}
