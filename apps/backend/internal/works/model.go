package works

type WorkVisibility string

const (
	WorkVisibilityPublic   WorkVisibility = "public"
	WorkVisibilityLoggedIn WorkVisibility = "logged_in"
	WorkVisibilityPrivate  WorkVisibility = "private"
)

type Work struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Visibility   WorkVisibility `json:"visibility"`
	ChapterCount string         `json:"chapter_count"`
}
