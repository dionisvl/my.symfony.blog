# web3main.pro — Astro Frontend

SSR blog frontend built with Astro 5, TailwindCSS v4, and `@astrojs/node` adapter (standalone mode).
Proxies API requests to the Go backend (`api-go`) inside the Docker network.

---

## Stack

- **Astro 5** — SSR, `@astrojs/node` standalone adapter
- **TailwindCSS v4** — via `@tailwindcss/vite`
- **Shiki** — syntax highlighting (`min-light` / `github-dark-default`)
- **pnpm** — package manager
- **Node 24**

---

## Running via Docker (recommended)

All services are orchestrated from the repo root.

```bash
# Dev (HMR, no bundling — many files in DevTools is normal)
make up

# Stop
make down
```

Dev URL: `http://<ASTRO_HOST>` (set via `ASTRO_HOST` in root `.env`, default: `phpqas.local`)

### Test production build locally

In dev mode Vite serves ~50 individual unbundled modules — this is expected behaviour, not a problem. Production build outputs a handful of minified chunks.

To run a prod build locally against your dev environment (from repo root):

```bash
make astro-build-local
```

Site will be available at `http://$ASTRO_HOST` as usual. To go back to dev mode — `make up`.

---

## Running locally (without Docker)

```bash
cp .env.example .env
# fill in API_URL, API_KEY, PUBLIC_SITE_URL, etc.

pnpm install
pnpm dev
```

---

## Environment variables

Configured via `.env` (dev) or `.env.production` (prod). See `.env.example` for the full list.

| Variable | Description |
|----------|-------------|
| `API_URL` | Go API URL, e.g. `http://api-go:8081` |
| `API_KEY` | API key, must match `GO_API_KEY` in root `.env` |
| `PUBLIC_SITE_URL` | Canonical site URL, used in OG/meta tags |
| `PUBLIC_SITE_TITLE` | Site title |
| `PUBLIC_SITE_AUTHOR` | Site author |
| `PUBLIC_SITE_DESC` | Meta description |
| `PUBLIC_GOOGLE_SITE_VERIFICATION` | Google Search Console meta tag (optional) |

Server-side variables (`API_URL`, `API_KEY`) are never exposed to the client.

---

## Scripts

```bash
pnpm dev          # Start dev server
pnpm build        # Type-check + build for production
pnpm preview      # Preview production build locally
pnpm lint         # ESLint
pnpm format       # Prettier (write)
pnpm format:check # Prettier (dry-run)
pnpm sync         # Regenerate Astro type declarations (.astro/)
```

---

## Project structure

```
frontend-astro/
├── src/
│   ├── config.ts              # SITE constants (reads from env)
│   ├── lib/
│   │   ├── api.ts             # API client (apiFetch + all endpoints)
│   │   └── types.ts           # TypeScript types for all API responses
│   ├── pages/
│   │   ├── index.astro        # / — Hero + Featured + Recent posts
│   │   ├── 404.astro
│   │   ├── search.astro       # /search?q=
│   │   ├── contacts.astro     # /contacts
│   │   ├── posts/
│   │   │   ├── index.astro               # /posts/ (page 1)
│   │   │   ├── page/[page].astro         # /posts/page/N/
│   │   │   └── [...slug]/index.astro     # /posts/<slug>/
│   │   ├── tags/
│   │   │   └── [tag]/[...page].astro     # /tags/<slug>/
│   │   ├── categories/
│   │   │   └── [slug]/[...page].astro    # /categories/<slug>/
│   │   └── api/
│   │       ├── comment.ts               # POST /api/comment → api-go
│   │       ├── subscribe.ts             # POST /api/subscribe → api-go
│   │       ├── send-contact.ts          # POST /api/send-contact → api-go
│   │       └── postlike/[postId].ts     # POST /api/postlike/:id → api-go
│   ├── components/
│   │   ├── Header.astro, Footer.astro, MobileMenu.astro
│   │   ├── Card.astro         # Post card
│   │   ├── Pagination.astro
│   │   ├── Breadcrumb.astro
│   │   ├── Tag.astro
│   │   ├── Datetime.astro     # Client-side date formatting (user's locale/timezone)
│   │   └── ...
│   └── assets/
│       └── fonts/             # Local fonts (woff2)
├── public/                    # Static assets
├── astro.config.ts
├── .env                       # Dev secrets (gitignored)
├── .env.production            # Prod secrets (gitignored)
└── .env.example               # Template
```

---

## Traefik routing (inside Docker)

| Host | Service | Notes |
|------|---------|-------|
| `$ASTRO_HOST` | `astro-frontend` | catch-all, priority 10 |
| `$ASTRO_HOST` + `/admin/*`, `/login`, `/logout` | `symfony` | priority 50 |
| `$API_GO_HOST` | `api-go` | dedicated subdomain, all paths |

Hosts are configured via `ASTRO_HOST` and `API_GO_HOST` in the root `.env`.

Astro proxies API calls server-side to `api-go` via `API_URL` (see `frontend-astro/.env`). Write endpoints (`/api/comment`, `/api/subscribe`, `/api/send-contact`, `/api/postlike/*`) are handled by Astro API routes — the browser never calls `api-go` directly and never sees `API_KEY`.
