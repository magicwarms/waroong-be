package user_types

import (
	"waroong-be/apps/user_types/interfaces"
	"waroong-be/apps/user_types/model"
)

type userTypeService struct {
	userTypeRepository interfaces.UserTypeRepository
}

// NewService is used to create a single instance of the service
func NewService(r interfaces.UserTypeRepository) interfaces.UserTypeService {
	return &userTypeService{
		userTypeRepository: r,
	}
}

func (s *userTypeService) GetUserTypeById(id uint) (*model.UserTypeModel, error) {
	// TODO:
	// use redis instead of DB
	getUserType, err := s.userTypeRepository.GetById(id)
	if err != nil {
		return &model.UserTypeModel{}, err
	}

	return getUserType, nil
}
