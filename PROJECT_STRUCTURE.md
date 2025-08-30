# GinFileHub 项目结构

```
.
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── auth/
│   │   ├── auth.go
│   │   └── apikey/
│   │       └── apikey.go
│   ├── files/
│   │   ├── handler.go
│   │   └── service.go
│   ├── routes/
│   │   └── routes.go
│   └── logger/
│       └── logger.go
├── pkg/
│   └── utils/
│       └── utils.go
├── web/
│   └── dist/
├── migrations/
│   └── 001_init.up.sql
├── configs/
│   └── .ginfilehub.yml
├── test/
│   ├── handlers/
│   │   ├── files_test.go
│   │   └── auth_test.go
│   └── data/
├── Dockerfile
└── docs/
    └── openapi.yaml