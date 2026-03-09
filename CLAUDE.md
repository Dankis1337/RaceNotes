# RaceNotes

## Проект
PWA "RaceNotes" — дневник велогонщика. Fullstack монорепо. PostgreSQL поднимается через Docker.

## Стек
- Backend: Go 1.21+, Gin, GORM, PostgreSQL 15+, JWT (bcrypt cost 14, токен 24ч)
- Frontend: Vue 3 (Composition API, script setup), Tailwind CSS, Vue Router 4, Pinia, Axios, Vite
- Инфраструктура: Docker, docker-compose, Nginx

## Структура
backend/cmd/server/main.go — точка входа
backend/internal/models/ — GORM модели
backend/internal/handlers/ — HTTP хендлеры
backend/internal/middleware/ — JWT middleware
backend/internal/services/ — бизнес-логика
frontend/src/api/ — Axios + API функции
frontend/src/components/ — переиспользуемые компоненты
frontend/src/views/ — страницы
frontend/src/router/ — маршруты
frontend/src/stores/ — Pinia
docker-compose.yml — postgres + backend + frontend

## Дизайн
Mobile-first, Tailwind, primary #83C082, accent #FAFFBA, bg #FFFFFF, text #2D2D2D

## API
POST /api/auth/register {name, username, email, password, height, weight}
POST /api/auth/login {username, password} -> {token}
GET/PUT /api/user/profile (JWT)
CRUD /api/setups (JWT)
CRUD /api/races?type=&is_completed=&setup_id= (JWT)
POST /api/calculator/tire-pressure (JWT)

## Модели
User: id, name, username (unique), email (unique), password (hashed), height (int см), weight (float кг), timestamps, soft delete
Setup: id, user_id, name, photo (nullable), bike_name, tires, components_description (nullable), timestamps, soft delete
Race: id, user_id, name, date, type (Road/MTB/Gravel/Cyclocross/Track), photo (nullable), setup_id (FK nullable), bike_name (nullable), tires (nullable), tire_pressure_front (float), tire_pressure_rear (float), temperature (int), conditions (Sunny/Cloudy/Rain/Snow), wind (None/Light/Moderate/Strong), road_conditions (Dry/Wet/Mud), nutrition_plan (text), result (string), rating (1-5), feelings (text), is_completed (bool default false), timestamps, soft delete

## Бизнес-правила
- setup_id != null -> bike_name/tires из Setup, иначе обязательны вручную
- Давление всегда в Race
- is_completed=false запланированная, true завершённая
- Мягкое удаление, пользователь видит только свои данные
- Калькулятор: упрощённая формула SRAM

## Правила
- Go: тонкие хендлеры, логика в services
- Vue: script setup, Tailwind, API через src/api/
- Git: коммиты по-английски, push после каждой фазы
- PostgreSQL через docker-compose, для локальной разработки тоже
