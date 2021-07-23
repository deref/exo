# Pagination

`exo` uses cursor-based pagination for all of its APIs so that clients can page through results in a consistent and stateless way. The pagination scheme is described below.

### Request parameters:

| Parameter | required? | Description |
| --------- | --------- | ----------- |
| `cursor` | N | A cursor previously returned by another request, which represents a specific location in the collection. |
| `prev` | N | Maximum number of results to return before the cursor. In most cases, the default will be `0`. *Cannot be specified if `next` is specified.*  |
| `next` | N | Maximum number of results to return after the cursor. Each collection will set a reasonable default if not specified. *Cannot be specified if `prev` is specified.* |

Both `prev` and `next` may be constrained to some maximum value on a per-collection basis. In this case, setting a higher value would have no effect.

Collections may specify other mechanism for filtering or ordering results. For example, a collection could allow a `since` timestamp parameter to be specified, which would return a page of results that are at least as recent as the timestamp supplied. However, only the cursor is necessary to retrieve subsequent pages of data. In these cases, the `next` parameter should be honored as the initial limit on the number of results that could be returned.

### Response parameters:

| Parameter | Description |
| --------- | ----------- |
| `prevCursor` | Cursor to be passed in the requests for the previous page of data. |
| `nextCursor` | Cursor to be passed in the requests for the next page of data. |
| `items` | Array of results. |

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
    "prevCursor": "AXrPduJbkSzMEDnDEfCF1A",
    "nextCursor": "AXrPdwWGH8hCVn90aNxSzA",
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
    "cursor": "AXrPdwWGH8hCVn90aNxSzA",
    "next": 500
}
```
