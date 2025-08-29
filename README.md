## Запуск приложения
В папке backend создайте файл .env, аналогичный .env.example

Запустите проект с помощью команды:

   ```bash
   docker compose up --build -d
   ```

После запуска API будет доступен по адресу: [http://localhost:8080](http://localhost:8080), где также будет доступна документация Swagger.

---

## Тестирование

Есть тесты game_service

Запуск
   ```bash
   go test -v
   ```
