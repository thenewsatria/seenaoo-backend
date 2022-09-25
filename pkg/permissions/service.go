package permissions

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertPermission(permission *models.Permission) (*models.Permission, error)
	FetchPermissionsByItemCategory(permissionItemCategory *models.PermissionByItemCategory) (*[]models.Permission, error)
	FetchPermissionById(permissionId *models.PermissionById) (*models.Permission, error)
	FetchAllPermissions() (*[]models.Permission, error)
	FetchPermissionsDistinctItemCategory() (*[]string, error)
	FetchPermissionByName(permissionName *models.PermissionByName) (*models.Permission, error)
	UpdatePermission(permission *models.Permission) (*models.Permission, error)
	RemovePermission(permission *models.Permission) (*models.Permission, error)
}

type service struct {
	repository Repository
}

func (s *service) FetchPermissionById(permissionId *models.PermissionById) (*models.Permission, error) {
	return s.repository.ReadPermissionById(permissionId)
}

func (s *service) FetchPermissionByName(permissionName *models.PermissionByName) (*models.Permission, error) {
	return s.repository.ReadPermissionByName(permissionName)
}

func (s *service) FetchPermissionsByItemCategory(permissionItemCategory *models.PermissionByItemCategory) (*[]models.Permission, error) {
	return s.repository.ReadPermissionsByItemCategory(permissionItemCategory)
}

func (s *service) InsertPermission(permission *models.Permission) (*models.Permission, error) {
	return s.repository.CreatePermission(permission)
}

func (s *service) RemovePermission(permission *models.Permission) (*models.Permission, error) {
	return s.repository.DeletePermission(permission)
}

func (s *service) UpdatePermission(permission *models.Permission) (*models.Permission, error) {
	return s.repository.UpdatePermission(permission)
}

func (s *service) FetchAllPermissions() (*[]models.Permission, error) {
	return s.repository.ReadAllPermissions()
}

func (s *service) FetchPermissionsDistinctItemCategory() (*[]string, error) {
	return s.repository.ReadPermissionsDistinctItemCategory()
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
