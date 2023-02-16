package model

type Club struct {
	ID       int    `json:"id,omitempty"`
	ClubName string `json:"clubname"`
	Points   int    `json:"points,omitempty"`
	Standing int    `json:"standing,omitempty"`
}

type RecordGameRequest struct {
	ClubHomeName string `validate:"required" json:"clubhomename"`
	ClubAwayName string `validate:"required" json:"clubawayname"`
	Score        string `validate:"required" json:"score"`
}
