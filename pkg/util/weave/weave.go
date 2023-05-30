package weave

import (
	"os/exec"
)

func Attach(id, ip string) error {
	_, err := exec.Command("weave", "attach", ip, id).Output()
	return err
}

func Expose(ip string) error {
	_, err := exec.Command("weave", "expose", ip).Output()
	return err
}

func Reset() error {
	_, err := exec.Command("weave", "reset").Output()
	return err
}
