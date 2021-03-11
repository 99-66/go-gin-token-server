# Golang Gin JWT Token Server

JWT 토큰을 발급, 검증, 재발급하는 API

* 재발급을 하는 경우에 새로운 AccessToken, RefreshToken을 생성하여 반환하고, 
* 재발급에 요청한 RefreshToken은 재사용을 막기 위해 Redis에 넣어 이전 토큰 만료시간까지 재사용을 하지 못하게 한다

### 환경변수 설정

```shell
export GIN_MODE=release
export PORT=8000

# JWT
export JWT_ACCESS_SECRET=test
export JWT_REFRESH_SECRET=test

# Redis Configuration
export REDIS_HOST=192.168.0.1
export REDIS_PORT=6379
export REDIS_DB=0
export REDIS_PASSWORD=password
```

### Create Token

ID, Domain, Roles를 Body로 전송하여 토큰을 생성한다
```go
[Request]

POST /api/token/create
{
    "id":1,
    "domain":"example.com",
    "roles": ["admin", "read-only"]
}
```
```go
[Response]

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJkb21haW4iOiJleGFtcGxlLmNvbSIsImV4cCI6MTYxNTQ1MTEwNCwiaWF0IjoxNjE1NDUwODA0LCJyb2xlcyI6WyJhZG1pbiIsInJlYWQtb25seSJdLCJ1c2VyX2lkIjoxfQ.5-qeIBJ5tMRV6iOWZ-ZdATFbaYBQ-EPbncSVZkKpKWM",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb21haW4iOiJleGFtcGxlLmNvbSIsImV4cCI6MTYxNTg4MjgwNCwiaWF0IjoxNjE1NDUwODA0LCJyb2xlcyI6WyJhZG1pbiIsInJlYWQtb25seSJdLCJ1c2VyX2lkIjoxfQ.riB-cfBY05cV5e18YBZ7b7UJOtZU04bLdTCeygElwTY",
    "access_expire": 1615451104,
    "refresh_expire": 1615882804,
    "iat": 1615450804,
    "id": 1,
    "domain": "example.com",
    "roles": ["admin", "read-only"]
}
```

### Verify Token

Authorization 헤더에 토큰값을 넣어서 전송한다
```go
[Request]

GET /api/token/verify
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJkb21haW4iOiJnb29nbGUuY29tIiwiZXhwIjoxNjE1NDQ3MjAzLCJpYXQiOjE2MTU0NDY5MDMsInJvbGVzIjpbImFkbWluIiwicmVhZC1vbmx5Il0sInVzZXJfaWQiOjF9.pbl_bSTHzMMcGjI7G59N6JAEpe-QK-_nk03KUqe4N5o
```
```go
[Response]

{
    "authorized": true,
    "domain": "example.com",
    "exp": 1615452000,
    "iat": 1615451700,
    "id": 1,
    "roles": ["admin", "read-only"]
}
```


### Refresh Token
Request Body에 refreshToken을 넣어서 전송한다
```go
[Request]

POST /api/token/refresh
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb21haW4iOiJnb29nbGUuY29tIiwiZXhwIjoxNjE1ODgyNDg5LCJpYXQiOjE2MTU0NTA0ODksInJvbGVzIjpbImFkbWluIiwicmVhZC1vbmx5Il0sInVzZXJfaWQiOjF9.iiVGA6YpG4qlJWDstj_adq-Rvhvb7wqHvQjJKbvb_u0"
}
```

```go
[Response]

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJkb21haW4iOiJleGFtcGxlLmNvbSIsImV4cCI6MTYxNTQ1MTM5MSwiaWF0IjoxNjE1NDUxMDkxLCJyb2xlcyI6WyJhZG1pbiIsInJlYWQtb25seSJdLCJ1c2VyX2lkIjoxfQ.xrYm1PQqguzqeMl7Y5BJNFTF5--E6yT5EhNYrWFuH3Q",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb21haW4iOiJleGFtcGxlLmNvbSIsImV4cCI6MTYxNTg4MzA5MSwiaWF0IjoxNjE1NDUxMDkxLCJyb2xlcyI6WyJhZG1pbiIsInJlYWQtb25seSJdLCJ1c2VyX2lkIjoxfQ.vrkpk1Mqf5JRBkZCRtTVN1R14vq_AvTy90ZaddlD03M",
    "access_expire": 1615451391,
    "refresh_expire": 1615883091,
    "iat": 1615451091,
    "id": 1,
    "domain": "example.com",
    "roles": ["admin", "read-only"]
}
```
---