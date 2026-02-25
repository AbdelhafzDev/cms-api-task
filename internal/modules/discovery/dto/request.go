package dto

type SearchRequest struct {
	Query       string `json:"q" validate:"required,min=1,max=255"`
	ProgramType string `json:"program_type" validate:"omitempty,oneof=podcast documentary"`
	Category    string `json:"category" validate:"omitempty,max=255"`
	Language    string `json:"language" validate:"omitempty,max=50"`
	Page        int    `json:"page" validate:"omitempty,min=1"`
	PerPage     int    `json:"per_page" validate:"omitempty,min=1,max=100"`
}

func NewSearchRequest(q, programType, category, language string, page, perPage int) SearchRequest {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return SearchRequest{
		Query:       q,
		ProgramType: programType,
		Category:    category,
		Language:    language,
		Page:        page,
		PerPage:     perPage,
	}
}

type ListRequest struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
}

func NewListRequest(cursorStr string, limit int) ListRequest {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return ListRequest{Cursor: cursorStr, Limit: limit}
}

type PathID struct {
	ID string `validate:"required,uuid"`
}
