# GO demo
By the [SPEC](./SPEC.md), provide endpoints to access task resouces. 
## Makefile
Use `make` command to manage project
```
go-demo - Makefile commands

make [options]

release     Build docker container
build       Build local
run         Run go-demo at local
test        Run test cases
start       Start docker compose
redis       Start docker compose only Redis, before you run applicaion at local
clean       Remove object files and cached files
format      Format sources
docs        Gernerate /docs
create      curl test to create task by picking random word
update      curl test to update task by example 'make update id=10 text=movie status=1'
delete      curl test to delete task by example 'make delete id=10'
list        curl test to list tasks
```

## Environment
Development base on
- GO: 1.18
- GIN: v1.7.7

### Redis
Data will save on Redis. There are 3 ways to coordinate with app:
1. Docker Compose contains a Redis service, just run `make start` to start app and Redis both.
2. Use `make redis` to start Redis image only.
3. Start your own Redis at port 6379

## Development
**conf/app.ini** keeps app parameters, list major keys as below:
1. HTTP Port: 8000
2. Redis Port: 6379

## Implementation
<img width="517" alt="image" src="https://github.com/shenmengkai/go-demo/assets/15992122/fdde246c-34f5-4289-8edc-71e56be030c5">

the major implementation is in `/internal`
| module                                                           | description                                    |
| :--------------------------------------------------------------- | :--------------------------------------------- |
| [/internal/router](./internal/router/router.go)                  | routing table to relative handler              |
| [/internal/middleware](./internal/middleware/task_middleware.go) | Handler of gin router, process request content |
| [/internal/service](./internal/service/task_service.go)          | Busniess logic, define the behavior            |
| [/internal/repo](./internal/repo/task_repo.go)                   | Data source interface                          |
| [/internal/models](./internal/models/task_model.go)              | data type definition                           |

## Swagger Documentation
Open http://localhost:8000/swagger/index.html when running application
<img width="797" alt="截圖 2024-03-03 下午3 56 04" src="https://github.com/shenmengkai/go-demo/assets/15992122/cf912768-ceaa-40d3-bb10-3db25130b50a">

## Unit Test
2 major logic, middleware and service have been covered by test cases

<img width="603" alt="image" src="https://github.com/shenmengkai/go-demo/assets/15992122/496ff5c8-673f-4bd7-86e5-f2b2ed2f09e7">


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

![go-demo](https://github.com/shenmengkai/go-demo/assets/15992122/2156613b-e548-40ba-8098-809fc280cfe7)

## Lack of features or issues could be improved
1. Authenication or API_KEY
2. Database
  Currently use redis for quick implementation, and lack of consistansy, chance to hit race condition, change to database to achieve atomic access 
