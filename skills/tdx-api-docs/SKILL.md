---
name: tdx-api-docs
description: Use when answering questions about this repository's TDX HTTP or WebSocket API, including endpoints, parameters, response shapes, units, examples, implementation status, or integration guidance based on API_接口文档.md and API_集成指南.md
---

# TDX API Docs

## Overview

Use this skill when the user is asking how to call the repo's API, what an endpoint returns, whether an endpoint exists, or how the documented behavior maps to the implementation.

## Sources Of Truth

Read these files in this order:

1. `API_接口文档.md`
2. `API_集成指南.md`
3. `web/server.go`
4. `web/server_api_extended.go`
5. `web/server_ws.go`

If the markdown and code disagree, treat the code as the current behavior and say the markdown was outdated or has been updated.

## Workflow

1. Find the endpoint in `API_接口文档.md`.
2. If the user is asking how to integrate the API, which endpoint group to choose, or timeout/retry guidance, read `API_集成指南.md`.
3. Confirm the route is actually registered in `web/server.go`.
4. If the user asks whether something is implemented, inspect the matching handler in `web/server.go`, `web/server_api_extended.go`, or `web/server_ws.go`.
5. Answer with the exact method, path, required params, optional params, and a minimal example.
6. When useful, call out data units explicitly:
   - price: `厘`
   - volume: `手`
   - order size in quote depth: `股`

## API Notes

- HTTP base URL: `http://${ENV_TDX_API_HOST}:${ENV_TDX_API_PORT}`
- WebSocket base URL: `ws://${ENV_TDX_API_HOST}:${ENV_TDX_API_PORT}`
- Server env vars:
  - `ENV_TDX_API_HOST`: listen host, local default `localhost`, container default `0.0.0.0`
  - `ENV_TDX_API_PORT`: listen port, default `8080`
- Unified JSON envelope:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- WebSocket quote stream endpoint: `GET /ws/quote`
- `code=0` means success, `code=-1` means failure

## Response Guidance

- Prefer concise endpoint-oriented answers over long prose.
- If the user asks for examples, provide `curl`, JavaScript `fetch`, or WebSocket snippets matching the documented interface.
- If the user asks for "有没有实现", "支持吗", or similar, verify in code instead of trusting the markdown alone.
- If the user asks "怎么接", "用哪个接口组合", "超时怎么配", or similar, prefer `API_集成指南.md` before expanding into implementation details.
- If the user asks for batch, index, history, or task APIs, also inspect `web/server_api_extended.go`.

## Do Not

- Do not invent fields that are not present in the doc or handler output.
- Do not assume the API uses standard yuan/share units; this repo frequently returns `厘` and `手`.
- Do not claim a route exists unless it is both documented or requested and registered in the web server code.
