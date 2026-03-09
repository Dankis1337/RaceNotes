import { ref, computed } from 'vue'

const STORAGE_KEY = 'racenotes_lang'

const currentLang = ref(localStorage.getItem(STORAGE_KEY) || 'ru')

export function setLanguage(lang) {
  currentLang.value = lang
  localStorage.setItem(STORAGE_KEY, lang)
}

export function getLanguage() {
  return currentLang.value
}

const translations = {
  // ===== Navigation =====
  nav_races: { en: 'Races', ru: 'Гонки' },
  nav_setups: { en: 'Setups', ru: 'Сетапы' },
  nav_calculator: { en: 'Calculator', ru: 'Калькулятор' },
  nav_profile: { en: 'Profile', ru: 'Профиль' },

  // ===== Offline =====
  offline_mode: { en: 'Offline mode', ru: 'Офлайн режим' },
  syncing: { en: 'Syncing...', ru: 'Синхронизация...' },
  pending_sync: { en: '{n} pending sync', ru: '{n} ожидает синхронизации' },

  // ===== Login =====
  app_name: { en: 'RaceNotes', ru: 'RaceNotes' },
  login_subtitle: { en: "Cyclist's race diary", ru: 'Дневник велогонщика' },
  username: { en: 'Username', ru: 'Логин' },
  password: { en: 'Password', ru: 'Пароль' },
  enter_username: { en: 'Enter username', ru: 'Введите логин' },
  enter_password: { en: 'Enter password', ru: 'Введите пароль' },
  signing_in: { en: 'Signing in...', ru: 'Вход...' },
  sign_in: { en: 'Sign In', ru: 'Войти' },
  no_account: { en: "Don't have an account?", ru: 'Нет аккаунта?' },
  sign_up: { en: 'Sign Up', ru: 'Регистрация' },
  login_failed: { en: 'Login failed', ru: 'Ошибка входа' },

  // ===== Register =====
  create_account: { en: 'Create your account', ru: 'Создайте аккаунт' },
  name: { en: 'Name', ru: 'Имя' },
  your_name: { en: 'Your name', ru: 'Ваше имя' },
  choose_username: { en: 'Choose username', ru: 'Придумайте логин' },
  email: { en: 'Email', ru: 'Email' },
  email_placeholder: { en: 'email@example.com', ru: 'email@example.com' },
  min_6_chars: { en: 'Min 6 characters', ru: 'Минимум 6 символов' },
  height_cm: { en: 'Height (cm)', ru: 'Рост (см)' },
  weight_kg: { en: 'Weight (kg)', ru: 'Вес (кг)' },
  creating_account: { en: 'Creating account...', ru: 'Создание...' },
  already_have_account: { en: 'Already have an account?', ru: 'Уже есть аккаунт?' },
  registration_failed: { en: 'Registration failed', ru: 'Ошибка регистрации' },

  // ===== Races List =====
  hi_rider: { en: 'Hi, {name}!', ru: 'Привет, {name}!' },
  your_race_diary: { en: 'Your race diary', ru: 'Твой дневник гонок' },
  yesterday_races: { en: 'You have {n} race(s) from yesterday without results.', ru: 'У вас {n} гонок вчера без результатов.' },
  fill_results: { en: 'Fill in results', ru: 'Заполнить результаты' },
  later: { en: 'Later', ru: 'Позже' },
  status: { en: 'Status', ru: 'Статус' },
  type: { en: 'Type', ru: 'Тип' },
  setup: { en: 'Setup', ru: 'Сетап' },
  all: { en: 'All', ru: 'Все' },
  planned: { en: 'Planned', ru: 'Запланирована' },
  completed: { en: 'Completed', ru: 'Завершена' },
  done: { en: 'Done', ru: 'Готово' },
  no_races_yet: { en: 'No races yet', ru: 'Гонок пока нет' },
  tap_plus_race: { en: 'Tap + to add your first race', ru: 'Нажмите + чтобы добавить первую гонку' },

  // ===== Race Card =====
  // (uses planned/done from above)

  // ===== Race Detail =====
  back: { en: 'Back', ru: 'Назад' },
  bike: { en: 'Bike', ru: 'Велосипед' },
  tires: { en: 'Tires', ru: 'Покрышки' },
  front_pressure: { en: 'Front Pressure', ru: 'Давление спереди' },
  rear_pressure: { en: 'Rear Pressure', ru: 'Давление сзади' },
  other_components: { en: 'Other Components', ru: 'Другие компоненты' },
  weather_conditions: { en: 'Weather & Conditions', ru: 'Погода и условия' },
  wind_label: { en: 'Wind', ru: 'Ветер' },
  nutrition_plan: { en: 'Nutrition Plan', ru: 'План питания' },
  results: { en: 'Results', ru: 'Результаты' },
  result: { en: 'Result', ru: 'Результат' },
  rating: { en: 'Rating', ru: 'Оценка' },
  feelings: { en: 'Feelings', ru: 'Ощущения' },
  fill_in_results: { en: 'Fill in Results', ru: 'Заполнить результаты' },
  delete_race_q: { en: 'Delete Race?', ru: 'Удалить гонку?' },
  cannot_undo: { en: 'This action cannot be undone.', ru: 'Это действие нельзя отменить.' },
  cancel: { en: 'Cancel', ru: 'Отмена' },
  delete_btn: { en: 'Delete', ru: 'Удалить' },
  race_deleted: { en: 'Race deleted', ru: 'Гонка удалена' },

  // ===== Race Form =====
  edit_race: { en: 'Edit Race', ru: 'Редактировать гонку' },
  new_race: { en: 'New Race', ru: 'Новая гонка' },
  basic_info: { en: 'Basic Info', ru: 'Основная информация' },
  race_name: { en: 'Race Name', ru: 'Название гонки' },
  race_name_placeholder: { en: 'e.g. Spring Classic 2026', ru: 'напр. Весенняя классика 2026' },
  date: { en: 'Date', ru: 'Дата' },
  photo: { en: 'Photo', ru: 'Фото' },
  bike_setup: { en: 'Bike Setup', ru: 'Сетап велосипеда' },
  use_saved_setup: { en: 'Use saved setup', ru: 'Использовать сохранённый сетап' },
  select_setup: { en: 'Select Setup', ru: 'Выберите сетап' },
  choose_setup: { en: 'Choose a setup...', ru: 'Выберите сетап...' },
  bike_name: { en: 'Bike Name', ru: 'Название велосипеда' },
  bike_name_placeholder: { en: 'e.g. Canyon Aeroad', ru: 'напр. Canyon Aeroad' },
  tires_placeholder: { en: 'e.g. Continental GP5000 28mm', ru: 'напр. Continental GP5000 28mm' },
  tire_pressure: { en: 'Tire Pressure', ru: 'Давление в шинах' },
  calculate: { en: 'Calculate', ru: 'Рассчитать' },
  front_bar: { en: 'Front (bar)', ru: 'Переднее (bar)' },
  rear_bar: { en: 'Rear (bar)', ru: 'Заднее (bar)' },
  other_components_section: { en: 'Other Components', ru: 'Другие компоненты' },
  other_components_placeholder: { en: 'e.g. Saddle height, handlebar position, gearing...', ru: 'напр. Высота седла, положение руля, передачи...' },
  temperature: { en: 'Temperature (°C)', ru: 'Температура (°C)' },
  weather: { en: 'Weather', ru: 'Погода' },
  road: { en: 'Road', ru: 'Дорога' },
  nutrition_plan_section: { en: 'Nutrition Plan', ru: 'План питания' },
  nutrition_placeholder: { en: 'e.g. Gel every 30 min, 750ml water per hour...', ru: 'напр. Гель каждые 30 мин, 750 мл воды в час...' },
  is_completed: { en: 'Completed', ru: 'Завершена' },
  result_placeholder: { en: 'e.g. 3rd place, 2:45:30', ru: 'напр. 3-е место, 2:45:30' },
  feelings_placeholder: { en: 'How did you feel during the race?', ru: 'Как вы себя чувствовали во время гонки?' },
  saving: { en: 'Saving...', ru: 'Сохранение...' },
  save_changes: { en: 'Save Changes', ru: 'Сохранить изменения' },
  create_race: { en: 'Create Race', ru: 'Создать гонку' },
  race_updated: { en: 'Race updated', ru: 'Гонка обновлена' },
  race_created: { en: 'Race created', ru: 'Гонка создана' },
  failed_save_race: { en: 'Failed to save race', ru: 'Не удалось сохранить гонку' },

  // ===== Calculator Modal (in Race Form) =====
  tire_pressure_calculator: { en: 'Tire Pressure Calculator', ru: 'Калькулятор давления в шинах' },
  rider_weight: { en: 'Rider Weight (kg)', ru: 'Вес райдера (кг)' },
  bike_weight: { en: 'Bike Weight (kg)', ru: 'Вес велосипеда (кг)' },
  tire_width: { en: 'Tire Width (mm)', ru: 'Ширина покрышки (мм)' },
  tire_type: { en: 'Tire Type', ru: 'Тип покрышки' },
  surface: { en: 'Surface', ru: 'Покрытие' },
  conditions: { en: 'Conditions', ru: 'Условия' },
  calculating: { en: 'Calculating...', ru: 'Расчёт...' },
  recommended_pressure: { en: 'Recommended Pressure', ru: 'Рекомендуемое давление' },
  front: { en: 'Front', ru: 'Переднее' },
  rear: { en: 'Rear', ru: 'Заднее' },
  tips: { en: 'Tips', ru: 'Советы' },
  apply_to_race: { en: 'Apply to Race', ru: 'Применить к гонке' },
  calculation_failed: { en: 'Calculation failed', ru: 'Ошибка расчёта' },
  calculator_error: { en: 'Calculator error', ru: 'Ошибка калькулятора' },

  // ===== Tire Type options =====
  clincher: { en: 'Clincher', ru: 'Клинчер' },
  tubeless: { en: 'Tubeless', ru: 'Бескамерные' },
  tubular: { en: 'Tubular', ru: 'Трубки' },

  // ===== Surface options =====
  surface_road: { en: 'Road', ru: 'Шоссе' },
  surface_gravel: { en: 'Gravel', ru: 'Гравий' },
  surface_mixed: { en: 'Mixed', ru: 'Микс' },
  surface_cobblestone: { en: 'Cobblestone', ru: 'Брусчатка' },

  // ===== Conditions options =====
  cond_dry: { en: 'Dry', ru: 'Сухо' },
  cond_wet: { en: 'Wet', ru: 'Мокро' },
  cond_mud: { en: 'Mud', ru: 'Грязь' },
  cond_snow: { en: 'Snow', ru: 'Снег' },

  // ===== Weather options (race form chips) =====
  sunny: { en: 'Sunny', ru: 'Солнечно' },
  cloudy: { en: 'Cloudy', ru: 'Облачно' },
  rain: { en: 'Rain', ru: 'Дождь' },
  snow: { en: 'Snow', ru: 'Снег' },

  // ===== Wind options =====
  wind_none: { en: 'None', ru: 'Нет' },
  wind_light: { en: 'Light', ru: 'Слабый' },
  wind_moderate: { en: 'Moderate', ru: 'Умеренный' },
  wind_strong: { en: 'Strong', ru: 'Сильный' },

  // ===== Road conditions =====
  road_dry: { en: 'Dry', ru: 'Сухо' },
  road_wet: { en: 'Wet', ru: 'Мокро' },
  road_mud: { en: 'Mud', ru: 'Грязь' },

  // ===== Race types =====
  type_road: { en: 'Road', ru: 'Шоссе' },
  type_mtb: { en: 'MTB', ru: 'MTB' },
  type_gravel: { en: 'Gravel', ru: 'Гравел' },
  type_cyclocross: { en: 'Cyclocross', ru: 'Циклокросс' },
  type_track: { en: 'Track', ru: 'Трек' },

  // ===== Setups List =====
  my_setups: { en: 'My Setups', ru: 'Мои сетапы' },
  no_setups_yet: { en: 'No setups yet', ru: 'Сетапов пока нет' },
  create_first_setup: { en: 'Create your first bike setup', ru: 'Создайте свой первый сетап' },
  components: { en: 'Components', ru: 'Компоненты' },

  // ===== Setup Detail =====
  delete_setup_q: { en: 'Delete Setup?', ru: 'Удалить сетап?' },
  setup_deleted: { en: 'Setup deleted', ru: 'Сетап удалён' },

  // ===== Setup Form =====
  edit_setup: { en: 'Edit Setup', ru: 'Редактировать сетап' },
  new_setup: { en: 'New Setup', ru: 'Новый сетап' },
  setup_name: { en: 'Setup Name', ru: 'Название сетапа' },
  setup_name_placeholder: { en: 'e.g. Road Race Setup', ru: 'напр. Гоночный сетап' },
  bike_name_setup_placeholder: { en: 'e.g. Canyon Aeroad CF SLX', ru: 'напр. Canyon Aeroad CF SLX' },
  components_description: { en: 'Components Description', ru: 'Описание компонентов' },
  components_placeholder: { en: 'Shimano Ultegra Di2, Zipp 303 wheels...', ru: 'Shimano Ultegra Di2, колёса Zipp 303...' },
  create_setup: { en: 'Create Setup', ru: 'Создать сетап' },
  setup_updated: { en: 'Setup updated', ru: 'Сетап обновлён' },
  setup_created: { en: 'Setup created', ru: 'Сетап создан' },
  failed_save_setup: { en: 'Failed to save setup', ru: 'Не удалось сохранить сетап' },

  // ===== Profile =====
  profile_updated: { en: 'Profile updated', ru: 'Профиль обновлён' },
  failed_update_profile: { en: 'Failed to update profile', ru: 'Не удалось обновить профиль' },
  username_cannot_change: { en: 'Username cannot be changed', ru: 'Логин нельзя изменить' },
  edit_profile: { en: 'Edit Profile', ru: 'Редактировать профиль' },
  save: { en: 'Save', ru: 'Сохранить' },
  log_out: { en: 'Log Out', ru: 'Выйти' },
  language: { en: 'Language', ru: 'Язык' },

  // ===== Photo Upload =====
  file_too_large: { en: 'File must be under 10 MB', ru: 'Файл должен быть менее 10 МБ' },
  upload_failed: { en: 'Upload failed', ru: 'Ошибка загрузки' },
  uploading: { en: 'Uploading...', ru: 'Загрузка...' },
  tap_add_photo: { en: 'Tap to add photo', ru: 'Нажмите чтобы добавить фото' },
}

export function t(key, params) {
  const entry = translations[key]
  if (!entry) return key
  let text = entry[currentLang.value] || entry.en || key
  if (params) {
    Object.keys(params).forEach(k => {
      text = text.replace(`{${k}}`, params[k])
    })
  }
  return text
}

export function useI18n() {
  const lang = computed(() => currentLang.value)

  function tr(key, params) {
    // Access currentLang.value to create reactivity
    const _ = currentLang.value
    return t(key, params)
  }

  return { t: tr, lang, setLanguage, getLanguage }
}
