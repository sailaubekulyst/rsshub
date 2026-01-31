# RSSHub 📰

RSSHub — это CLI-приложение для агрегации RSS-лент.
Оно периодически загружает RSS-фиды, парсит статьи и сохраняет их в PostgreSQL, используя фоновой worker pool.

Проект разработан на Go и разворачивается с помощью Docker Compose.

---

## ✨ Возможности

* 📥 Добавление RSS-фидов
* 📃 Просмотр списка фидов
* 📰 Просмотр последних статей
* 🔄 Фоновая агрегация RSS-лент
* ⏱ Динамическое изменение интервала обновления
* 👷 Динамическое изменение количества воркеров
* 🛑 Корректное завершение работы (graceful shutdown)

---

## 🛠 Технологии

* Go
* PostgreSQL
* Docker / Docker Compose
* Concurrency (goroutines, channels, context)
* Worker Pool
* XML (RSS)

⚠️ Внешние библиотеки используются **только для PostgreSQL**.

---

## 📂 Структура проекта

```
.
├── cmd/
│   └── rsshub/           # CLI команды
├── internal/
│   ├── domain/           # Чистая бизнес-логика
│   │   ├── feed/
│   │   ├── article/
│   │   └── aggregator/
│   ├── aggregator/       # Реализация worker pool и ticker
│   ├── rss/              # RSS XML парсер
│   ├── postgres/         # PostgreSQL repositories
│   └── config/           # Загрузка конфигурации
├── migrations/           # SQL миграции
├── docker-compose.yml
├── main.go
└── README.md
```

---

## ⚙️ Конфигурация

Используются переменные окружения:

```env
# CLI App
CLI_APP_TIMER_INTERVAL=3m
CLI_APP_WORKERS_COUNT=3

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changeme
POSTGRES_DBNAME=rsshub
```

---

## 🐳 Запуск через Docker Compose

```bash
docker-compose up --build
```

Остановка:

```bash
docker-compose down
```

---

## 🏗 Сборка проекта

```bash
go build -o rsshub .
```

Проверка на data race:

```bash
go run -race main.go
```

---

## 🚀 CLI Команды

### ▶ Запуск фоновой агрегации

```bash
rsshub fetch
```

Вывод:

```
The background process for fetching feeds has started (interval = 3 minutes, workers = 3)
```

Остановка:

```bash
Ctrl + C
```

---

### ➕ Добавить RSS-фид

```bash
rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"
```

---

### 📃 Список RSS-фидов

```bash
rsshub list
```

С ограничением:

```bash
rsshub list --num 5
```

---

### ❌ Удалить RSS-фид

```bash
rsshub delete --name "tech-crunch"
```

---

### 📰 Показать статьи

```bash
rsshub articles --feed-name "tech-crunch"
```

С лимитом:

```bash
rsshub articles --feed-name "tech-crunch" --num 5
```

---

### ⏱ Изменить интервал обновления (на лету)

```bash
rsshub set-interval 2m
```

Вывод:

```
Interval of fetching feeds changed from 3 minutes to 2 minutes
```

⚠️ Работает только если `rsshub fetch` уже запущен.

---

### 👷 Изменить количество воркеров (на лету)

```bash
rsshub set-workers 5
```

Вывод:

```
Number of workers changed from 3 to 5
```

---

### ❓ Помощь

```bash
rsshub --help
```

---

## 🗄 База данных

### Таблица `feeds`

| Поле       | Тип       | Описание                   |
| ---------- | --------- | -------------------------- |
| id         | UUID      | Primary Key                |
| name       | TEXT      | Уникальное имя             |
| url        | TEXT      | URL RSS                    |
| created_at | TIMESTAMP | Дата добавления            |
| updated_at | TIMESTAMP | Дата последнего обновления |

---

### Таблица `articles`

| Поле         | Тип       | Описание        |
| ------------ | --------- | --------------- |
| id           | UUID      | Primary Key     |
| feed_id      | UUID      | FK → feeds.id   |
| title        | TEXT      | Заголовок       |
| link         | TEXT      | Ссылка          |
| description  | TEXT      | Описание        |
| published_at | TIMESTAMP | Дата публикации |
| created_at   | TIMESTAMP | Дата сохранения |

---

## 🔒 Безопасность и стабильность

* ❌ Нет data race (`sync.Mutex`, `context`)
* ❌ Нет утечек goroutines
* ❌ Нет дублирующих ticker'ов
* ✅ Graceful shutdown
* ✅ Один активный fetch-процесс

---

## 🧠 Рекомендации

* Не делайте слишком частые запросы к RSS-сайтам
* Следите за логами при запуске
* Используйте Ctrl+C для безопасной остановки

---

## 📌 Примеры RSS

* [https://techcrunch.com/feed/](https://techcrunch.com/feed/)
* [https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)
* [https://feeds.bbci.co.uk/news/world/rss.xml](https://feeds.bbci.co.uk/news/world/rss.xml)
* [https://www.theverge.com/rss/index.xml](https://www.theverge.com/rss/index.xml)
* [http://feeds.arstechnica.com/arstechnica/index](http://feeds.arstechnica.com/arstechnica/index)

---

## ✅ Статус

Проект соответствует требованиям задания и готов к защите 🎓
