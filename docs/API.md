# SubNode API Documentation

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Endpoints

### Auth

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string",
  "code": "string"
}
```

Response:
```json
{
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "admin|moderator|member",
    "created_at": "string"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}
```

Response:
```json
{
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "admin|moderator|member",
    "created_at": "string"
  }
}
```

#### Me
```http
GET /api/auth/me
Authorization: Bearer <token>
```

Response:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "role": "admin|moderator|member",
  "created_at": "string"
}
```

### Nodes

#### List Nodes
```http
GET /api/nodes
```

Response:
```json
[
  {
    "id": "string",
    "name": "string",
    "slug": "string",
    "description": "string",
    "rules": "string",
    "host_id": "string",
    "created_at": "string",
    "updated_at": "string"
  }
]
```

#### Get Node
```http
GET /api/nodes/:slug
```

Response:
```json
{
  "id": "string",
  "name": "string",
  "slug": "string",
  "description": "string",
  "rules": "string",
  "host_id": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

#### Create Node (Admin)
```http
POST /api/nodes
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "string",
  "description": "string",
  "rules": "string"
}
```

### Posts

#### List Posts
```http
GET /api/nodes/:slug/posts?page=1&limit=20
```

Response:
```json
[
  {
    "id": "string",
    "node_id": "string",
    "author_id": "string",
    "title": "string",
    "body": "string",
    "url": "string",
    "vote_count": 0,
    "author": {
      "id": "string",
      "username": "string",
      "email": "string",
      "role": "string"
    },
    "created_at": "string",
    "updated_at": "string"
  }
]
```

#### Get Post
```http
GET /api/posts/:id
```

#### Create Post
```http
POST /api/nodes/:slug/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "string",
  "body": "string",
  "url": "string",
  "tags": ["string"]
}
```

#### Vote Post
```http
POST /api/posts/:id/vote
Authorization: Bearer <token>
Content-Type: application/json

{
  "vote": "up|down"
}
```

### Invites (Admin)

#### Create Invite
```http
POST /api/invites
Authorization: Bearer <token>
Content-Type: application/json

{
  "max_uses": 1,
  "expires": "2024-12-31T23:59:59Z",
  "role": "member|moderator"
}
```

#### List Invites
```http
GET /api/invites
Authorization: Bearer <token>
```

#### Revoke Invite
```http
DELETE /api/invites/:id
Authorization: Bearer <token>
```

### Health

#### Health Check
```http
GET /health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

#### Readiness Check
```http
GET /ready
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

## Error Responses

All errors follow this format:
```json
{
  "error": "Error message"
}
```

Common status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error
