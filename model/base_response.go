package model

type SuccessResult struct {
	Code   int    `json:"code" example:"200"`
	Status string `json:"status" example:"Ok"`
	Data   any    `json:"data"`
}

type BadRequestResult struct {
	Code   int    `json:"code" example:"400"`
	Status string `json:"status" example:"BAD REQUEST"`
}

type NotFoundResult struct {
	Code   int    `json:"code" example:"404"`
	Status string `json:"status" example:"NOT FOUND"`
}

type InternalServerErrorResult struct {
	Code   int    `json:"code" example:"500"`
	Status string `json:"status" example:"INTERNAL SERVER ERROR"`
}
