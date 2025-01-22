package service

import "fmt"

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	return &Service{config: config}
}

func (s *Service) ListenAndServe() error {
	fmt.Println("start listening")
	return nil
}
