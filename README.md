# REST API с JWT аутентификацией

## Маршруты
- POST   /auth/signup              - регистрация
- GET    /auth/signin              - вход в аккаунт
- GET    /auth/logout              - выход из аккаунта
- GET    /users/:user_id           - получение пользователя по id
- GET    /users/                   - получение списка всех пользователей
- PATCH  /users/:user_id           - изменение пользователя по id
- DELETE /users/:user_id           - удаление пользователя по id

## JWT 
- Access токен тип JWT, алгоритм SHA512;
- В Payload Access токена хранится guid пользователя и время, после которого истекает срок действия токена;
- Refresh токен тип произвольный. Состоит из двух частей: строки из 16 случайных символов и 8 последних символов Access токена, который был сгенерирован вместе с ним;
- В куки Refresh токен хранится в base64
- В БД Refresh токен хранится в виде bcrypt хеша
- После операции Refresh, Refresh токен в бд заменяется новым
- Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним. Это реализовано путём сравнения последних 8 символов у токенов.
