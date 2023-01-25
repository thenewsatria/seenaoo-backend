package userprofiles

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertProfile(userProfile *models.UserProfile) (*models.UserProfile, error, bool)
	FetchProfileById(userProfileId *models.UserProfileByID) (*models.UserProfile, error)
	FetchProfileByOwner(userProfileOwner *models.UserProfileByOwner) (*models.UserProfile, error)
	UpdateProfile(userProfile *models.UserProfile) (*models.UserProfile, error, bool)
}

type service struct {
	repository Repository
}

// FetchProfileById implements Service
func (s *service) FetchProfileById(userProfileId *models.UserProfileByID) (*models.UserProfile, error) {
	return s.repository.ReadProfileById(userProfileId)
}

// FetchProfileByOwner implements Service
func (s *service) FetchProfileByOwner(userProfileOwner *models.UserProfileByOwner) (*models.UserProfile, error) {
	return s.repository.ReadProfileByOwner(userProfileOwner)
}

// InsertProfile implements Service
func (s *service) InsertProfile(userProfile *models.UserProfile) (*models.UserProfile, error, bool) {
	return s.repository.CreateProfile(userProfile)
}

func (s *service) UpdateProfile(userProfile *models.UserProfile) (*models.UserProfile, error, bool) {
	return s.repository.UpdateProfile(userProfile)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
