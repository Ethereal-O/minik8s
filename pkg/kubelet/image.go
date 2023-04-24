package kubelet

import (
	"github.com/docker/docker/api/types"
)

func PullImage(name string) error {
	events, err := Client.ImagePull(Ctx, name, types.ImagePullOptions{
		All: false,
	})
	if err != nil {
		return err
	}
	waitForPullComplete(events)
	return nil
}
