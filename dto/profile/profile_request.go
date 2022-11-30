package profiledto

import "Stage2Backend/models"

type ProfileResponse struct {
	ID       int                         `json:"id" gorm:"primary_key:auto_increment"`
	UserID   int                         `json:"user_id"`
	User     models.UsersProfileResponse `json:"user"`
	Location string                      `json:"location" gorm:"type: text"`
}
