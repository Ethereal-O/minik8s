package weave

import (
	"os/exec"
)

func Attach(id, ip string) error {
	_, err := exec.Command("weave", "attach", ip, id).Output()
	return err
}
