package works

func findPublicWorkByID(id string) (Work, bool) {
	work, ok := Works[id]
	if !ok {
		return Work{}, false
	}

	if work.Visibility != WorkVisibilityPublic {
		return Work{}, false
	}

	return work, true
}
