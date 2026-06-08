# Документация REST API сервиса поиска мест отдыха

## Общая информация

### Технологический стек
- **Язык программирования**: Go 1.21+
- **Веб-фреймворк**: Chi Router
- **ORM**: GORM
- **База данных**: PostgreSQL
- **Аутентификация**: JWT токены

### Архитектура проекта
Проект построен по принципам чистой архитектуры с разделением на слои:
- **Repository** - работа с базой данных
- **Service** - бизнес-логика
- **Handler** - обработка HTTP запросов
- **Middleware** - промежуточные обработчики

## Структура проекта
```
api_tp/
├── internal/
│   ├── config/      # Конфигурация приложения
│   ├── handlers/    # HTTP обработчики
│   ├── middleware/  # Промежуточные обработчики
│   ├── models/      # Модели данных
│   ├── repository/  # Репозитории для работы с БД
│   ├── server/      # Конфигурация сервера
│   └── service/     # Бизнес-логика
├── pkg/
│   └── database/    # Подключение к БД
└── main.go          # Точка входа
```

## Текущие эндпоинты API

### Версия API: v1

#### 1. Проверка здоровья сервиса
```
GET /health
GET /v1/check
GET /v1/auth/check
GET /v1/api/users/check
```
**Ответ:**
```json
{
  "status": "healthy",
  "timestamp": "Tue, 02 Jan 2024 10:00:00 UTC"
}
```

#### 2. Аутентификация (публичные маршруты)

##### Регистрация пользователя
```
POST /v1/auth/register
```
**Тело запроса (ожидается):**
```json
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "Имя пользователя"
}
```

##### Авторизация пользователя
```
POST /v1/auth/login
```
**Тело запроса (ожидается):**
```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

#### 3. Пользователи (защищенные маршруты - требуется JWT токен)

##### Получение всех пользователей
```
GET /v1/api/users/
```
**Заголовок:**
```
Authorization: Bearer <JWT_TOKEN>
```

##### Создание пользователя
```
POST /v1/api/users/
```
**Заголовок:**
```
Authorization: Bearer <JWT_TOKEN>
```

##### Получение пользователя по ID
```
GET /v1/api/users/{id}
```
**Заголовок:**
```
Authorization: Bearer <JWT_TOKEN>
```

## Модель данных

### Пользователь (UserT)
```go
type UserT struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Аутентификация
Защищенные маршруты требуют JWT токен, который должен передаваться в заголовке:
```
Authorization: Bearer <ваш_jwt_токен>
```

## Запуск сервиса

### Требования
1. Go 1.21 или выше
2. PostgreSQL 12+
3. Установленные переменные окружения (детали в `config` пакете)

### Запуск
```bash
go run main.go
```

Сервер запускается на порту, указанном в конфигурации (по умолчанию 8080).

## Будущие реализации эндпоинтов

### 1. Работа с местами отдыха

#### Категории мест
```
GET    /v1/api/categories          # Получение всех категорий
POST   /v1/api/categories          # Создание категории
GET    /v1/api/categories/{id}     # Получение категории по ID
PUT    /v1/api/categories/{id}     # Обновление категории
DELETE /v1/api/categories/{id}     # Удаление категории
```

#### Сами места отдыха
```
GET    /v1/api/places              # Получение всех мест
POST   /v1/api/places              # Создание места
GET    /v1/api/places/{id}         # Получение места по ID
PUT    /v1/api/places/{id}         # Обновление места
DELETE /v1/api/places/{id}         # Удаление места
GET    /v1/api/places/search       # Поиск мест по параметрам
GET    /v1/api/places/category/{id} # Места по категории
```

### 2. Отзывы и рейтинги

```
POST   /v1/api/places/{id}/reviews     # Добавление отзыва
GET    /v1/api/places/{id}/reviews     # Получение отзывов места
PUT    /v1/api/reviews/{id}           # Обновление отзыва
DELETE /v1/api/reviews/{id}           # Удаление отзыва
POST   /v1/api/places/{id}/rating     # Добавление/обновление рейтинга
```

### 3. Избранное

```
GET    /v1/api/user/favorites         # Получение избранного пользователя
POST   /v1/api/places/{id}/favorite  # Добавление в избранное
DELETE /v1/api/places/{id}/favorite  # Удаление из избранного
```

### 4. Фильтрация и поиск

```
GET /v1/api/places/filter             # Расширенная фильтрация
```
**Параметры запроса:**
- `category` - ID категории
- `min_price`, `max_price` - диапазон цен
- `rating` - минимальный рейтинг
- `location` - географическое положение
- `amenities` - удобства (wi-fi, парковка и т.д.)

### 5. Административные функции

```
GET    /v1/admin/users              # Получение всех пользователей (админ)
PUT    /v1/admin/users/{id}/status  # Изменение статуса пользователя
GET    /v1/admin/places/pending    # Места ожидающие модерации
PUT    /v1/admin/places/{id}/approve # Одобрение места
DELETE /v1/admin/places/{id}        # Удаление места (админ)
```

### 6. Статистика и аналитика

```
GET /v1/api/analytics/popular-places    # Популярные места
GET /v1/api/analytics/user-activity     # Активность пользователя
GET /v1/admin/analytics/overview       # Общая статистика (админ)
```

### 7. Геолокационные функции

```
GET /v1/api/places/nearby              # Ближайшие места
POST /v1/api/location/suggest         # Предложения по локации
```

### 8. Медиа-контент

```
POST   /v1/api/places/{id}/photos     # Добавление фотографий
DELETE /v1/api/photos/{id}           # Удаление фотографии
GET    /v1/api/places/{id}/photos    # Получение фотографий места
```

## Планируемые модели данных

### Место отдыха (Place)
```go
type Place struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Description string
    CategoryID  uint      // Связь с категорией
    Address     string
    Latitude    float64
    Longitude   float64
    PriceRange  string    // "низкий", "средний", "высокий"
    Rating      float64   `gorm:"default:0"`
    CreatedBy   uint      // ID пользователя-создателя
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Отзыв (Review)
```go
type Review struct {
    ID        uint      `gorm:"primaryKey"`
    PlaceID   uint      `gorm:"not null"`
    UserID    uint      `gorm:"not null"`
    Rating    int       `gorm:"check:rating >= 1 AND rating <= 5"`
    Comment   string
    CreatedAt time.Time
}
```

### Категория (Category)
```go
type Category struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"uniqueIndex;not null"`
    Description string
    Icon        string    // URL иконки
}
```

## Безопасность и валидация

1. **Валидация входных данных** на всех эндпоинтах
2. **Rate limiting** для предотвращения DDoS атак
3. **CORS** настройка для веб-клиентов
4. **Хеширование паролей** с использованием bcrypt
5. **JWT токены** с ограниченным временем жизни
6. **Ролевая модель** (пользователь, модератор, администратор)

## Мониторинг и логирование

1. **Structured logging** с использованием zap или logrus
2. **Метрики Prometheus** для мониторинга
3. **Health checks** расширенные
4. **Tracing** с использованием OpenTelemetry

## Деплоймент

### Docker
Планируется создание Dockerfile и docker-compose для развертывания:
- Приложение API
- PostgreSQL
- Redis (для кэширования)
- Nginx (как reverse proxy)

### Kubernetes
Конфигурации для развертывания в Kubernetes кластере.

---

*Примечание: Это предварительная документация. Реализация новых эндпоинтов будет сопровождаться обновлением документации.*