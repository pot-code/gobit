package api

type CursorPaginationReq struct {
	Limit  int    `json:"limit" query:"limit" validate:"min=0,max=200"`
	Cursor string `json:"cursor" query:"cursor"`
}

type CursorPaginationResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor"`
}

// NewCursorPaginationResponse create new cursor based pagination response
//
// cursor: base64 encoded cursor string
func NewCursorPaginationResponse(data interface{}, cursor string) *CursorPaginationResponse {
	return &CursorPaginationResponse{
		Data:       data,
		NextCursor: cursor,
	}
}

type OffsetPaginationReq struct {
	Page  int `query:"page" validate:"min=0"`
	Limit int `query:"limit" validate:"min=0,max=200"`
}

type OffsetPagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

type OffsetPaginationResponse struct {
	Data             interface{} `json:"data"`
	OffsetPagination `json:"pagination"`
}

// NewOffsetPaginationResponse create new offset based pagination response
//
func NewOffsetPaginationResponse(data interface{}, total, page, pages int) *OffsetPaginationResponse {
	return &OffsetPaginationResponse{
		Data: data,
		OffsetPagination: OffsetPagination{
			total, page, pages,
		},
	}
}
