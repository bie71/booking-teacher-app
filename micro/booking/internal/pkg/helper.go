package pkg

type Pagination struct {
	Status       string `json:"status"`
	StartDateStr string `json:"start_date"`
	EndDateStr   string `json:"end_date"`
	PageStr      string `json:"page"`
	LimitStr     string `json:"limit"`
	TeacherID    string `json:"teacher_id"`
}

type Paginate struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationPage struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	Limit       int `json:"limit"`
}

type ResponsePaginate struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationPage `json:"pagination"`
}
