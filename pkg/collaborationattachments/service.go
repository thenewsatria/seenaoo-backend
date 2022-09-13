package collaborationattachments

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertCollaborationAttachment(collabAtt *models.CollaborationAttachment) (*models.CollaborationAttachment, error)
	FetchCollaborationAttachment(collabAttId *models.CollaborationAttachmentById) (*models.CollaborationAttachment, error)
	RemoveCollaborationAttachment(collabAtt *models.CollaborationAttachment) (*models.CollaborationAttachment, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchCollaborationAttachment(collabAttId *models.CollaborationAttachmentById) (*models.CollaborationAttachment, error) {
	return s.repository.ReadCollabAttachment(collabAttId)
}

func (s *service) InsertCollaborationAttachment(collabAtt *models.CollaborationAttachment) (*models.CollaborationAttachment, error) {
	return s.repository.CreateCollabAttachment(collabAtt)
}

func (s *service) RemoveCollaborationAttachment(collabAtt *models.CollaborationAttachment) (*models.CollaborationAttachment, error) {
	return s.repository.DeleteCollabAttachment(collabAtt)
}

func NewFunc(r Repository) Service {
	return &service{
		repository: r,
	}
}
