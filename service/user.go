package service

import (
	"segment_service/domain"
)

type UserRepository interface {
	GetAll() (map[domain.User][]domain.Segment, error)
	GetSegmentsById(id int) ([]domain.Segment, error)
	ClearUserSegments(id int) error
	AddSegment(id int, segment string) error
	DeleteSegment(id int, segment string) error
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepository: ur,
	}
}

func (us *UserService) GetAll() (map[domain.User][]domain.Segment, error) {
	return us.userRepository.GetAll()
}

func (us *UserService) GetSegmentsById(id int) ([]domain.Segment, error) {
	return us.userRepository.GetSegmentsById(id)
}

func (us *UserService) ClearUserSegments(id int) error {
	return us.userRepository.ClearUserSegments(id)
}

func (us *UserService) AddSegment(id int, segment string) error {
	return us.userRepository.AddSegment(id, segment)
}

func (us *UserService) DeleteSegment(id int, segment string) error {
	return us.userRepository.DeleteSegment(id, segment)
}
