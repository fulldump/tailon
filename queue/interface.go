package queue

import (
	"encoding/json"
)

type JSON = json.RawMessage

type Queue interface {
	Write(JSON) error
	Read() (JSON, error)
}

type Info struct {
	Name  string
	Queue Queue
}

type Service interface {
	GetQueue(name string) (Queue, error)
	ListQueues() ([]string, error)
	CreateQueue(name string) (Queue, error)
	DeleteQueue(name string) error
}
