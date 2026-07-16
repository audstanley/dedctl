# games-frontend

SvelteKit web dashboard for managing Steam dedicated game servers.

## Tech Stack

- **Svelte 5** — UI framework with runes (`$state`, `$derived`, `$effect`)
- **SvelteKit** — Full-stack framework with file-based routing
- **Vite** — Build tool and dev server
- **Tailwind CSS v4** — Utility-first CSS framework
- **Flowbite-Svelte** — Pre-built Svelte component library
- **TypeScript** — Type safety
- **ApexCharts** — Charts for dashboard analytics

## Project Structure

```
frontend/games-frontend/
├── src/
│   ├── lib/
│   │   ├── api/
│   │   │   └── client.ts     # API client (fetch wrapper, auth handling)
│   │   ├── stores/
│   │   │   ├── auth.ts       # Auth state (token, user, login/logout)
│   │   │   └── games.ts     # Game store (fetch, status, control operations)
│   │   ├── assets/
│   │   │   └── favicon.svg
│   │   └── index.ts
│   └── routes/
│       ├── +layout.svelte    # Root layout
│       ├── +page.svelte      # Login page
│       ├── dashboard/
│       │   ├── +layout.svelte
│       │   └── +page.svelte  # Game servers grid
│       ├── games/
│       │   ├── [name]/
│       │   │   ├── +page.svelte      # Game server control panel
│       │   │   └── logs/
│       │   │       └── +page.svelte  # Real-time log viewer (SSE)
│       │   └── +layout.svelte
│       └── admin/
│           ├── +layout.svelte
│           └── settings/
│               └── +page.svelte  # Admin settings page
├── static/
│   └── robots.txt
├── .env.example              # Environment variable template
├── svelte.config.js          # SvelteKit configuration
├── vite.config.ts            # Vite configuration
├── tailwind.config.js        # Tailwind CSS configuration
├── postcss.config.js         # PostCSS configuration
└── package.json
```

## Pages

| Route | Description |
|-------|-------------|
| `/` | Login page with username/password form |
| `/dashboard` | Game servers grid showing all available servers |
| `/games/[name]` | Individual game server control panel (start/stop/restart) |
| `/games/[name]/logs` | Real-time log viewer via Server-Sent Events |
| `/admin/settings` | Admin settings for global metadata and cover art |

## State Management

- **`$lib/stores/auth.ts`** — Manages authentication state (JWT token, user info) with `localStorage` persistence. Exposes `login()`, `logout()`, `isAuthenticated()`, and a `subscribe` function for reactive updates.
- **`$lib/stores/games.ts`** — `GamesStore` class wrapping Svelte's `writable` stores for game list and status. Provides methods for fetching games, updating status, and controlling servers (start/stop/restart).

## API Client

**`$lib/api/client.ts`** — Typed API client wrapping `fetch`:

- Auto-includes JWT Bearer token in requests
- Handles 401 responses (redirects to login)
- Exports types: `GameInfo`, `ServerInfo`, `AuthResponse`

## Environment Variables

Copy `.env.example` to `.env` and set:

```bash
VITE_API_BASE_URL=http://localhost:8080
```

## Development

```bash
npm install
npm run dev
```

Dev server runs on port **5174** (configured in `package.json`).

## Build

```bash
npm run build
```

## Preview

```bash
npm run preview
```

## Code Quality

```bash
# Lint check
npm run lint

# Format files
npm run format

# Type check
npm run check
```
