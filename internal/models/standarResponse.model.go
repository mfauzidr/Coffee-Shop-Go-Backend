package models

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Meta    *Meta       `json:"meta,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Meta struct {
	Total     int `json:"totalData,omitempty"`
	TotalPage int `json:"totalPage,omitempty"`
	Page      int `json:"page,omitempty"`
	NextPage  int `json:"nextPage,omitempty"`
	PrevPage  int `json:"prevPage,omitempty"`
}
