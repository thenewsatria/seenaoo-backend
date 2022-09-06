package collaborations

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error)
	GetCollaboration(collaborationId *models.CollaborationById) (*models.Collaboration, error)
	UpdateCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error)
	RemoveCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error)
}

type service struct {
	repository Repository
}

func (s *service) GetCollaboration(collaborationId *models.CollaborationById) (*models.Collaboration, error) {
	return s.repository.ReadCollaboration(collaborationId)
}

func (s *service) InsertCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error) {
	return s.repository.CreateCollaboration(collaboration)
}

func (s *service) RemoveCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error) {
	return s.repository.DeleteCollaboration(collaboration)
}

func (s *service) UpdateCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error) {
	return s.repository.UpdateCollaboration(collaboration)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
