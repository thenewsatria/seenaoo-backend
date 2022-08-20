package refreshtokens

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertRefreshToken(refreshToken *models.RefreshToken) (*models.RefreshToken, error)
	FetchRefreshTokenByUsername(refreshTokenUsername *models.RefreshTokenByUsersUsernameRequest) (*models.RefreshToken, error)
	UpdateRefreshToken(refreshToken *models.RefreshToken) (*models.RefreshToken, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertRefreshToken(refreshToken *models.RefreshToken) (*models.RefreshToken, error) {
	return s.repository.CreateRefreshToken(refreshToken)
}

func (s *service) FetchRefreshTokenByUsername(refreshTokenUsername *models.RefreshTokenByUsersUsernameRequest) (*models.RefreshToken, error) {
	return s.repository.ReadRefreshTokenByUsersUsername(refreshTokenUsername)
}

func (s *service) UpdateRefreshToken(refreshToken *models.RefreshToken) (*models.RefreshToken, error) {
	return s.repository.UpdateRefreshToken(refreshToken)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
