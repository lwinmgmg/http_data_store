package views

type FileRead struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
}
