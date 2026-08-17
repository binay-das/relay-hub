package types

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqPayLoad struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Collection struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SavedRequest struct {
	ID           int64             `json:"id"`
	UserID       int64             `json:"user_id"`
	CollectionID *int64            `json:"collection_id"`
	Name         string            `json:"name"`
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

type RequestHistory struct {
	ID              int64             `json:"id"`
	UserID          int64             `json:"user_id"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Headers         map[string]string `json:"headers"`
	RequestBody     string            `json:"request_body"`
	StatusCode      int               `json:"status_code"`
	StatusText      string            `json:"status_text"`
	ElapsedMS       float64           `json:"elapsed_ms"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    string            `json:"response_body"`
	BodyType        string            `json:"body_type"`
	Error           bool              `json:"error"`
	Message         string            `json:"message"`
	CreatedAt       string            `json:"created_at"`
}

type SaveRequestPayload struct {
	CollectionID *int64            `json:"collection_id"`
	Name         string            `json:"name"`
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
}

type CollectionPayload struct {
	Name string `json:"name"`
}
