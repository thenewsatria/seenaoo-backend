package collaborations

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error, bool)
	FetchCollaboration(collaborationId *models.CollaborationById) (*models.Collaboration, error)
	UpdateCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error, bool)
	RemoveCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error)
	CheckIsCollaborator(collaborationItemIdAndCollaborator *models.CollaborationByItemIdAndCollaborator) (bool, error)
	FetchCollaborationByItemIdAndCollaborator(collaborationItemAndCollaborator *models.CollaborationByItemIdAndCollaborator) (*models.Collaboration, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchCollaboration(collaborationId *models.CollaborationById) (*models.Collaboration, error) {
	return s.repository.ReadCollaboration(collaborationId)
}

func (s *service) InsertCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error, bool) {
	return s.repository.CreateCollaboration(collaboration)
}

func (s *service) RemoveCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error) {
	return s.repository.DeleteCollaboration(collaboration)
}

func (s *service) UpdateCollaboration(collaboration *models.Collaboration) (*models.Collaboration, error, bool) {
	return s.repository.UpdateCollaboration(collaboration)
}

func (s *service) CheckIsCollaborator(collaborationItemIdAndCollaborator *models.CollaborationByItemIdAndCollaborator) (bool, error) {
	collaboration, err := s.repository.ReadCollaborationsByItemIdAndCollaborator(collaborationItemIdAndCollaborator)
	if err != nil {
		return false, err
	}
	return collaboration.Status == "ACCEPTED", nil
}

func (s *service) FetchCollaborationByItemIdAndCollaborator(
	collaborationItemAndCollaborator *models.CollaborationByItemIdAndCollaborator) (*models.Collaboration, error) {
	return s.repository.ReadCollaborationsByItemIdAndCollaborator(collaborationItemAndCollaborator)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
