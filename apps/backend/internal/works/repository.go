package works

var Works = map[string]Work{
	"1": {
		ID:         "1",
		Title:      "Test Work",
		Summary:    "This is the first work served by the new Fick backend.",
		CreatedAt:  "2026-09-07T00:00:00.000Z",
		UpdatedAt:  "2026-09-07T00:00:00.000Z",
		Visibility: WorkVisibilityPublic,
		SeriesID:   nil,
	},
	"2": {
		ID:         "2",
		Title:      "Locked Test Work",
		Summary:    "This work requires an authenticated user.",
		CreatedAt:  "2026-09-07T00:00:00.000Z",
		UpdatedAt:  "2026-09-07T00:00:00.000Z",
		Visibility: WorkVisibilityLoggedIn,
		SeriesID:   nil,
	},
}
