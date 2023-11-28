package soundx

import (
	"fmt"
	"strconv"
)

func parseId(resource, id string) (int64, error) {
	if n, err := strconv.ParseInt(id, 10, 64); err != nil {
		return 0, fmt.Errorf("invalid %s id: %w", resource, err)
	} else {
		return n, nil
	}
}
