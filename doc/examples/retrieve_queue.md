# Retrieve queue

Curl example:

```sh
curl "https://example.com/v1/queues/my-queue"
```


HTTP request/response example:

```http
GET /v1/queues/my-queue HTTP/1.1
Host: example.com



HTTP/1.1 200 OK
Content-Length: 49
Content-Type: text/plain; charset=utf-8
Date: Mon, 15 Aug 2022 02:08:13 GMT

{
    "len": 0,
    "name": "my-queue",
    "reads": 0,
    "writes": 0
}
```


