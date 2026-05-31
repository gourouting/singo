# Singo

Singo: Simple Single Golang Web Service

go-crud has officially been renamed to Singo!

Build web services with Singo: use the simplest architecture to implement a practical framework that can serve massive numbers of users.

https://github.com/Gourouting/singo

## Changelog

1. API testing is now supported.
2. Go 1.25 is now supported. Please install this Go version before using this project.

## Video Tutorial

[Let's Build a G Site! Golang Full-Stack Programming Live Tutorial](https://space.bilibili.com/10/channel/detail?cid=78794)

## Example Projects Built with Singo

Giligili, a Bilibili-like G site: https://github.com/Gourouting/giligili

Token-based mobile login example built with Singo: https://github.com/bydmm/singo-token-exmaple

## Purpose

This project uses a series of popular Golang components. You can use it as a foundation to quickly build RESTful web APIs.

## Features

This project integrates many components required for API development:

1. [Gin](https://github.com/gin-gonic/gin): a lightweight web framework that claims to have the fastest routing in Golang.
2. [GORM](https://gorm.io/index.html): an ORM tool. This project is intended to be used with MySQL.
3. [Gin-Session](https://github.com/gin-contrib/sessions): a session management tool for the Gin framework.
4. [Go-Redis](https://github.com/go-redis/redis): a Golang Redis client.
5. [godotenv](https://github.com/joho/godotenv): an environment variable tool for development environments, making environment variables easier to use.
6. [Gin-Cors](https://github.com/gin-contrib/cors): a CORS middleware for the Gin framework.
7. [httpexpect](https://github.com/gavv/httpexpect): an API testing tool.
8. Basic internationalization (i18n) functionality implemented in this project.
9. This project stores login state with cookie-based sessions. You can switch to token authentication if needed.

This project has already implemented some common code for reference and reuse:

1. A user model.
2. The `/api/v1/user/register` user registration endpoint.
3. The `/api/v1/user/login` user login endpoint.
4. The `/api/v1/user/me` user profile endpoint, which requires a session after login.
5. The `/api/v1/user/logout` user logout endpoint, which requires a session after login.

This project has also pre-created a series of folders for the following modules:

1. The `api` folder is the controller layer of the MVC framework. It coordinates all parts to complete each task.
2. The `model` folder stores database models and database operation code.
3. The `service` folder handles more complex business logic. Modeling business logic can effectively improve business code quality, such as user registration, top-ups, and order placement.
4. The `serializer` folder stores common JSON models and converts database models from `model` into JSON objects required by the API.
5. The `cache` folder contains Redis cache-related code.
6. The `auth` folder contains permission control code.
7. The `util` folder contains common utility helpers.
8. The `conf` folder stores static configuration files. Translation-related configuration files are placed in `locales`.

## Godotenv

The project depends on the following environment variables at startup. You can also create a `.env` file in the project root to set environment variables more conveniently. This is recommended for development environments.

```shell
DB_USER="db_user" # MySQL user
DB_PASSWORD="db_password" # MySQL password
DB_HOST="127.0.0.1" # MySQL host
DB_PORT="3306" # MySQL port
DB_NAME="db_name" # MySQL database name
DB_CHARSET="utf8" # MySQL charset
DB_PARSE_TIME="True" # Parse MySQL time values
DB_LOC="Local" # MySQL time zone
REDIS_ADDR="127.0.0.1:6379" # Redis port and address
REDIS_PW="" # Redis connection password
REDIS_DB="" # Redis database, from 0 to 10
SESSION_SECRET="setOnProducation" # Session secret. It must be set and must not be leaked.
GIN_MODE="debug"
```

## Go Mod

This project uses [Go Mod](https://github.com/golang/go/wiki/Modules) to manage dependencies.

```shell
go mod init singo
export GOPROXY=http://mirrors.aliyun.com/goproxy/
go run main.go // automatically installs dependencies
```

## Run

```shell
go run main.go
```

After the project starts, it runs on port 3000. You can change this; see the Gin documentation for details.

## API Testing

[New] This project includes built-in API tests.

#### Usage

0. Make sure you are in the project root directory.
1. Create a test-specific environment variable file in the `test` directory.

```shell
cp test/.env.example test/.env
```

2. Modify the environment variables in `test/.env` to ensure that MySQL and Redis can be connected normally.
3. Run tests from the project root and enable `-v` to check whether the tests are running correctly.

```shell
go test -v ./test
```

4. After confirming that the tests run correctly, remove the `-v` flag and check whether the tests pass.

```shell
go test ./test
ok      singo/test      (cached)
```
