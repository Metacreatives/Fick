package works

import "context"

func findPublicWorkByID(
	ctx context.Context,
	repository *Repository,
	id string,
) (Work, bool, error) {
	work, found, err := repository.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return Work{}, false, err
	}

	if !found {
		return Work{}, false, nil
	}

	if work.Visibility != WorkVisibilityPublic {
		return Work{}, false, nil
	}

	return work, true, nil
}
