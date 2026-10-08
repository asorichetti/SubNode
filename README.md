# SubNode

A self-hosted, invite-only community platform inspired by **Reddit** and **Lobsters**.

## 🎯 Vision

SubNode brings the best of Reddit-style threading and Lobsters-style technical curation into a **self-hosted, invite-only** experience. The person hosting a SubNode instance controls who gets in through invite codes—no public signups, no noise, just your community.

## ✨ Key Features

- **Reddit-style**: Sub-nodes (subreddits), posts, comments, upvotes/downvotes
- **Lobsters-inspired**: Tag-based discovery, clean minimal UI, community moderation
- **Invite-only**: Host generates invite codes; only invited users can join
- **Self-hosted**: Runs on your local Kubernetes cluster (or even a single machine)
- **Lightweight**: Minimal resource footprint, easy to maintain

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│                 Kubernetes                   │
│                                              │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │  Frontend│  │  Backend │  │  Postgres  │  │
│  │  (Nginx) │  │  (Go)    │  │  (DB)      │  │
│  └──────────┘  └──────────┘  └───────────┘  │
│                                              │
│  ┌──────────────────────────────────────┐    │
│  │         Ingress Controller            │    │
│  └──────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

## 🚀 Quick Start

```bash
# Deploy to your local cluster
make deploy

# Access at http://localhost
```

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for full setup instructions.

## 📁 Project Structure

```
SubNode/
├── backend/          # Go API server
│   ├── cmd/          # Entry points
│   ├── internal/     # Business logic
│   ├── pkg/          # Shared packages
│   └── migrations/   # Database migrations
├── frontend/         # React + TypeScript UI
│   ├── src/
│   └── public/
├── k8s/              # Kubernetes manifests / Helm chart
│   ├── helm/
│   └── manifests/
├── docs/             # Documentation
└── compose/          # Docker Compose for local dev
```

## 🛠️ Tech Stack

- **Backend**: Go (Gin framework)
- **Frontend**: React + TypeScript + Tailwind CSS
- **Database**: PostgreSQL
- **Cache**: Redis (optional)
- **Deployment**: Kubernetes (Helm)
- **Auth**: JWT + Invite code system

## 🔑 Invite System

The host generates one-time or multi-use invite codes. These codes:
- Link a new user to the hosting instance
- Can expire or be revoked
- Can be single-use or multi-use
- Are the ONLY way to join (no public registration)

## 📜 License

MIT
