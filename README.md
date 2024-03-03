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
create      curl test to create task by picking random word
update      curl test to update task by example 'make update id=10 text=movie status=1'
delete      curl test to delete task by example 'make delete id=10'
list        curl test to list tasks
```

## Environment
### Dev
Development base on
- GO: 1.18
- GIN: v1.7.7
### Runtime
**conf/app.ini** keeps all parameters, list major keys as below:
1. HTTP Port: 8000
2. Redis Port: 6379
  Data save on Redis. the Docker Compose already contains a Redis service. If run at local, use `make redis` first, or start your own redis at port 6379

## Implementation
<img width="517" alt="image" src="https://github.com/shenmengkai/gogolook2024/assets/15992122/fdde246c-34f5-4289-8edc-71e56be030c5">


the major implementation is in `/internal`
| module                                        | description                                    |
| :-------------------------------------------- | :--------------------------------------------- |
| [/internal/router](./internal/router )        | routing table to relative handler              |
| [/internal/middleware](./internal/middleware) | Handler of gin router, process request content |
| [/internal/service](./internal/service)       | Busniess logic, define the behavior            |
| [/internal/repo](./internal/repo)             | Data source interface                          |
| [/internal/models](./internal/models)         | data type definition                           |

## Unit Test
2 Major logic middleware and service have been covered by test cases

<img width="603" alt="image" src="https://github.com/shenmengkai/gogolook2024/assets/15992122/496ff5c8-673f-4bd7-86e5-f2b2ed2f09e7">


`make test` use *gotestsum* for better output and do automatically installation, or you prefer go test just by
```
go test ./internal/...
```
## Integration Test
make provide commends to manually request on endpoints by curl

- List tasks
```
make list
```

- Create task with random word by system
```
make create
```

- Update task with specific _text_ or _status_
```
make update id=17 text=swim status=1
```

- Delete task by id
```
make delete id=17
```

![gogolook2024](https://github.com/shenmengkai/gogolook2024/assets/15992122/2156613b-e548-40ba-8098-809fc280cfe7)

## Lack of features or issues could be improved
1. Authenication or API_KEY
2. Database
  Currently use redis for quick implementation, and lack of consistansy, chance to hit race condition, change to database to achieve atomic access 
3. swagger/api doc auto generation
