package lms

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type LMSService struct {
	repo LMSRepository
	ch   *amqp.Channel
}

func NewLMSService(repo LMSRepository, ch *amqp.Channel) *LMSService {
	return &LMSService{repo: repo, ch: ch}
}

func (s *LMSService) CheckAndUpdateVersion(name, path string) {
	lms, err := s.repo.GetByName(name)
	if err != nil {
		// not exist -> create
		lms = &LMS{Name: name, Path: path}
		if err := s.repo.Create(lms); err != nil {
			log.Printf("Failed to create LMS %s: %v", name, err)
			return
		}
		log.Printf("Created LMS %s", name)
		return
	}
}
