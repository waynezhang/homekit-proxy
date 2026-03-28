# HomeKit Proxy HTTP API

## Authentication

Auth is optional. When `HOMEKIT_PROXY_USER` and `HOMEKIT_PROXY_PASSWORD` are set, all `/s/*` endpoints require authentication.

Set `HOMEKIT_PROXY_AUTH_TOKEN` to a secret value and pass it as a bearer token:

```
Authorization: Bearer <HOMEKIT_PROXY_AUTH_TOKEN>
```

---

## Data

### `GET /s/all`
Returns all accessories and automations.

**Response:**
```json
{
  "now": "2026-03-28T00:00:00Z",
  "name": "bridge name",
  "characteristics": [
    {
      "id": 1,
      "name": "Living Room - on",
      "type": "on",
      "value": "true",
      "min": "",
      "max": "",
      "step": "",
      "area": "living room",
      "icon": ""
    }
  ],
  "automations": [
    {
      "id": 1,
      "name": "Morning Routine",
      "cmd": "/path/to/script.sh",
      "cron": "0 7 * * *",
      "tolerance": 0,
      "last_run": "2026-03-28T07:00:00Z",
      "last_error": "",
      "next_run": "2026-03-29T07:00:00Z",
      "enabled": true,
      "group": "morning"
    }
  ]
}
```
Note: `next_run` is zero-value when automation is disabled.

---

## Control

### `POST /s/c/{id}`
Updates a characteristic value. `{id}` is the integer characteristic ID from `/s/all`.

**Request:**
```json
{ "value": "true" }
```
Value is a string regardless of characteristic type (e.g. `"22.5"` for temperature, `"true"` for boolean).

**Response:**
```json
{ "result": "OK" }
```

---

### `POST /s/a/{id}`
Enables or disables an automation. State is persisted to the action log DB and survives restarts.

**Request:**
```json
{ "value": "true" }
```
`value` must be parseable as a boolean (`"true"`, `"false"`, `"1"`, `"0"`).

**Response:**
```json
{ "result": "OK" }
```

---

## Error Responses

| Status | Condition |
|---|---|
| `400` | Invalid method, malformed body, invalid ID, non-boolean value |
| `401` | Missing or invalid bearer token (protected endpoints) |
| `500` | Command execution failure, DB error |
