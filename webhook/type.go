package webhook

type Message struct {
	Username string   `json:"username,omitempty"`
	Content  string   `json:"content,omitempty"`
	Embeds   *[]Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string   `json:"title,omitempty"`
	Url         string   `json:"url,omitempty"`
	Description string   `json:"description,omitempty"`
	Color       string   `json:"color,omitempty"`
	Fields      *[]Field `json:"fields,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
