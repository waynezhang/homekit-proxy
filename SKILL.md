---
name: homekit-proxy
description: A skill to communicate with HomeKit Proxy
metadata:
  openclaw:
    requires:
      env:
        - HOMEKIT_PROXY_BASE_URL
        - HOMEKIT_PROXY_AUTH_TOKEN
---

# HomeKit Proxy — AI Agent Skill Reference

## Overview

HomeKit Proxy exposes a local HTTP API to read and control HomeKit accessories and automations. Use this skill to query device state and issue control commands.

Base URL: `http://host:port`

---

## Authentication

If the server has auth enabled, include a bearer token on all `/s/*` requests:

```
Authorization: Bearer <token>
```

If no token is set server-side, omit the header.

---

## Skill: List All Devices and Automations

**When to use:** To discover available accessories (characteristics) and automations, or to look up IDs before issuing a control command.

**Request:**
```
GET /s/all
```

**Key response fields:**

| Field | Description |
|---|---|
| `now` | Server's current time (ISO 8601) |
| `name` | Bridge name |
| `characteristics[]` | List of controllable accessory attributes |
| `characteristics[].id` | Integer ID used for control commands |
| `characteristics[].name` | Human-readable label (e.g. `"Living Room - on"`) |
| `characteristics[].type` | Attribute type (e.g. `"on"`, `"brightness"`, `"temperature"`) |
| `characteristics[].value` | Current value as a string |
| `characteristics[].min/max/step` | Constraints for numeric values |
| `characteristics[].area` | Room or area grouping |
| `automations[]` | List of scheduled automations |
| `automations[].id` | Integer ID used to enable/disable |
| `automations[].name` | Human-readable label |
| `automations[].cron` | Schedule in cron format |
| `automations[].enabled` | Whether the automation is active |
| `automations[].next_run` | ISO 8601 timestamp; zero-value when disabled |

---

## Skill: Control a Device Characteristic

**When to use:** To turn a device on/off, change brightness, set temperature, etc.

**Steps:**
1. Call `GET /s/all` to find the characteristic `id` for the target device/attribute.
2. Issue the control command using that `id`.

**Request:**
```
POST /s/c/{id}
Content-Type: application/json

{ "value": "<string>" }
```

**Value encoding:**
- Boolean attributes: `"true"` or `"false"`
- Numeric attributes: `"22.5"`, `"80"`, etc. (always a string)

**Success response:**
```json
{ "result": "OK" }
```

**Example — turn on Living Room light (id: 1):**
```
POST /s/c/1
{ "value": "true" }
```

---

## Skill: Enable or Disable an Automation

**When to use:** To pause or resume a scheduled automation.

**Steps:**
1. Call `GET /s/all` to find the automation `id`.
2. Send the enable/disable command.

**Request:**
```
POST /s/a/{id}
Content-Type: application/json

{ "value": "true" }   // "true"/"false"/"1"/"0"
```

**Success response:**
```json
{ "result": "OK" }
```

State is persisted and survives server restarts.

---

## Error Handling

| HTTP Status | Meaning |
|---|---|
| `400` | Bad request — invalid ID, malformed JSON, non-boolean value, wrong method |
| `401` | Unauthorized — missing or invalid bearer token |
| `500` | Server error — command execution failed, database error |

On error, retry only if the cause is transient (e.g. network). Do not retry `400` errors without correcting the request.

---

## Decision Guide

```
User wants to know device state or find an ID
  → GET /s/all

User wants to control a device (on/off, brightness, temperature, etc.)
  → GET /s/all (find characteristic id)
  → POST /s/c/{id} with { "value": "..." }

User wants to enable/disable a scheduled automation
  → GET /s/all (find automation id)
  → POST /s/a/{id} with { "value": "true" or "false" }
```
