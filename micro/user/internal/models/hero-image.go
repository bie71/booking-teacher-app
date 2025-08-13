package models

type HeroImage struct {
	Id       int    `gorm:"primary_key auto_increment" json:"id"`
	KeyImage string `gorm:"unique" json:"key_image"`
	ImageUrl string `json:"image_url"`
}
