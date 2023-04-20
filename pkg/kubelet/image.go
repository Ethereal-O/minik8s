package kubelet

import (
	"github.com/docker/docker/api/types"
)

func PullImage(name string, config *PullConfig) error {
	events, err := Client.ImagePull(Ctx, name, types.ImagePullOptions{
		All: config.All,
	})
	if err != nil {
		return err
	}
	if config.Verbose {
		parseAndPrintPullEvents(events, name)
	}
	return nil
}
