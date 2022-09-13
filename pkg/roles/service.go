package roles

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertRole(role *models.Role) (*models.Role, error)
	FetchRoleById(roleId *models.RoleById) (*models.Role, error)
	FetchRoleBySlug(roleSlug *models.RoleBySlug) (*models.Role, error)
	FetchRolesByOwner(roleOwner *models.RoleByOwner) (*[]models.Role, error)
	UpdateRole(role *models.Role) (*models.Role, error)
	DeleteRole(role *models.Role) (*models.Role, error)
}

type service struct {
	repository Repository
}

func (s *service) DeleteRole(role *models.Role) (*models.Role, error) {
	return s.repository.DeleteRole(role)
}

func (s *service) FetchRoleById(roleId *models.RoleById) (*models.Role, error) {
	return s.repository.ReadRoleById(roleId)
}

// FetchRoleBySlug implements Service
func (s *service) FetchRoleBySlug(roleSlug *models.RoleBySlug) (*models.Role, error) {
	return s.repository.ReadRoleBySlug(roleSlug)
}

// FetchRolesByOwner implements Service
func (s *service) FetchRolesByOwner(roleOwner *models.RoleByOwner) (*[]models.Role, error) {
	return s.repository.ReadRolesByOwner(roleOwner)
}

// InsertRole implements Service
func (s *service) InsertRole(role *models.Role) (*models.Role, error) {
	return s.repository.CreateRole(role)
}

// UpdateRole implements Service
func (s *service) UpdateRole(role *models.Role) (*models.Role, error) {
	return s.repository.UpdateRole(role)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
