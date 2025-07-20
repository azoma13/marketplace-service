# marketplace-service

## Описание проекта:
Сервис, который реализует условный маркетплейс. 

## Стек:
    Язык сервиса: Go;
    База данных: PostgreSQL;
    Деплоя зависимостей и самого сервиса: Docker, Docker-Compose.
`Примечание: есть работающее приложение на публичном хостинге по ip:45.12.228.205. Ниже все примеры указаны с его использованием.`

## Сервис предоставляет следующие конечные точки API:
- регистрация пользователя;
- авторизация пользователя; 
- размещение нового объявления (только для авторизованных пользователей). При размещении объявления есть следующие ограничения:
    - Длина заголовка не меньше 4 и не больше 255 символов;
    - Длина текста описания не меньше 4 и не больше 1024 символов;
    - Цена должна быть в диапазоне от 49.99 до 999999999.99;
    - Формат изображения: .jpg,.jpeg. Дополнительно можно указать и другие форматы в файле .env параметра APP_ALLOWED_FILE_EXTENSIONS;
    - Размер изображения ограничен 5Мб. Возможно изменить в файле .env параметр APP_MAX_IMAGE_SIZE(указывается в Мб).
- отображения ленты объявлений. При отсутствии параметров в запросе, список отсортирован по дате добавления. 
    - Реализовано пагинация;
    - Дополнительно реализовано возможность указать количество объявлений на странице(для удобства дефолтное значение 2, но можно указать 10, 25, 50); 
    - Реализована возможность изменения типа и направления сортировки(по дате создания и цене);
    - Реализована возможность фильтрации по цене;
    - Для авторизованных пользователей дополнительно возвращается признак принадлежности объявления текущему пользователю.

## Инструкция по запуску проекта через Docker
1. Описать файл `.env` согласно примеру `.env.example`
2. Сервис собирается и запускается командой:
`docker-compose -f docker-compose.yml up -d`

## Примеры
Примеры возможных запросов:
- [Регистрация](#sign-up)
- [Авторизация](#sign-in)
- [Размещение нового объявления](#create)
- [Отображения ленты объявлений](#feed-ad)

### Регистрация <a name="sign-up"></a>
Запрос:
```
curl --location 'http://45.12.228.205:8080/auth/sign-up' \
--header 'Content-Type: application/json' \
--data '{
    "username": "example",
    "password": "Example123!"
}'
```
Ответ: `201 Created`
```json
{
    "id": 2,
    "username": "example"
}
```

### Авторизация <a name="sign-in"></a>
Запрос:
```
curl --location 'http://45.12.228.205:8080/auth/sign-in' \
--header 'Content-Type: application/json' \
--data '{
    "username": "example",
    "password": "Example123!"
}'
```
Ответ: `200 OK`

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI5Nzg0MDIsImlhdCI6MTc1Mjk3NjYwMiwic3ViIjoidG9rZW4iLCJVc2VySWQiOjJ9.xPT6vp43FyYolQhadc73VwUN0sE4DIb7XqDLpCGYgOk"
}
```

### Размещение нового объявления <a name="create"></a>
Запрос:
```
curl --location 'http://45.12.228.205:8080/api/v1/advertise/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI5Nzg0MDIsImlhdCI6MTc1Mjk3NjYwMiwic3ViIjoidG9rZW4iLCJVc2VySWQiOjJ9.xPT6vp43FyYolQhadc73VwUN0sE4DIb7XqDLpCGYgOk' \
--data '{
    "title": "Автожир Cavalon 914, синий",
    "description": "Год выпуска 2014",
    "image": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
    "price": 5000000.00
}'
```
Ответ: `201 Created`
```json
{
    "id": 9,
    "title": "Автожир Cavalon 914, синий",
    "description": "Год выпуска 2014",
    "image": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
    "price": 5000000,
    "user_id": 2,
    "created_at": "2025-07-20T01:58:22.002287Z"
}
```

### Отображения ленты объявлений <a name="feed-ad"></a>
Запрос:
```
curl --location 'http://45.12.228.205:8080/api/v1/advertise/feed-ad'
```
Ответ не авторизованного пользователя: `200 OK`
```json
[
    {
        "title": "Автожир Cavalon 914, синий",
        "description": "Год выпуска 2014",
        "image_url": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
        "price": 5000000,
        "author_username": "example"
    },
    {
        "title": "Автожир Cavalon 914, синий",
        "description": "Год выпуска 2014",
        "image_url": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
        "price": 6000000,
        "author_username": "azoma13"
    }
]
```
Ответ авторизованного пользователя: `200 OK`
```json 
[
    {
        "title": "Автожир Cavalon 914, синий",
        "description": "Год выпуска 2014",
        "image_url": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
        "price": 5000000,
        "author_username": "example",
        "is_author": true
    },
    {
        "title": "Автожир Cavalon 914, синий",
        "description": "Год выпуска 2014",
        "image_url": "https://sun9-73.userapi.com/s/v1/if1/hL1udYK9MvmB2aUMFRp-nRzE4Sx2lBqi0dmFrJEUn1UggQ4eF7Qqes7gQo8xZ_umc2G5xLZN.jpg?quality=96&as=32x24,48x36,72x54,108x81,160x120,240x180,360x270,480x360,540x405,640x480,720x540,1080x810,1280x960,1440x1080,2000x1500&from=bu&u=5dsrWc-KRlLjdtOxvWLQ-fTHc0pf-2jH-29Z9Y7aspo&cs=1080x0",
        "price": 6000000,
        "author_username": "azoma13"
    }
]
```
Запрос с указанием типа и направления сортировки, страницы, фильтрации по цене:
```
curl --location 'http://45.12.228.205:8080/api/v1/advertise/feed-ad?sort=price_desc&page=2&currency_price=60%3B1000000' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI5Nzg0MDIsImlhdCI6MTc1Mjk3NjYwMiwic3ViIjoidG9rZW4iLCJVc2VySWQiOjJ9.xPT6vp43FyYolQhadc73VwUN0sE4DIb7XqDLpCGYgOk'
```
Ответ: `200 OK`
```json 
[
    {
        "title": "Фитнес-браслет Xiaomi Mi Band 4",
        "description": "Вес (g): 11",
        "image_url": "https://sun9-48.userapi.com/s/v1/if1/sTN-RbVZY4T4V4TFPj4E76l_cencxkmO9DI1m3D0QHLVu2iKCX6-JwWmY6fDOkTUemwZqv_S.jpg?quality=96&as=32x21,48x32,72x48,108x72,160x107,240x160,360x240,480x320,540x360,604x403&from=bu&cs=604x0",
        "price": 2750,
        "author_username": "azoma13"
    },
    {
        "title": "Усилитель сигнала Xiaomi Mi Wi-Fi Amplifier PRO",
        "description": "Поддержка MIMO: есть",
        "image_url": "https://sun9-8.userapi.com/s/v1/ig1/b_Ba4KmfD5iDw-dmjRhp1TSRC7JRPwUSoTcw3vyuh61fuTZXgfbFNxkbdKOlNwxDLElFpRqV.jpg?quality=96&as=32x32,48x48,72x72,108x108,160x160,240x240,360x360,480x480,500x500&from=bu&u=IiUh3yyYsUpmTcKxeYXalLKaus8WKminhs6GKN3z-8o&cs=500x0",
        "price": 1190,
        "author_username": "azoma13"
    }
]
```