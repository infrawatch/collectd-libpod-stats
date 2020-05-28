/*Package container represents running and paused containers in Podman

Podman stores container related information in one of two places:
1. /var/lib/containers/storage/overlay-containers (root)
2. $HOME/.local/share/containers/storage/overlay-containers (rootless)

Within these directories, the 'containers.json' file associates information regarding container
id, names and image ids.

Here this information is used to build up a lightweight representation of a container. The advantage
of this approach rather than using varlink is that 1) no running podman daemon is needed and 2)
retrieving container information is quick. In podman, the varlink GetContainer() function does an
in depth analysis on the container filesystem and when sizes are large, this can take a substantial
amount of time.
*/
package containers

import (
	"encoding/json"
	"fmt"
)

//Container represents a container within Podman
type Container struct {
	Names   []string `json:"names"`
	ID      string   `json:"id"`
	Image   string   `json:"image"`
	Created string   `json:"created"`
}

// NewListFromJSON create map[container.ID]*Container from json data
// data is expected to be in the format of podman's containers.json
func NewListFromJSON(data json.RawMessage) ([]*Container, error) {
	dataList := []json.RawMessage{}
	cList := []*Container{}
	err := json.Unmarshal(data, &dataList)
	if err != nil {
		return nil, &Error{
			msg:    "failed unmarshalling data to list",
			subErr: err,
		}
	}

	for _, cJSON := range dataList {
		c := Container{}
		err = json.Unmarshal(cJSON, &c)
		if err != nil {
			return nil, &Error{
				msg:    "failed unmarshalling data",
				subErr: err,
			}
		}

		if c.ID == "" {
			return nil, &Error{
				msg: "invalid ID",
			}
		}

		cList = append(cList, &c)
	}

	return cList, nil
}

//Error error type for all container function failures
type Error struct {
	msg    string
	subErr error
}

func (cce *Error) Error() string {
	if cce.subErr != nil {
		return fmt.Sprintf("container: %s [%s]", cce.msg, cce.subErr)
	}
	return fmt.Sprintf("container: %s", cce.msg)
}
