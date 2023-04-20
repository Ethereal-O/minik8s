package kubelet

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/object"
	"strings"
)

// ------------------Container------------------

// toVolumeBinds returns the binds of volumes
func toVolumeBinds(pod *object.Pod, target *object.Containers) []string {
	// Get volume devices, create a map
	// Mapping volume name to its source
	volumes := make(map[string]string)
	// Now we only support HostPath
	for _, volume := range pod.Spec.Volumes {
		if volume.Type == "hostPath" {
			volumes[volume.Name] = volume.Path
		}
	}

	var volumeBinds []string
	for _, volumeMount := range target.VolumeMounts {
		volumeName := volumeMount.Name
		// If the specified volume device is existent, and is hostPath(we only support this type temporarily)
		if device, exists := volumes[volumeName]; exists && device == "hostPath" {
			// Volume bind rule: $(host path):$(container path)
			mountRule := device + ":" + volumeMount.MountPath
			volumeBinds = append(volumeBinds, mountRule)
		}
	}
	return volumeBinds
}

// ------------------Image------------------

func parseAndPrintPullEvents(events io.ReadCloser, imageName string) {
	d := json.NewDecoder(events)

	type Event struct {
		Status         string `json:"status"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
		Error string `json:"error"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Printf("EVENT: %+v\n", event)
	}

	if event != nil {
		if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", imageName)) {
			// new
			fmt.Printf("Image %s is new.\n", imageName)
		}
		if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", imageName)) {
			// up-to-date
			fmt.Printf("Image %s is up-to-date.\n", imageName)
		}
	}
}
