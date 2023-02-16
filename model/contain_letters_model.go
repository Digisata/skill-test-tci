package model

type ContainLettersRequest struct {
	FirstWord string `validate:"required" json:"first_word" example:"cat"`
	SecondWord string `validate:"required" json:"second_word" example:"artarctica"`
}