package service

import "segment_service/domain"

type SegmentRepository interface {
	GetAll() []domain.Segment
	Delete(segment string) error
	Create(segment string) error
	Distribute(segment string, percentage float64) error
}

type SegmentService struct {
	segmentRepository SegmentRepository
}

func NewSegmentService(sr SegmentRepository) *SegmentService {
	return &SegmentService{
		segmentRepository: sr,
	}
}

func (ss *SegmentService) GetAll() []domain.Segment {
	return ss.segmentRepository.GetAll()
}

func (ss *SegmentService) Delete(segment string) error {
	return ss.segmentRepository.Delete(segment)
}

func (ss *SegmentService) Create(segment string) error {
	return ss.segmentRepository.Create(segment)
}

func (ss *SegmentService) Distribute(segment string, percentage float64) error {
	if percentage < 0 || percentage > 100 {
		return domain.ErrInvalidPercentage
	}

	return ss.segmentRepository.Distribute(segment, percentage)
}
