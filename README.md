# Url Shortener

## Цель проекта
- Научиться работать с k8s
- Практика масштабирования микросервисов
- Практика масштабирования бд
- Практика работа с брокером сообщений
- Сделать грамотное ограничение трафика


## Описание проекта
Проект состоит из трех частей
1. Link creator
2. Link redirector
3. Log consumer

### Link creator
[Readme.md](https://github.com/Zrossiz/shortener/blob/main/LinkCreator/README.md)
#### Функциональность:
- Принимать оригинальный адрес
- Создавать хеш
- Сохранять хеш в бд (postgres, redis)
- Обрабатывать дубли

### Link redirector
[Readme.md](https://github.com/Zrossiz/shortener/blob/main/LinkRedirector/README.md)  
#### Функциональность:
- Принимать хеш
- Искать оригинальный урл по хешу
- Перенаправлять на оригинальный адрес
- Отправлять сообщения kafka о переходе по ссылке

### Log consumer
[Readme.md](https://github.com/Zrossiz/shortener/blob/main/LogConsumer/README.md)  
#### Функциональность:
- Читать сообщения из кафка
- Сохранять данные в базу данных
- Получение данных для анализа

## Тесты
1. Чистая архитектура Робрета Мартина
2. Dependency inversion

Эти два принципа лежат в основе приложения, что позволяет писать изолированные unit тесты  
### Процент покрытия Redirector:
- 68.8%
- redirector/internal/delivery/rest         coverage: 92.3%
- redirector/internal/repository/postgresql coverage: 84.6% 
- redirector/internal/repository/redis      coverage: 83.3% 
- redirector/internal/service               coverage: 100.0% 
- redirector/pkg/config                     coverage: 100.0% 
- redirector/pkg/logger                     coverage: 92.3% 


### Процент покрытия Creator:
- 73.1%
- creator/internal/delivery/rest           coverage: 95.0%
- creator/internal/repository/postgresql   coverage: 86.7%
- creator/internal/repository/redis        coverage: 83.3%
- creator/internal/service                 coverage: 100.0%
- creator/pkg/config                       coverage: 100.0%
- creator/pkg/logger                       coverage: 92.3%
 
### Процент покрытия Consumer:
- 00.0%
- creator/pkg/config                       coverage: 100.0%
- creator/pkg/logger                       coverage: 92.3%