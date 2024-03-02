# Coding Exercise

Implement a Restful task list API as well as run this application in container.

- Spec
  - Fields of task:
      - id
          - Type: Int
          - Value: 0 (Incremental)
      - text
          - Type: String
      - status
          - Type: Bool
          - Value
              - 0=Incomplete
              - 1=Complete
  - Response headers
      - Content-Type=application/json
  - Unit Test
  - Manage codebase on Github

- Runtime Environment Requirement
    - If you choose Python
        - Python 3.9+
        - Flask 2.0.x or FastAPI 0.89+
    - If you choose Golang
        - Go 1.17.8+
        - Gin 1.7.7
    - Other language (JAVA, C#) is OK as well
    - Containerize your application (Use Docker)

- About Database
  - DB is a complex component that you can use in-memory mechanism to handle data storage problem.


### 1.  GET /tasks (list tasks)
```
{
    "result": [
        {"id": 1, "text": "task", "status": 0}
    ]
}
```

### 2.  POST /task  (create task)
```
request
{
  "text": "買晚餐"
}

response status code 201
{
    "result": {"text": "買晚餐", "status": 0, "id": 1}
}
```

### 3. PUT /task/<id> (update task)
```
request
{
  "text": "買早餐",
  "status": 1,
  "id": 1
}

response status code 200
{
  "result":{
    "text": "買早餐",
    "status": 1,
    "id": 1
  }
}
```

### 4. DELETE /task/<id> (delete task)
```
Response status code 200
```