# SLMS REST API

## Endpoints

### Get Short Link

> GET /api/shortlinks/:ID

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `ID` | `string` | The unique ID *or* the short identifyer of the short link |

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```