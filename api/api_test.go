package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/fulldump/apitest"
	"github.com/fulldump/biff"

	"github.com/fulldump/tailon/queue"
)

type JSON = map[string]interface{}

func TestAcceptance(t *testing.T) {

	biff.Alternative("Setup", func(a *biff.A) {

		qs := queue.NewMemoryService()

		h := Build("test version", "", qs)

		api := apitest.NewWithHandler(h)

		biff.Alternative("List queues", func(a *biff.A) {
			res := api.Request("GET", "/v1/queues").Do()

			biff.AssertEqualJson(res.BodyJson(), []string{})
		})

		biff.Alternative("Create queue", func(a *biff.A) {
			res := api.Request("POST", "/v1/queues").WithBodyJson(JSON{
				"name": "my-queue",
			}).Do()
			Save(res, "Create queue", ``)

			biff.AssertEqual(res.BodyString(), "")
			biff.AssertEqual(res.StatusCode, http.StatusCreated)

			biff.Alternative("List queues", func(a *biff.A) {
				res := api.Request("GET", "/v1/queues").Do()
				Save(res, "List queues", ``)

				biff.AssertEqualJson(res.BodyJson(), []string{"my-queue"})
			})
			biff.Alternative("Retrieve queue", func(a *biff.A) {
				res := api.Request("GET", "/v1/queues/my-queue").Do()
				Save(res, "Retrieve queue", ``)

				biff.AssertEqual(res.StatusCode, http.StatusOK)
				biff.AssertEqualJson(res.BodyJson(), JSON{
					"name": "my-queue",
					"len":  0,
				})
			})
			biff.Alternative("Write messages", func(a *biff.A) {

				body := strings.Join([]string{
					`{"id":1,"message":"element 1"}`,
					`{"id":2,"message":"element 2"}`,
					`{"id":3,"message":"element 3"}`,
				}, "\n")

				res := api.Request("POST", "/v1/queues/my-queue:write").
					WithBodyString(body).Do()
				Save(res, "Write messages", ``)

				biff.AssertEqual(res.StatusCode, http.StatusOK)
				biff.AssertEqual(res.BodyString(), "")

				biff.Alternative("Read messages", func(a *biff.A) {
					res := api.Request("GET", "/v1/queues/my-queue:read").
						WithHeader("Limit", "3").Do()
					Save(res, "Read messages", ``)

					dec := json.NewDecoder(strings.NewReader(res.BodyString()))

					for i := 1; i <= 3; i++ {
						m := JSON{}
						dec.Decode(&m)
						biff.AssertEqualJson(m, JSON{"id": i, "message": "element " + strconv.Itoa(i)})
					}
				})
			})

		})

	})

}
