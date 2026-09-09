package chapters

type Chapter struct {
	ID            string `json:"id"`
	WorkID        string `json:"work_id"`
	Number        string `json:"number"`
	Title         string `json:"title"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ContentFormat string `json:"content_format"`
	ContentRaw    string `json:"content_raw"`
}
