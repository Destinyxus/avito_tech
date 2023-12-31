# SEGMENTS SERVICE

## О сервисе 

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)


## Запуск сервиса

Перед запуском не забудьте объявить переменную среды POSTGRES и создать БД с названием на ваш выбор (данное название запросто замените в переменных среды и docker-compose)
````
export POSTGRES="ваш путь до дб в формате postgres://postgres:password@host/dbname?sslmode=disable"
````

Помимо этого, для корректного запуска сервиса необходимо указать флаг -path
````
go run main.go -path "путь до конфиг файла"
````

## Функционал


 ### Создать сегмент

Создание сегмента происходит по эндпоинту /create-segment методом POST

Пример запроса:
````
{   
"slug":"AVITO_VOICE_MESSAGES"
    }
````

### Удалить сегмент 

Удаление сегмента происходит по эндпоинту /delete-segment методом DELETE

Пример запроса:
````
{   
"slug":"AVITO_VOICE_MESSAGES"
    }
````
Стоит отметить, что при удалении, сегмент не исчезает из таблицы, а просто помечается как неактивный.

### Создать пользователя

Создание пользователя происходит по эндпоинту /create-user методом POST

Пример запроса:

````
{   
"name":"Vladimir"
    }
````

### Добавление пользователя в сегмент

Добавление пользователя в сегмент происходит по эндпоинту /addUser методом POST

Пример запроса:

````
{   
"id":"1", // id пользователя
"slug":["AVITO_VOICE_MESSAGES"],
"ttl":2
    }
````
Поле TTL опционально и указывает на срок активации сегмента у пользователя (формат дней)

### Удаление из сегмента/ов пользователя

Удаление происходит по эндпоинту /deleteFromUser методом DELETE

Пример запроса:

````
{   
"id":"1",
"slug":["AVITO_VOICE_MESSAGES"],
    }
````

### Получение активных сегментов

Получение активных сегментов происходит по эндпоинту /activeSegments методом GET

Пример запроса:

````
{   
"id":"1" // id пользователя
    }
````

### Получение отчета в формате CSV

Получение отчета по эндпоинту /getCSV, метод GET

Пример запроса:

````
{   
"user_id":1,
"year":2023,
"month":8
    }
````
На данный момент функционал позволяет получить отчет только за один месяц.

