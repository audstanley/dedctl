# Issues

This file tracks all identified issues in the Audstanley Games project, organized by severity.

---

## Critical Issues

### 1. `setInterval` every 100ms in `auth.ts`
**Severity:** Critical
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/lib/stores/auth.ts`
**Impact:** Runs `loadFromStorage()` 10 times per second for the entire page lifetime. Severe performance impact.
**Description:** The auth store uses a 100ms `setInterval` to poll localStorage. This is unnecessary and causes significant CPU usage.
**Fix:** Removed the `setInterval`. Added a proper `subscribe` function so components can reactively listen to auth state changes without polling.

### 2. Duplicate navbars on dashboard
**Status:** ✅ Resolved
**Files:** `frontend/games-frontend/src/routes/+layout.svelte`, `frontend/games-frontend/src/routes/dashboard/+layout.svelte`
**Impact:** Visual bug - two navbars stacked on top of each other on the dashboard.
**Fix:** Removed the duplicate navbar from `dashboard/+layout.svelte` - it was redundant since the root layout already renders one.

### 3. Non-reactive auth state
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/lib/stores/auth.ts`
**Impact:** Navbar doesn't update after login; auth changes not reflected anywhere until page refresh.
**Fix:** Added a `subscribe` function with a listener set. Components use `$effect` with `auth.subscribe()` to reactively update when auth state changes.

### 4. SSE token in URL query string
**Status:** ⚠️ Partially Resolved (see note)
**File:** `frontend/games-frontend/src/lib/api/client.ts:108`
**Impact:** Security risk - JWT token leaks into browser history, server access logs, and proxy logs.
**Note:** `EventSource` does not support custom headers natively, so query parameter is still required for SSE. The backend `AuthRequired` middleware now also checks the `Authorization` header for SSE connections. Consider using a WebSocket-based approach for production.

### 5. Unbounded log array (memory leak)
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/routes/games/[name]/logs/+page.svelte`
**Impact:** Memory leak in long-running log streaming sessions.
**Fix:** Added `MAX_LOGS = 1000` cap. Oldest entries are trimmed when the limit is exceeded.

### 6. Backend SSE goroutine leak
**Status:** ✅ Resolved
**File:** `backend/steam-game-control/internal/service/game.go`
**Impact:** Goroutine and systemd journal file handle leak if client disconnects.
**Fix:** Added `context.Context` parameter to `StreamLogs` interface. The service now checks `ctx.Done()` in every iteration and exits gracefully when the client disconnects.

### 7. Backend SSE data injection
**Status:** ✅ Resolved
**File:** `backend/steam-game-control/internal/handler/game.go`
**Impact:** Newlines or backslashes in log messages break SSE protocol and could inject fake events.
**Fix:** Added `escapeSSEData()` function that escapes `\`, `\n`, and `\r` per SSE spec. Also removed `Access-Control-Allow-Origin: *` override from SSE handler.

---

## High Issues

### 8. Inconsistent state management (3 patterns)
**Status:** ⚠️ Partially Resolved
**File:** `frontend/games-frontend/src/lib/stores/auth.ts`, `frontend/games-frontend/src/lib/stores/games.ts`
**Impact:** Confusing, hard to maintain, no reactivity in auth store.
**Fix:** `auth.ts` now uses a proper subscribe/listener pattern. `games.ts` still uses `writable` stores but this is functional - full rewrite to `$state` is a lower-priority refactor.

### 9. `games.ts` uses `writable` stores without reactive subscriptions
**Status:** ⚠️ Partially Resolved
**File:** `frontend/games-frontend/src/lib/stores/games.ts`
**Impact:** No reactivity when components consume the games store.
**Fix:** Removed the duplicate `401` check. Extracted `handleAuthError()` to eliminate string-based 401 detection. Full reactive store rewrite to Svelte 5 `$state` is a lower-priority refactor.

### 10. D-Bus connection panics server on failure
**Status:** Not Fixed
**File:** `backend/steam-game-control/internal/service/game.go:40`
**Impact:** Server crashes if D-Bus is unavailable.
**Description:** `NewGameService()` calls `panic(err)` if it cannot connect to systemd D-Bus.

### 11. No admin enforcement
**Status:** Not Fixed
**File:** `backend/steam-game-control/internal/handler/game.go`
**Impact:** `is_admin` flag is stored in JWT but never enforced. All users have full access.
**Description:** The `is_admin` claim is returned in login response but no middleware checks it.

### 12. Default JWT secret in config
**Status:** Not Fixed
**File:** `backend/steam-game-control/configs/config.yaml`
**Impact:** Weak security - default secret is `"your-secret-key-here"`.
**Description:** The shipped config contains a well-known default JWT signing secret.

## High Issues (Resolved)

### 32. No rate limiting on login endpoint
**Status:** Not Fixed
**File:** `backend/steam-game-control/internal/handler/auth.go`
**Impact:** Enables brute-force attacks against the login endpoint.

## Medium Issues

### 13. String-based 401 detection
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/lib/stores/games.ts`
**Impact:** Brittle auth check - breaks if error message format changes.
**Fix:** Replaced string-based `error.message.includes('401')` with proper `response.status === 401` check in `ApiClient.request()`. The games store now uses `handleAuthError()` which checks `isAuthenticated()` instead of string matching.

### 14. No token expiration/refresh
**Severity:** Medium
**File:** `frontend/games-frontend/src/lib/stores/auth.ts`
**Impact:** Auth never expires automatically; no refresh mechanism.
**Description:** Tokens are stored forever in localStorage with no expiry check or refresh logic.

### 15. `@html` rendering of log messages (potential XSS)
**Severity:** Medium
**File:** `frontend/games-frontend/src/routes/games/[name]/logs/+page.svelte`
**Impact:** Potential XSS if log messages contain malicious HTML/JS.
**Description:** Log messages are rendered with `{@html ...}` without escaping.

### 16. `fetchGames()` called redundantly on game detail pages
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/routes/games/[name]/+page.svelte`
**Impact:** Unnecessary API calls - re-fetches entire game list every time a game detail page is visited.
**Fix:** Now only calls `fetchGames()` if the current game is not already in the cached list.

### 17. Error alert in game control page is not dismissable
**Status:** ✅ Resolved
**File:** `frontend/games-frontend/src/routes/games/[name]/+page.svelte`
**Impact:** UX issue - errors stay on screen forever.
**Fix:** The Alert component already uses `dismissable` prop in the code.

### 18. `ListGames` does O(n) scan per status check
**Status:** Not Fixed
**File:** `backend/steam-game-control/internal/service/game.go`
**Impact:** Inefficient - `GetGameStatus` calls `ListUnits()` and iterates the full list every time.
**Description:** Could be optimized with a map lookup instead of linear scan.

---

## Low / Style Issues

### 19. `icon` package installed but unused
**Severity:** Low
**File:** `frontend/games-frontend/package.json`
**Impact:** Unnecessary dependency in node_modules.
**Description:** `flowbite-svelte-icons` is installed but no icons from it are used in the codebase.

### 20. `app.d.ts` has all interfaces commented out
**Severity:** Low
**File:** `frontend/games-frontend/src/app.d.ts`
**Impact:** Type declarations missing for Error, Locals, PageData, etc.
**Description:** SvelteKit type declarations are all commented out/stubbed.

### 21. Duplicate `401` condition check in `games.ts`
**Severity:** Low
**File:** `frontend/games-frontend/src/lib/stores/games.ts`
**Impact:** Code duplication.
**Description:** `error.message.includes('401')` appears twice in the same condition.

### 22. Dashboard cards always show "Ready" badge regardless of actual status
**Severity:** Low
**File:** `frontend/games-frontend/src/routes/dashboard/+page.svelte`
**Impact:** Misleading UI - status badge doesn't reflect actual game state.

### 23. No loading states on game control buttons
**Severity:** Low
**File:** `frontend/games-frontend/src/routes/games/[name]/+page.svelte`
**Impact:** No visual feedback while start/stop/restart API calls are in progress.

### 24. `extractMessage` regex could be fragile
**Severity:** Low
**File:** `frontend/games-frontend/src/routes/games/[name]/logs/+page.svelte`
**Impact:** Log lines with multiple `]` characters may be parsed incorrectly.
**Description:** `/\] (.+)$/` regex extracts after the last `]`, which may not be correct for all log formats.

### 25. `formatTimestamp` silently fails on non-numeric input
**Severity:** Low
**File:** `frontend/games-frontend/src/routes/games/[name]/logs/+page.svelte`
**Impact:** Invalid timestamps produce `Invalid Date` without warning.
**Description:** `parseInt(timestamp)` silently returns `NaN` on non-numeric input.

### 26. No environment variable for API base URL
**Severity:** Low
**File:** `frontend/games-frontend/src/lib/api/client.ts`
**Impact:** API URL is hardcoded based on `window.location.hostname`, fragile across environments.

### 27. No health check endpoint
**Severity:** Low
**File:** `backend/steam-game-control/`
**Impact:** No `/health` or `/ready` endpoint for load balancers or orchestration.

### 28. No logging middleware
**Severity:** Low
**File:** `backend/steam-game-control/`
**Impact:** No request logging for audit or debugging.

### 29. `go.mod` specifies Go 1.26.1
**Severity:** Low
**File:** `backend/steam-game-control/go.mod`
**Impact:** Go 1.26.1 does not exist yet. Likely a typo for 1.21 or 1.22.

### 30. Documentation-code mismatch
**Severity:** Low
**Files:** `README.md`, `IMPLEMENTATION_SUMMARY.md` vs actual code
**Impact:** Docs mention LevelDB, bcrypt-ready structure, `/auth/register` endpoint - none exist in code.

### 31. Empty `games` array vs nil inconsistency
**Severity:** Low
**File:** `backend/steam-game-control/internal/handler/game.go`
**Impact:** Minor - `ListGames` returns nil on empty result but handler converts to `[]string{}`.

### 32. No rate limiting on login endpoint
**Severity:** High
**File:** `backend/steam-game-control/internal/handler/auth.go`
**Impact:** Enables brute-force attacks against the login endpoint.

---

## Status Legend

| Symbol | Meaning |
|--------|---------|
| 🔴 | Critical - Must fix before production |
| 🟠 | High - Should fix in next sprint |
| 🟡 | Medium - Plan to fix |
| 🟢 | Low / Style - Nice to have |
