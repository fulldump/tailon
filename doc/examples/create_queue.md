# Create queue

Curl example:

```sh
curl -X POST "https://example.com/v1/queues" \
-d '{
    "name": "my-queue"
}'
```


HTTP request/response example:

```http
POST /v1/queues HTTP/1.1
Host: example.com

{
    "name": "my-queue"
}

HTTP/1.1 201 Created
Content-Length: 0
Date: Mon, 15 Aug 2022 02:08:13 GMT


```


