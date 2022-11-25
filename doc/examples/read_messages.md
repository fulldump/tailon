# Read messages

Curl example:

```sh
curl "https://example.com/v1/queues/my-queue:read" \
-H "Limit: 3"
```


HTTP request/response example:

```http
GET /v1/queues/my-queue:read HTTP/1.1
Host: example.com
Limit: 3



HTTP/1.1 200 OK
Content-Length: 93
Content-Type: text/plain; charset=utf-8
Date: Mon, 15 Aug 2022 02:08:13 GMT

{"id":1,"message":"element 1"}
{"id":2,"message":"element 2"}
{"id":3,"message":"element 3"}

```


