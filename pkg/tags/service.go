package tags

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertTag(tag *models.Tag) (*models.Tag, error, bool)
	FetchTagById(tagId *models.TagById) (*models.Tag, error)
	FetchTagByName(tagName *models.TagByName) (*models.Tag, error)
	RemoveTag(tag *models.Tag) (*models.Tag, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertTag(tag *models.Tag) (*models.Tag, error, bool) {
	return s.repository.CreateTag(tag)
}

func (s *service) FetchTagById(tagId *models.TagById) (*models.Tag, error) {
	return s.repository.ReadTagById(tagId)
}

func (s *service) FetchTagByName(tagName *models.TagByName) (*models.Tag, error) {
	return s.repository.ReadTagByTagName(tagName)
}

func (s *service) RemoveTag(tag *models.Tag) (*models.Tag, error) {
	return s.repository.DeleteTag(tag)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
