# todo-server

## Примеры запросов

Добавление задачи "Buy milk" с приоритетом 2:

```sh
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"name": "Buy milk", "priority": 2}' \
    localhost:9000/todo
```

Получение задачи с идентификатором 1:

```sh
curl -X GET localhost:9000/todo/1
```

Получение всего списка задач:

```sh
curl -X GET localhost:9000/todo
```

Получение списка выполненных задач с приоритетом 2:

```sh
curl -X GET 'localhost:9000/todo?done=true&priority=2'
```

Установление/снятие отметки о выполнении задачи с идентификатором 1:

```sh
curl -X PUT localhost:9000/todo/1
```

Удаление задачи с идентификатором 1:

```sh
curl -X DELETE localhost:9000/todo/1
```
