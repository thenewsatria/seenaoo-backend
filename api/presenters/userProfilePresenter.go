package presenters

type UserProfile struct {
	DisplayName     string `bson:"display_name" json:"displayName" validate:"required,min=1,max=64"`
	AvatarImagePath string `bson:"avatar_image_path" json:"avatarImagePath" validate:"required"`
	BannerImagePath string `bson:"banner_image_path" json:"bannerImagePath" validate:"required"`
	Biography       string `bson:"biography" json:"biography" validate:"omitempty,min=5,max=255"`
	IsVerified      bool   `bson:"is_verified" json:"isVerified" validate:"boolean"`
}
