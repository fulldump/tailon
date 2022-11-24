package api

import (
	"fmt"

	"github.com/fulldump/box"

	"github.com/fulldump/tailon/statics"
)

func Build(version, staticsDir string) *box.B {

	b := box.NewBox()

	v1 := b.Resource("/v1")

	v1.Resource("/queues").
		WithActions(
			box.Get(func() (string, error) {
				return "list queues", fmt.Errorf("soy un error")
			}),
			box.Post(func() string {
				return "create queue"
			}),
		)

	v1.Resource("/queues/{queue_id}").
		WithActions(
			box.Get(func() string {
				return "return queue"
			}),
			box.Delete(func() string {
				return "delete queue"
			}),
			box.Action(func() string {
				return "read from queue"
			}).WithName("read"),
			box.ActionPost(func() string {
				return "write to queue"
			}).WithName("write"),
		)

	b.Resource("/release").
		WithActions(box.Get(func() string {
			return version
		}))

	// Mount statics
	b.Resource("/*").
		WithActions(
			box.Get(statics.ServeStatics(staticsDir)).WithName("serveStatics"),
		)

	return b
}
