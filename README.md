# Gogolook2024 - interview homework project
By the [SPEC](./SPEC.md), provide endpoints to access task resouces. 
## Makefile
Use `make` command to manage project
```
gogolook2024 - Makefile commands

make [options]

release     Build docker container
build       Build local
run         Run gogolook2024 at local
test        Run test cases
start       Start docker compose
redis       Start docker compose only Redis, before you run applicaion at local
clean       Remove object files and cached files
format      Format sources
```
## Environment
### Dev
Development base on
- GO: 1.18
- GIN: v1.9.1
### Runtime
**conf/app.ini** keeps all parameters, list major keys as below:
1. HTTP Port: 8000
2. Redis Port: 6379
  Data save on Redis. the Docker Compose already contains a Redis service. If run at local, use `make redis` first, or start your own redis at port 6379

## Implementation
the major implementation is in `/internal`
| module               | description                                    |
| -------------------- | ---------------------------------------------- |
| /internal/router     | routing table to relative handler              |
| /internal/middleware | Handler of gin router, process request content |
| /internal/service    | Busniess logic, define the behavior            |
| /internal/repo       | Data source interface                          |
| /internal/models     | data type definition                           |
|                      |                                                |

## Test
2 Major logic middleware and service have been covered by test cases

<img width="480" alt="image" src="https://github.com/shenmengkai/gogolook2024/assets/15992122/ccf261f8-393d-4285-a557-48f09fe776ed">

`make test` use *gotestsum* for better output and do automatically installation, or you prefer go test just by
```
go test ./internal/...
```



## Lack of features or issuses could be improve
1. Authenication
2. Race condition on data access
3. swagger/api doc auto generation
