# Write messages

Curl example:

```sh
curl -X POST "https://example.com/v1/queues/my-queue:write" \
-d '{"id":1,"message":"element 1"}
{"id":2,"message":"element 2"}
{"id":3,"message":"element 3"}'
```


HTTP request/response example:

```http
POST /v1/queues/my-queue:write HTTP/1.1
Host: example.com

{"id":1,"message":"element 1"}
{"id":2,"message":"element 2"}
{"id":3,"message":"element 3"}

HTTP/1.1 200 OK
Content-Length: 0
Date: Mon, 15 Aug 2022 02:08:13 GMT


```


