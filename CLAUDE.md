# RaceNotes

## Проект
PWA "RaceNotes" — дневник велогонщика. Учёт сетапов, планирование гонок, анализ результатов. Акцент на давление в шинах, погоду, план питания. Fullstack монорепо.

## Стек
### Backend
- Go 1.21+, Gin, GORM, PostgreSQL 15+, JWT (bcrypt cost 14, токен 24ч)

### Frontend
- Vue 3 (Composition API, script setup), Tailwind CSS, Vue Router 4, Pinia, Axios, Vite

### Структура монорепо
backend/
  cmd/server/main.go
  internal/
    models/
    handlers/
    middleware/
    services/
  go.mod
frontend/
  src/
    api/
    components/
    views/
    router/
    stores/
    App.vue
    main.js
  package.json
  vite.config.js
docker-compose.yml
nginx.conf

## Дизайн фронтенда
- Mobile-first, стандартные Tailwind компоненты
- Primary: #83C082, Accent: #FAFFBA, Background: #FFFFFF, Text: #2D2D2D
- rounded-xl на карточках, rounded-lg на кнопках, shadow-md
- Чистый, спортивный, минималистичный

## API
- POST /api/auth/register {name, username, email, password, height, weight}
- POST /api/auth/login {username, password} -> {token}
- GET /api/user/profile (JWT)
- PUT /api/user/profile (JWT)
- CRUD /api/setups (JWT)
- CRUD /api/races?type=&is_completed=&setup_id= (JWT)
- POST /api/calculator/tire-pressure (JWT)

## Модели
User: id, name, username (unique), email (unique), password (hashed), height (int, см), weight (float, кг), created_at, updated_at, deleted_at
Setup: id, user_id, name, photo (URL nullable), bike_name, tires, components_description (text nullable), created_at, updated_at, deleted_at
Race: id, user_id, name, date, type (Road/MTB/Gravel/Cyclocross/Track), photo (nullable), setup_id (FK nullable), bike_name (nullable), tires (nullable), tire_pressure_front (float nullable), tire_pressure_rear (float nullable), temperature (int nullable), conditions (Sunny/Cloudy/Rain/Snow nullable), wind (None/Light/Moderate/Strong nullable), road_conditions (Dry/Wet/Mud nullable), nutrition_plan (text nullable), result (string nullable), rating (int 1-5 nullable), feelings (text nullable), is_completed (bool default false), created_at, updated_at, deleted_at

## Бизнес-правила
- setup_id != null -> bike_name и tires берутся из Setup
- setup_id = null -> bike_name и tires обязательны
- Давление всегда в Race
- is_completed=false -> запланированная гонка
- is_completed=true -> завершённая гонка
- Калькулятор: упрощённая формула SRAM, возвращает front/rear pressure + рекомендации
- Мягкое удаление (deleted_at) для setups и races
- Пользователь видит только свои данные

## Правила
- Go: тонкие хендлеры, логика в services, GORM для БД
- Vue: Composition API script setup, Tailwind, API через src/api/ слой
- Git: коммиты по-английски, push после каждой фазы
