# SLMS REST API

## Parameters

Parameters with `default` formated names are **required** and parameters with *`italic`* formated names are ***optional***.

The parameter type consists of   
1. the passing method 
   - `path`: as URL path subdirectory
   - `query`: as URL query
   - `json-body`: as json formated key-value pair in the request body
2. and the data type:
   - `int`: numeral value
   - `string`: character sequence
   <!-- - `bool`: boolean -->

## Status Codes and Error Formats

The API uses the standart HTTP status codes, as defined in [RFC 2616](https://tools.ietf.org/html/rfc2616#section-10) and [RFC 6585](https://tools.ietf.org/html/rfc6585).

An error response from the API contains the status code as header and an error description as JSON body. Example error response:

```
< HTTP/1.1 400 Bad Request
< Content-Type: application/json
```
```json
{
  "error": {
    "code": 400,
    "message": "json: cannot unmarshal string into Go struct field authRequestData.session of type int"
  }
}
```

## Rate Limits
<!-- TODO: Write stuff. -->

---

## Endpoints

### Get Short Link List

> GET /api/shortlinks

*The list of short links are ordered descending by `created` date.*

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| *`from`* | `query`: `int` | Start item of the list (from top). |
| *`limit`* | `query`: `int` | Maximum ammount of items in list. |

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:16:49 GMT
< Content-Type: application/json
< Content-Length: 715
```
```json
{
  "n": 3,
  "results": [
    {
      "id": 3,
      "root_link": "https://someurl.example/somedoc.txt",
      "short_link": "RyajU4cH",
      "created": "2019-03-04T18:42:07Z",
      "accesses": 2,
      "edited": "2019-03-04T17:43:26Z"
    },
    {
      "id": 2,
      "root_link": "https://github.com/zekroTJA/vplan2019/tree/dev",
      "short_link": "vplan2",
      "created": "2019-03-04T08:56:09Z",
      "accesses": 32,
      "edited": "2019-03-04T08:56:15Z"
    },
    {
      "id": 1,
      "root_link": "https://github.com/zekroTJA/shinpuru/releases/tag/0.9.0",
      "short_link": "sp09",
      "created": "2019-02-23T11:02:37Z",
      "accesses": 12,
      "edited": "2019-03-04T00:37:02Z"
    }
  ]
}
```

#### Response

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:10:11 GMT
< Content-Type: application/json
< Content-Length: 179
```
```json
{
  "id": 12,
  "root_link": "https://github.com/zekroTJA/slms",
  "short_link": "slms",
  "created": "2018-10-07T12:20:36Z",
  "accesses": 9,
  "edited": "2019-04-02T09:11:19Z"
}
```

---

### Get Short Link

> GET /api/shortlinks/:ID

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `ID` | `path`: `string` | The unique ID *or* the short identifier of the short link. |

#### Response

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:10:11 GMT
< Content-Type: application/json
< Content-Length: 179
```
```json
{
  "id": 12,
  "root_link": "https://github.com/zekroTJA/slms",
  "short_link": "slms",
  "created": "2018-10-07T12:20:36Z",
  "accesses": 9,
  "edited": "2019-04-02T09:11:19Z"
}
```

---

### Create Short Link

> POST /api/shortlinks

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `root_link` | `json-body`: `string` | The root link. |
| *`short_link`* | `json-body`: `string` | The short link identifier.<br>If this argument is not passed, a new identifier of random characters will be created. |

#### Response

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:23:57 GMT
< Content-Type: application/json
< Content-Length: 166
```
```json
{
  "id": 23,
  "root_link": "http://zekro.de",
  "short_link": "B2ffM7Tk",
  "created": "2019-04-02T22:24:02Z",
  "accesses": 0,
  "edited": "2019-04-02T22:24:02Z"
}
```

---

### Modify Short Link

> POST /api/shortlinks/:ID

*You can only modify the `root_link` **or** the `short_link` in one request.*

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `ID` | `path`: `string` | The unique ID *or* the short identifier of the short link. |
| *`root_link`* | `json-body`: `string` | Pass this to modify the root link. |
| *`short_link`* | `json-body`: `string` | Pas this to modify the short identifier. |

#### Response

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:31:35 GMT
< Content-Type: application/json
< Content-Length: 180
```
```json
{
  "id": 23,
  "root_link": "https://zekro.de/src/logo.png",
  "short_link": "B2ffM7Tk",
  "created": "2019-04-02T22:24:02Z",
  "accesses": 0,
  "edited": "2019-04-02T22:24:02Z"
}
```

---

### Delete Short Link

> DELETE /api/shortlinks/:ID

*This endpoint **only** fetches the short link by **ID** and not also by short identifier.*

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `ID` | `path`: `int` | The unique ID of the short link. |

#### Response

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:33:21 GMT
< Content-Length: 0
```