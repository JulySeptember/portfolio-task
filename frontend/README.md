# Frontend

Serverless Task Management App frontend built with
Next.js App Router.

---

# Tech Stack

| Layer | Technology |
| --- | --- |
| Framework | Next.js |
| Language | TypeScript |
| Styling | Tailwind CSS v4 |
| UI Components | shadcn/ui + Radix UI |
| API State Management | React Query |
| Global State Management | Zustand |
| Validation | Zod |
| Icons | Lucide React |

---

# Development

Install dependencies:

```bash
npm install
```

Start development server:

```bash
npm run dev
```

Application:

```text
http://localhost:3000
```

---

# Environment Variables

Create `.env.local`.

Example:

```env
NEXT_PUBLIC_API_URL=https://xxxxxxxx.execute-api.ap-northeast-1.amazonaws.com

NEXT_PUBLIC_COGNITO_DOMAIN=xxxxxxxx.auth.ap-northeast-1.amazoncognito.com

NEXT_PUBLIC_COGNITO_CLIENT_ID=xxxxxxxx

NEXT_PUBLIC_COGNITO_REDIRECT_URI=http://localhost:3000/auth/callback

NEXT_PUBLIC_COGNITO_LOGOUT_URI=http://localhost:3000/login
```

---

# Directory Structure

```text
src/
├── app
├── components
├── features
├── lib
├── providers
└── store
```

---

# App Structure

```text
app/
├── (auth)
├── (protected)
├── api
└── auth
```

| Directory | Role |
| --- | --- |
| (auth) | login pages |
| (protected) | authenticated pages |
| api | route handlers |
| auth | Cognito callback handling |

---

# UI

This project uses:

```text
- Tailwind CSS
- shadcn/ui
- Radix UI
```

Included base UI components:

```text
- badge
- button
- card
- dialog
- dropdown-menu
- input
- label
- skeleton
- table
- textarea
- toaster
```

---

# State Management

## React Query

Used for API state management.

Responsibilities:

```text
- API fetch
- cache management
- loading state
- mutation
- optimistic updates
- refetch
```

---

## Zustand

Used for lightweight frontend global state.

Responsibilities:

```text
- authentication state
```

---

# Authentication

Authentication uses AWS Cognito Hosted UI.

Frontend authentication flow:

```text
Login button
  ↓
Cognito Hosted UI
  ↓
redirect (/auth/callback)
  ↓
token exchange
  ↓
localStorage storage
  ↓
Authorization: Bearer <id_token>
```

Current implementation stores tokens in localStorage
for development simplicity.

Production-grade deployments should consider
HttpOnly cookie based authentication strategies.

---

# JWT Validation

JWT signature validation is delegated to
API Gateway JWT Authorizer.

```text
Client
  ↓
API Gateway JWT Authorizer
  ↓
validated claims
  ↓
Lambda
```

---

# User Bootstrap

After login, frontend calls bootstrap API
to synchronize Cognito users with users table.

Behavior:

```text
- INSERT on first login
- UPDATE on existing users
- Cognito sub used as auth_user_id
```

---

# API

Backend API endpoints:

```text
/api/v1/*
```

Health check endpoint:

```text
/health
```

Swagger/OpenAPI:

```text
/api/docs
```

---

# Features

## Tasks

Implemented / planned capabilities:

```text
- create task
- update task
- update task status
- delete task
- pagination
- filtering
- sorting
```

---

# Frontend Architecture

This project adopts feature-based architecture.

```text
features/
├── auth
└── tasks
```

Each feature internally manages:

```text
- api
- hooks
- components
- schemas
- types
- utils
```

Benefits:

```text
- scalability
- separation of concerns
- maintainability
- feature isolation
```

---

# Deployment

Frontend deployment architecture:

```text
CloudFront
  ↓
Next.js hosting environment
```

Static assets may additionally be distributed via:

```text
CloudFront
  ↓
S3
```

---

# Notes

This project focuses on:

```text
- modern frontend architecture
- scalable feature organization
- serverless integration
- authentication flow design
- API-driven UI development
```