// +build !windows

package utils

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func DirWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, errors.New("path doesn't exist")
	}

	if !info.IsDir() {
		return false, errors.New("path isn't a directory")
	}

	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, errors.New("write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, errors.New("unable to get stat")
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, errors.New("user doesn't have permission to write to this directory")
	}
	return true, nil
}

func Ping(address string, secondsTimeout int) error {
	ping, err := exec.LookPath("ping")
	if err != nil {
		return err
	}
	out, _, err := Command(ping, address, "-c", "1", "-W", strconv.Itoa(secondsTimeout))
	if err != nil {
		return err
	}
	if strings.Contains(out, "Unknown host") {
		return errors.New("unknown host")
	}
	if strings.Contains(out, "100.0% packet loss") {
		return errors.New("destination host unreachable")
	}
	return nil
}
