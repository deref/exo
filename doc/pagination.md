# Pagination

`exo` uses cursor-based pagination for all of its APIs so that clients can page through results in a consistent and stateless way. The pagination scheme is described below.

### Request parameters:

| Parameter | required? | Description |
| --------- | --------- | ----------- |
| `cursor` | N | A cursor previously returned by another request, which represents a specific location in the collection |
| `previous` | N | Maximum number of results to return before the cursor. In most cases, the default will be `0`. |
| `next` | N | Maximum number of results to return at or after the cursor. Each resource type will set a reasonable default, which is not guaranteed to be consistent between different resource types. |

Collections may specify other mechanism for filtering or ordering results. For example, logs allow a `since` timestamp parameter to be specified, which will return a page of results that are at least as recent as the timestamp supplied. However, only the cursor is necessary to retrieve subsequent pages of data. In these cases, the `next` parameter should be honored as the initial limit on the number of results that could be returned.

### Response parameters:

| Parameter | Description |
| --------- | ----------- |
| `prev`    | Cursor to be passed in the requests for the previous page of data. |
| `next`    | Cursor to be passed in the requests for the next page of data. |
| `data`    | Array of results. |


## Example:

```
POST /_exo/workspace/get-events?id=abc-123

{
    "ref": "webserver",
    "after": "2021-07-21 16:40:00.4965432",
    "next": 500
}

---

{
    "prev": "01fb5kjft4m19sb4bapr4y2mk8",
    "next": "01fb5khvq2qvh9tv36bhge2w3s",
    "data": [
        ...
    ]
}
```

To obtain the next page of results, the client could perform the following request:

```
POST /_exo/workspace/get-events?id=abc-123

{
    "ref": "webserver",
    "cursor": "01fb5khvq2qvh9tv36bhge2w3s",
    "next": 500
}
```
