# Pagination

`exo` uses cursor-based pagination for all of its APIs so that clients can page through results in a consistent and stateless way. The pagination scheme is described below.

### Request parameters:

| Parameter | required? | Description |
| --------- | --------- | ----------- |
| `cursor` | N | A cursor previously returned by another request, which represents a specific location in the collection |
| `prev` | N | Maximum number of results to return before the cursor. In most cases, the default will be `0`. |
| `next` | N | Maximum number of results to return at or after the cursor. Each resource type will set a reasonable default, which is not guaranteed to be consistent between different resource types. |

Collections may specify other mechanism for filtering or ordering results. For example, logs allow a `since` timestamp parameter to be specified, which will return a page of results that are at least as recent as the timestamp supplied. However, only the cursor is necessary to retrieve subsequent pages of data. In these cases, the `next` parameter should be honored as the initial limit on the number of results that could be returned.

### Response parameters:

| Parameter | Description |
| --------- | ----------- |
| `prevCursor` | Cursor to be passed in the requests for the previous page of data. |
| `nextCursor` | Cursor to be passed in the requests for the next page of data. |
| `data` | Array of results. |

Whereas `prev` results are exclusive of the cursor (i.e. `prev` specifies the number of results _less than_ the cursor), the `next` results are always inclusive of the cursor when `next` is > 0. The implication here is that that server's `prevCursor` can be derived from the lowest-ordered result returned, whereas the `nextCursor` must be derived from some value outside the result set so that the first element of the following request does not include the last element of the previous request. A simple way to ensure this is to take the key from the last element in the result set and increment the value by a single bit.

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
    "prevCursor": "01fb5kjft4m19sb4bapr4y2mk8",  <1>
    "nextCursor": "01fb5khvq2qvh9tv36bhge2w3w",  <2>
    "data": [
        {
            id: "01fb5kjft4m19sb4bapr4y2mk8",    <1>
            ...
        },
        ...
        {
            id: "01fb5khvq2qvh9tv36bhge2w3s",    <2>
            ...
        }
    ]
}
```

1. `prevCursor` is derived from the sort key of the first result item.
2. `nextCursor` is larger than the sort key of the last result item. In this case, its binary representation is one bit greater than the key of the last result item. 

To obtain the next page of results, the client could perform the following request:

```
POST /_exo/workspace/get-events?id=abc-123

{
    "ref": "webserver",
    "cursor": "01fb5khvq2qvh9tv36bhge2w3w",
    "next": 500
}

---

{
    "prevCursor": "01fb5khvq2qvh9tv36bhge2w50",
    "nextCursor": <elided>,
    "data": [
        {
            id: "01fb5khvq2qvh9tv36bhge2w50",
            ...
        },
    ]
}
```

Note that even though the cursor of `01fb5khvq2qvh9tv36bhge2w3w` was supplied, the first sort key present in the dataset that is _greater than or equal to_ the key implied from the cursor is `01fb5khvq2qvh9tv36bhge2w50`.
