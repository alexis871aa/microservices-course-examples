#!/bin/bash

# Директория со скриптом
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR"
COMPOSE_DIR="$SCRIPT_DIR/../compose"

# Используем переданную переменную ENV_SUBST или системный envsubst
if [ -z "$ENV_SUBST" ]; then
  if ! command -v envsubst &> /dev/null; then
    echo "❌ Ошибка: envsubst не найден в системе и не передан через ENV_SUBST!"
    echo "Запустите скрипт через task env:generate"
    exit 1
  fi
  ENV_SUBST=envsubst
fi

# Функция загрузки переменных из .env файла
load_env_file() {
  local env_file="$1"

  if [ ! -f "$env_file" ]; then
    echo "❌ Ошибка: Файл $env_file не найден!"
  exit 1
fi

  echo "📋 Загружаем переменные из $env_file"
set -a
  source "$env_file"
set +a
}

# Функция для обработки шаблона и создания .env файла
process_template() {
  local service=$1
  local template="$TEMPLATE_DIR/${service}.env.template"
  local output="$COMPOSE_DIR/${service}/.env"

  echo "Обработка шаблона для сервиса $service..."

  if [ ! -f "$template" ]; then
    echo "⚠️ Шаблон $template не найден, пропускаем..."
    return 0
  fi

  # Создаем директорию, если она еще не существует
  mkdir -p "$(dirname "$output")"

  # Используем envsubst для замены переменных в шаблоне
  $ENV_SUBST < "$template" > "$output"

  echo "✅ Создан файл $output"
}

# Функция для обработки любого шаблона с указанием пути вывода
process_custom_template() {
  local template="$1"
  local output="$2"

  if [ ! -f "$template" ]; then
    echo "⚠️ Шаблон $template не найден, пропускаем..."
    return 1
  fi

  echo "🔄 Обработка шаблона: $template"

  # Создаем директорию, если она еще не существует
  mkdir -p "$(dirname "$output")"

  # Используем envsubst для замены переменных в шаблоне
  $ENV_SUBST < "$template" > "$output"

  echo "✅ Создан файл: $output"
  return 0
}

# Обработка сервисов (основной режим)
process_services() {
  # Проверяем наличие ENV_FILE
  if [ -z "$ENV_FILE" ]; then
    echo "❌ Ошибка: Переменная ENV_FILE не задана!"
    exit 1
  fi

  load_env_file "$ENV_FILE"

  # Определяем список сервисов из переменной окружения
  if [ -z "$SERVICES" ]; then
    echo "⚠️ Переменная SERVICES не задана. Нет сервисов для обработки."
    exit 0
  fi

  # Разделяем список сервисов по запятой
  IFS=',' read -ra services <<< "$SERVICES"
  echo "🔍 Обрабатываем сервисы: ${services[*]}"

  # Обрабатываем шаблоны для всех указанных сервисов
  success_count=0
  skip_count=0
  for service in "${services[@]}"; do
    process_template "$service"
    if [ -f "$TEMPLATE_DIR/${service}.env.template" ]; then
      ((success_count++))
    else
      ((skip_count++))
    fi
  done

  if [ $success_count -eq 0 ]; then
    echo "⚠️ Ни один .env файл не создан. Проверьте список сервисов и наличие шаблонов."
  else
    echo "🎉 Генерация завершена: $success_count файлов создано, $skip_count шаблонов пропущено"
  fi
}

# Обработка шаблонов мониторинга
process_monitoring() {
  echo "🔄 Генерация конфигурационных файлов для системы мониторинга..."

  # Проверяем наличие ENV_FILE
  if [ -z "$ENV_FILE" ]; then
    echo "❌ Ошибка: Переменная ENV_FILE не задана!"
    exit 1
  fi

  load_env_file "$ENV_FILE"

  # Список шаблонов для обработки в формате "шаблон:вывод"
  local templates=(
    "$COMPOSE_DIR/core/otel/otel-collector-config.template.yaml:$COMPOSE_DIR/core/otel/otel-collector-config.yaml"
    "$COMPOSE_DIR/core/prometheus/prometheus.template.yml:$COMPOSE_DIR/core/prometheus/prometheus.yml"
    "$COMPOSE_DIR/core/grafana/provisioning/datasources/prometheus.template.yml:$COMPOSE_DIR/core/grafana/provisioning/datasources/prometheus.yml"
    "$COMPOSE_DIR/core/grafana/provisioning/dashboards/dashboards.template.yml:$COMPOSE_DIR/core/grafana/provisioning/dashboards/dashboards.yml"
  )

  success_count=0
  for template_pair in "${templates[@]}"; do
    IFS=':' read -r template output <<< "$template_pair"
    if process_custom_template "$template" "$output"; then
      ((success_count++))
    fi
  done

  echo "🎉 Обработка шаблонов мониторинга завершена: $success_count файлов создано"
}

# Определяем режим работы скрипта
MODE="${1:-services}"

case "$MODE" in
  services)
    process_services
    ;;
  monitoring)
    process_monitoring
    ;;
  *)
    echo "❌ Неизвестный режим работы: $MODE"
    echo "Доступные режимы: services, monitoring"
    exit 1
    ;;
esac
