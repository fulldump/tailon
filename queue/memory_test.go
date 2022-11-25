package queue

import (
	"testing"

	"github.com/fulldump/biff"
)

func TestMemoryService_ListQueues_Empty(t *testing.T) {

	s := NewMemoryService()

	result, err := s.ListQueues()
	biff.AssertNil(err)
	biff.AssertEqual(result, []string{})
}

func TestMemoryService_ListQueues_OneItem(t *testing.T) {

	s := NewMemoryService()
	s.Queues["my-queue"] = NewMemoryQueue()

	result, err := s.ListQueues()
	biff.AssertNil(err)

	expected := []string{"my-queue"}
	biff.AssertEqualJson(result, expected)
}

func TestMemoryService_CreateQueue(t *testing.T) {

	s := NewMemoryService()

	queue, err := s.CreateQueue("my-queue")
	biff.AssertNil(err)
	biff.AssertNotNil(queue)
	biff.AssertEqual(len(s.Queues), 1)
}

func TestMemoryService_CreateQueue_Twice(t *testing.T) {

	s := NewMemoryService()

	queue, err := s.CreateQueue("my-queue")
	biff.AssertNil(err)
	biff.AssertNotNil(queue)
	biff.AssertEqual(len(s.Queues), 1)

	queue, err = s.CreateQueue("my-queue")
	biff.AssertNotNil(err)
	biff.AssertNil(queue)
	biff.AssertEqual(len(s.Queues), 1)
}

func TestMemoryService_GetQueue(t *testing.T) {

	s := NewMemoryService()

	q1, err := s.CreateQueue("my-queue")
	biff.AssertNil(err)

	q2, err := s.GetQueue("my-queue")
	biff.AssertNil(err)

	biff.AssertEqual(q1, q2)
}

func TestMemoryService_GetQueue_NotExist(t *testing.T) {

	s := NewMemoryService()

	q, err := s.GetQueue("invented-queue")
	biff.AssertNotNil(err)
	biff.AssertNil(q)
}

func TestMemoryService_DeleteQueue(t *testing.T) {

	s := NewMemoryService()

	q, err := s.CreateQueue("my-queue")
	biff.AssertNil(err)
	biff.AssertNotNil(q)

	err = s.DeleteQueue("my-queue")
	biff.AssertNil(err)
	biff.AssertEqual(len(s.Queues), 0)
}

func TestMemoryService_DeleteQueue_NotExist(t *testing.T) {

	s := NewMemoryService()

	err := s.DeleteQueue("my-queue")
	biff.AssertNotNil(err)
}

func TestMemoryQueue_WriteAndRead(t *testing.T) {

	q := NewMemoryQueue()

	errWrite := q.Write(JSON(`{"my":"object"}`))
	biff.AssertNil(errWrite)

	item, errRead := q.Read()
	biff.AssertNil(errRead)
	biff.AssertEqualJson(item, map[string]interface{}{"my": "object"})
}
