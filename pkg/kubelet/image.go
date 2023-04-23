package kubelet

import (
	"github.com/docker/docker/api/types"
)

func PullImage(name string, config *PullConfig) error {
	_, err := Client.ImagePull(Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return err
	}
	return nil
}
