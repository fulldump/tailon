package api

import (
	"testing"

	"github.com/fulldump/apitest"
	"github.com/fulldump/biff"

	"github.com/fulldump/tailon/queue"
)

func TestAcceptance(t *testing.T) {

	biff.Alternative("Setup", func(a *biff.A) {

		qs := queue.NewMemoryService()

		h := Build("test version", "", qs)

		api := apitest.NewWithHandler(h)

		biff.Alternative("List queues", func(a *biff.A) {
			res := api.Request("GET", "/v1/queues").Do()
			Save(res, "List queues", ``)

			biff.AssertEqual(res.BodyJson(), "list queues")
		})

		biff.Alternative("Create queue", func(a *biff.A) {
			res := api.Request("POST", "/v1/queues").Do()
			Save(res, "Create queue", ``)

			biff.AssertEqual(res.BodyJson(), "create queue")
		})

	})

}
