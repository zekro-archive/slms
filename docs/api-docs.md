# SLMS REST API

## Authorization

Generally, every API endpoint request needs to be authorized.

For authorizing your request, you need to send a valid `Basic` authentication token *(which is defined in the servers configuration)* as `Authorization` header.

```
> POST /api/login HTTP/1.1
> Authorization: Basic yourAuthTokenHere
```

If you are requesting the [`POST /api/login`](#session-login) endpoint with a valid Basic Auhtorization header, you will receive a session token. This token can also be used to authenticate against the API instead of the authorization header. This cookie has a lifetime of 10 minutes after the login request and will not be extended after each following request.

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

Rate limits are applied on a per-route and per-connection basis. The rate limit counter are based on a simple [token bucket](https://en.wikipedia.org/wiki/Token_bucket) system.

Information about the current limiter status are passed yb each response in following headers:

| Header Name | Description |
|-------------|-------------|
| `X-RateLimit-Limit` | The absolute ammount of tokens which can be used in one burst. |
| `X-RateLimit-Remaining` | The remaining tokens after this request. |
| `X-RateLimit-Reset` | The UNIX timestamp you need to wait until you can make another request. This value is `0` as long as at least `1` token is available after the request. |

Here you can see an example response:
```
< HTTP/1.1 200 OK
< Date: Wed, 03 Apr 2019 19:24:35 GMT
< Content-Length: 0
< X-Ratelimit-Limit: 3
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1554297886
```

---

## Endpoints

- [Session Login](#session-login)  
  `POST /api/login`

- [Get Short Link List](#get-short-link-list)  
  `GET /api/shortlinks`

- [Create Short Link](#create-short-link)  
  `POST /api/shortlinks`

- [Get Short Link](#get-short-link)  
  `GET /api/shortlinks/:ID`

- [Modify Short Link](#modify-short-link)  
  `POST /api/shortlinks/:ID`

- [Delete Short Link](#delete-short-link)  
  `DELETE /api/shortlinks/:ID`



### Session Login

> POST /api/login

```
< HTTP/1.1 200 OK
< Date: Tue, 02 Apr 2019 20:44:25 GMT
< Content-Length: 0
< Set-Cookie: session=MTU1NDIzNzg2NX...; expires=Tue, 02 Apr 2019 20:54:25 GMT; path=/
```

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