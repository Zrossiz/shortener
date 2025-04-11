# Url Shortener

## Цель проекта
- Научиться работать с k8s
- Практика масштабирования микросервисов
- Практика масштабирования бд
- Практика работы с брокером сообщений
- Сделать грамотное ограничение трафика

## Описание проекта
Проект состоит из трех частей
1. Link creator
2. Link redirector
3. Log consumer

### Быстрый старт
В каждом проекте есть .env.example.   
Скопируйте переменные в .env  
Запустите через docker compose -f docker-compose.local.yaml up --build  

### Link creator
[Readme.md](https://github.com/Zrossiz/url-shortener.generator/README.md)
#### Функциональность:
- Принимать оригинальный адрес
- Создавать хеш
- Сохранять хеш в бд (postgres, redis)
- Обрабатывать дубли

### Link redirector
[Readme.md](https://github.com/Zrossiz/url-shortener.logger/README.md)  
#### Функциональность:
- Принимать хеш
- Искать оригинальный урл по хешу
- Перенаправлять на оригинальный адрес
- Отправлять сообщения kafka о переходе по ссылке

### Log consumer
[Readme.md](https://github.com/Zrossiz/url-shortener.redirector/README.md)  
#### Функциональность:
- Читать сообщения из кафка
- Сохранять данные в базу данных
- Получение данных для анализа

## Тесты
1. Чистая архитектура Робрета Мартина
2. Dependency inversion

Эти два принципа лежат в основе приложения, что позволяет писать изолированные unit тесты  
### Процент покрытия Redirector:
- 65+%

### Процент покрытия Creator:
- 70+%

### Процент покрытия Log Consumer:
- 60+%