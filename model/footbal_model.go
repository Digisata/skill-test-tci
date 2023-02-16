package model

type Club struct {
	ID       int    `json:"id,omitempty"`
	ClubName string `json:"clubname"`
	Points   int    `json:"points,omitempty"`
	Standing int    `json:"standing,omitempty"`
}

type RecordGameRequest struct {
	ClubHomeName string `validate:"required" json:"clubhomename" example:"Chelsea"`
	ClubAwayName string `validate:"required" json:"clubawayname" example:"Man Utd"`
	Score        string `validate:"required" json:"score" example:"1 : 2"`
}
