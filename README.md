Mai rooms - это веб-приложение и telegram-бот, созданные для упрощения просмотра расписания по аудитории и бронирования свободных аудитории Московского авиационного института.

## Документация по API

### Получить все корпусы МАИ

#### Запрос

```http
GET /api/buildings
```

#### Ответ

```json
[
    {
        "id": 18,
        "name": "1"
    },
    {
        "id": 1,
        "name": "11"
    },
    {
        "id": 17,
        "name": "2"
    },
    {
        "id": 2,
        "name": "24"
    },
    {
        "id": 9,
        "name": "24Б"
    },
    {
        "id": 4,
        "name": "3"
    },
    {
        "id": 14,
        "name": "4"
    },
    {
        "id": 15,
        "name": "5"
    },
    {
        "id": 3,
        "name": "7"
    },
    {
        "id": 5,
        "name": "9"
    },
    {
        "id": 10,
        "name": "Берн. 14"
    },
    {
        "id": 6,
        "name": "ГАК"
    },
    {
        "id": 7,
        "name": "ГУК А"
    },
    {
        "id": 8,
        "name": "ГУК Б"
    },
    {
        "id": 13,
        "name": "ГУК В"
    },
    {
        "id": 16,
        "name": "Орш. А"
    },
    {
        "id": 12,
        "name": "Орш. Б"
    },
    {
        "id": 11,
        "name": "Орш. В"
    }
]
```

### Получить все аудитории МАИ для заданного корпуса

#### Запрос

```http
GET /api/rooms/{buildingId}
```

#### Пример

```http
GET /api/rooms/11
```

#### Ответ

```json
[
    {
        "id": 290,
        "name": "Орш. В-101"
    },
    {
        "id": 689,
        "name": "Орш. В-104"
    },
    {
        "id": 291,
        "name": "Орш. В-204"
    },
    {
        "id": 292,
        "name": "Орш. В-205"
    },
    {
        "id": 410,
        "name": "Орш. В-207"
    },
    {
        "id": 601,
        "name": "Орш. В-209"
    },
    {
        "id": 293,
        "name": "Орш. В-210"
    },
    {
        "id": 619,
        "name": "Орш. В-212"
    },
    {
        "id": 294,
        "name": "Орш. В-226"
    },
    {
        "id": 295,
        "name": "Орш. В-227"
    },
    {
        "id": 334,
        "name": "Орш. В-301"
    },
    {
        "id": 503,
        "name": "Орш. В-303"
    },
    {
        "id": 46,
        "name": "Орш. В-304"
    },
    {
        "id": 502,
        "name": "Орш. В-305"
    },
    {
        "id": 259,
        "name": "Орш. В-306"
    },
    {
        "id": 258,
        "name": "Орш. В-309"
    },
    {
        "id": 154,
        "name": "Орш. В-311"
    },
    {
        "id": 156,
        "name": "Орш. В-312"
    },
    {
        "id": 223,
        "name": "Орш. В-314"
    },
    {
        "id": 221,
        "name": "Орш. В-319"
    },
    {
        "id": 260,
        "name": "Орш. В-326"
    },
    {
        "id": 155,
        "name": "Орш. В-402"
    },
    {
        "id": 566,
        "name": "Орш. В-403"
    },
    {
        "id": 531,
        "name": "Орш. В-404"
    },
    {
        "id": 157,
        "name": "Орш. В-405"
    },
    {
        "id": 232,
        "name": "Орш. В-408"
    },
    {
        "id": 353,
        "name": "Орш. В-410"
    },
    {
        "id": 179,
        "name": "Орш. В-411"
    },
    {
        "id": 54,
        "name": "Орш. В-413"
    },
    {
        "id": 55,
        "name": "Орш. В-414"
    },
    {
        "id": 530,
        "name": "Орш. В-423"
    },
    {
        "id": 632,
        "name": "Орш. В-424"
    },
    {
        "id": 633,
        "name": "Орш. В-426"
    },
    {
        "id": 296,
        "name": "Орш. В-501"
    },
    {
        "id": 354,
        "name": "Орш. В-502"
    },
    {
        "id": 45,
        "name": "Орш. В-504"
    },
    {
        "id": 184,
        "name": "Орш. В-505"
    },
    {
        "id": 289,
        "name": "Орш. В-506"
    },
    {
        "id": 335,
        "name": "Орш. В-507"
    },
    {
        "id": 621,
        "name": "Орш. В-508"
    },
    {
        "id": 180,
        "name": "Орш. В-510"
    },
    {
        "id": 661,
        "name": "Орш. В-511"
    },
    {
        "id": 352,
        "name": "Орш. В-512"
    },
    {
        "id": 600,
        "name": "Орш. В-514"
    },
    {
        "id": 386,
        "name": "Орш. В-527"
    },
    {
        "id": 48,
        "name": "Орш. В-536"
    },
    {
        "id": 501,
        "name": "Орш. В-601"
    },
    {
        "id": 224,
        "name": "Орш. В-603"
    },
    {
        "id": 222,
        "name": "Орш. В-605"
    },
    {
        "id": 411,
        "name": "Орш. В-606"
    },
    {
        "id": 669,
        "name": "Орш. В-608"
    },
    {
        "id": 227,
        "name": "Орш. В-610"
    },
    {
        "id": 229,
        "name": "Орш. В-611"
    },
    {
        "id": 670,
        "name": "Орш. В-619"
    },
    {
        "id": 433,
        "name": "Орш. В-621"
    }
]
```

### Получить расписание занятий на аудиторию по дате

#### Запрос

```http
GET /api/schedule
```

**Параметры:**
- `room` (обязательно): id аудитории.
- `date` (обязательно): дата, на которую необхидмо получить расписание.

#### Пример

```http
GET /api/schedule?room=54&date=13.04.2024
```

#### Ответ

```json
[
    {
        "lector": "КУЗНЕЦОВ ПАВЕЛ МИХАЙЛОВИЧ",
        "time_start": "09:00",
        "time_end": "10:30",
        "subject": "Математическое моделирование",
        "groups": "М3З-403Бк-20",
        "type": "ЛК"
    },
    {
        "lector": "КУЗНЕЦОВ ПАВЕЛ МИХАЙЛОВИЧ",
        "time_start": "10:45",
        "time_end": "12:15",
        "subject": "Математическое моделирование",
        "groups": "М3З-403Бк-20",
        "type": "ЛК"
    },
    {
        "lector": "ХОРОШКО АЛЕКСЕЙ ЛЕОНИДОВИЧ / ВИКУЛИН МАКСИМ АЛЕКСАНДРОВИЧ",
        "time_start": "13:00",
        "time_end": "14:30",
        "subject": "WEB технологии",
        "groups": "М3О-232Б-22",
        "type": "ЛР"
    },
    {
        "lector": "ВИКУЛИН МАКСИМ АЛЕКСАНДРОВИЧ / ХОРОШКО АЛЕКСЕЙ ЛЕОНИДОВИЧ",
        "time_start": "14:45",
        "time_end": "16:15",
        "subject": "WEB технологии",
        "groups": "М3О-232Б-22",
        "type": "ЛР"
    }
]
```

### Получить забронированные занятий на аудиторию по дате

#### Запрос

```http
GET /api/schedule
```

**Параметры:**
- `room` (обязательно): id аудитории.
- `date` (обязательно): дата, на которую необходимо получить расписание.

#### Пример

```http
GET /api/reserved_lesssons?date=13.04.2024&room=54
```

#### Ответ

```json
[
    {
        "reserver": "Малин Александр Сергеевич",
        "reserver_id": 1,
        "room_name": "Орш. В-413",
        "room_id": 54,
        "date": "13.04.2024",
        "time_start": "16:30",
        "time_end": "18:00",
        "comment": "Проведение мероприятия"
    }
]
```

### Авторизация

#### Запрос

```http
POST /api/auth/login
```

**Пример тела запроса (JSON):**
```json
{
    "username":"asmalin",
    "password":"12345"
}
```

#### Ответ

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTk4NzM5MTgsImlhdCI6MTcxOTg3MjExOCwidXNlcl9pZCI6MX0.PGaSX8OmCxHwusikIupa5RKoKRW9In4IIKFcrERkkf4",
    "user": {
        "id": 1,
        "username": "asmalin",
        "password": "$2a$10$MpG/dODBAjl98yOQ0HmmquQiJPD.kqK1wurJZV/HLcMRDRvaIemF6",
        "fullname": "Малин Александр Сергеевич",
        "role": "admin",
        "email": "asmalin@mai.ru",
        "tgUsername": ""
    }
}
```

### Получение всех забронированных занятий

#### Запрос

```http
GET /api/all_reserved_lesssons
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации. 

#### Ответ

```json
[
    {
        "reserver": "Малин Александр Сергеевич",
        "reserver_id": 1,
        "room_name": "Орш. В-413",
        "room_id": 54,
        "date": "13.04.2024",
        "time_start": "16:30",
        "time_end": "18:00",
        "comment": "Проведение мероприятия"
    }
]
```

### Забронировать аудиторию

#### Запрос

```http
POST /api/reserve
```

**Пример тела запроса (JSON):**
```json
{
    "roomId":54,
    "date":"13.04.2024",
    "startTime":"16:30",
    "endTime":"18:00",
    "comment":"Проведение мероприятия"
}
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации. 

### Отменить бронь

#### Запрос

```http
POST /api/cancelReservation
```

**Пример тела запроса (JSON):**
```json
{
    "roomId":54,
    "date":"13.04.2024",
    "startTime":"16:30"
}
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации.

### Получить всех пользователей

#### Запрос

```http
GET /api/users
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации.

#### Ответ

```json
[
    {
        "id": 5,
        "username": "iiivanov",
        "fullname": "Иванов Иван Иванович",
        "role": "user",
        "email": "iiivanov@mai.ru"
    },
    {
        "id": 1,
        "username": "asmalin",
        "fullname": "Малин Александр Сергеевич",
        "role": "admin",
        "email": "asmalin@mai.ru"
    }
]
```

### Удалить пользователя

#### Запрос

```http
DELETE /users/delete/{userId}
```

#### Пример

```http
DELETE /users/delete/5
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации.

#### Ответ

```json
{
    "id": 5
}
```

### Изменить данные пользователя

#### Запрос

```http
PATCH /users/delete/{userId}
```

#### Пример

```http
PATCH /api/users/update/6
```

**Пример тела запроса (JSON):**
```json
{
    "email":"mmmax123@mai.ru"
}
```

Этот запрос доступен только для авторизированных пользователей.
Для прохождения этапа аутентификации необходимо отправлять свой Bearer Token, полученный при авторизации.

#### Ответ

```json
{
    "id": 6
}
```

