package works

import (
	"encoding/json"
	"fick/backend/internal/store"
)

func findPublicWorkByID(id string) (Work, bool) {
	var works []Work

	if err := json.Unmarshal([]byte(store.Works), &works); err != nil {
		return Work{}, false
	}

	for _, work := range works {
		if work.ID != id {
			continue
		}

		if work.Visibility != WorkVisibilityPublic {
			return Work{}, false
		}

		return work, true
	}

	return Work{}, false
}
