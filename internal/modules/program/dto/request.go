package dto

type PathID struct {
	ID string `validate:"required,uuid"`
}

type CreateProgramRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description"`
	ProgramType string `json:"program_type" validate:"required,oneof=podcast documentary"`
	Duration    string `json:"duration"`
	Thumbnail   string `json:"thumbnail" validate:"omitempty,url,max=2048"`
	VideoURL    string `json:"video_url" validate:"omitempty,url,max=2048"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
	CategoryID  *int64 `json:"category_id"`
	LanguageID  *int64 `json:"language_id"`
}

type UpdateProgramRequest struct {
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	ProgramType *string `json:"program_type" validate:"omitempty,oneof=podcast documentary"`
	Duration    *string `json:"duration"`
	Thumbnail   *string `json:"thumbnail" validate:"omitempty,url,max=2048"`
	VideoURL    *string `json:"video_url" validate:"omitempty,url,max=2048"`
	Status      *string `json:"status" validate:"omitempty,oneof=active inactive"`
	CategoryID  *int64  `json:"category_id"`
	LanguageID  *int64  `json:"language_id"`
}

type ListProgramsRequest struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
}

func NewListProgramsRequest(cursorStr string, limit int) ListProgramsRequest {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return ListProgramsRequest{Cursor: cursorStr, Limit: limit}
}
