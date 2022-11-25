package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/fulldump/box"

	"github.com/fulldump/tailon/glueauth"
	"github.com/fulldump/tailon/queue"
	"github.com/fulldump/tailon/statics"
)

func InjectQueueService(qs queue.Service) box.I {
	return func(next box.H) box.H {
		return func(ctx context.Context) {
			next(SetQueueService(ctx, qs))
		}
	}
}

const QueueServiceKey = "6c6b0b2a-6c5b-11ed-b52b-e78c8b06a360"

func SetQueueService(ctx context.Context, qs queue.Service) context.Context {
	return context.WithValue(ctx, QueueServiceKey, qs)
}

func GetQueueService(ctx context.Context) queue.Service {
	return ctx.Value(QueueServiceKey).(queue.Service)
}

func Build(version, staticsDir string, qs queue.Service) *box.B {

	b := box.NewBox()

	v1 := b.Resource("/v1")

	v1.Resource("/queues").
		WithInterceptors(
			InjectQueueService(qs),
		).
		WithActions(
			box.Get(ListQueues),
			box.Post(CreateQueue),
		)

	v1.Resource("/queues/{queue_id}").
		WithActions(
			box.Get(RetrieveQueue),
			box.Delete(func() string {
				return "delete queue"
			}),
			box.Action(Read),
			box.ActionPost(Write),
		)

	b.Resource("/release").
		WithActions(box.Get(func() string {
			return version
		}))

	b.Resource("/me").
		WithInterceptors(glueauth.Require).
		WithActions(box.Get(func(ctx context.Context) *glueauth.GlueAuthentication {
			return glueauth.GetAuth(ctx)
		}))

	// Mount statics
	b.Resource("/*").
		WithActions(
			box.Get(statics.ServeStatics(staticsDir)).WithName("serveStatics"),
		)

	return b
}

func ListQueues(ctx context.Context) (interface{}, error) {
	s := GetQueueService(ctx)
	return s.ListQueues()
}

type CreateQueueInput struct {
	Name string `json:"name"`
}

func CreateQueue(ctx context.Context, input CreateQueueInput, w http.ResponseWriter) error {

	s := GetQueueService(ctx)

	_, err := s.CreateQueue(input.Name)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func RetrieveQueue(ctx context.Context, w http.ResponseWriter) (interface{}, error) {

	queueName := box.GetUrlParameter(ctx, "queue_id")

	s := GetQueueService(ctx)

	_, err := s.GetQueue(queueName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // todo: check required!!
		return nil, err
	}

	return queueName, nil
}

func Write(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// duplicated code:
	queueName := box.GetUrlParameter(ctx, "queue_id")
	s := GetQueueService(ctx)
	q, err := s.GetQueue(queueName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // todo: check required!!
		return err
	}

	j := json.NewDecoder(r.Body)

	for {
		message := queue.JSON{}

		err := j.Decode(&message)
		if err == io.EOF {
			return nil // all is ok
		}
		if err != nil {
			return err // some error decoding
		}

		err = q.Write(message)
		if err != nil {
			return err // somme error writting to queue
		}
	}

}

func Read(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// duplicated code:
	queueName := box.GetUrlParameter(ctx, "queue_id")
	s := GetQueueService(ctx)
	q, err := s.GetQueue(queueName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // todo: check required!!
		return err
	}

	// get limit
	limit := 1000 // Default limit
	if l, err := strconv.Atoi(r.Header.Get("Limit")); err == nil {
		limit = l
	}

	// j := json.NewEncoder(w)

	// f, isFlusher := w.(http.Flusher)

	for limit > 0 {
		limit--

		message, err := q.Read()
		if err != nil {
			return err // some error reading queue
		}

		w.Write(message)
		w.Write([]byte("\n"))

		// err = j.Encode(message)
		// if err != nil {
		// 	return err // some error encoding response
		// }

		// if isFlusher {
		// 	// f.Flush()
		// }
	}

	return nil
}
