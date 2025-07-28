Тестовое задание = Сервис для управления подписками, реализованный на Go. Поддерживает CRUDL-операции и расчет общей стоимости подписок за период. Использует PostgreSQL, логирование, конфигурацию через переменные окружения и Swagger-документацию.

Установка и запуск

Клонировать репозиторий.
Установить зависимости:go mod tidy


Применить миграции:goose -dir internal/migrations postgres "user=postgres dbname=subscriptions password=password host=localhost" up


Сгенерировать Swagger-документацию:swag init -g cmd/subscriptionsservice/main.go


Запустить:docker-compose up


Эндпоинты и примеры запросов:

POST
/subscriptions
Создать подписку


GET
/subscriptions/{id}
Получить подписку по ID


PUT
/subscriptions/{id}
Обновить подписку


DELETE
/subscriptions/{id}
Удалить подписку


GET
/subscriptions
Получить список подписок


POST
/subscriptions/total-cost
Рассчитать общую стоимость
