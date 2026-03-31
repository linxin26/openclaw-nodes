package core

type ImagePayload struct {
	Base64    string `json:"base64"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height,omitempty"`
	Size      int    `json:"size"`
	Timestamp int64  `json:"timestamp"`
}

type VideoPayload struct {
	Base64       string `json:"base64,omitempty"`
	Format       string `json:"format"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	Size         int    `json:"size,omitempty"`
	DurationMs   int    `json:"durationMs"`
	Timestamp    int64  `json:"timestamp"`
	IncludeAudio bool   `json:"includeAudio,omitempty"`
}
