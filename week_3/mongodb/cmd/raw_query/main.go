package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/olezhek28/microservices-course-examples/week_3/mongodb/internal/model"
)

func main() {
	// Создаем контекст для управления выполнением операций MongoDB
	// Контекст позволяет управлять жизненным циклом запросов, включая таймауты и отмену
	ctx := context.Background()

	// Загружаем переменные окружения из файла .env
	// Это безопасный способ хранения конфигурационных данных вне кода
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Не удалось загрузить файл .env: %v\n", err)
		return
	}

	// Получаем строку подключения к MongoDB из переменной окружения
	// Строка подключения содержит всю информацию, необходимую для соединения с базой данных
	// Например: mongodb://username:password@localhost:27017/database
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Println("Ошибка: переменная окружения MONGO_URI не установлена")
		return
	}

	// Подключение к MongoDB
	// --------------------------------

	// Создаем клиент MongoDB с настройками из строки подключения
	// Клиент - это основной объект для работы с MongoDB, через него выполняются все операции
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v\n", err)
		return
	}

	// Гарантируем отключение клиента при завершении работы программы
	// Это важно для корректного освобождения ресурсов и закрытия соединений
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("Ошибка при отключении от MongoDB: %v\n", cerr)
		}
	}()

	// Проверяем соединение с базой данных отправкой ping-запроса
	// Если сервер недоступен или есть проблемы с подключением, получим ошибку
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB недоступна, ошибка ping: %v\n", err)
		return
	}
	log.Println("Успешное подключение к MongoDB")

	// Настройка коллекции и индексов
	// --------------------------------

	// Получаем доступ к коллекции "notes" в базе данных "example"
	// В MongoDB база данных и коллекции создаются автоматически при первом обращении
	// База данных - это контейнер для коллекций (аналог базы данных в SQL)
	// Коллекция - это группа документов (аналог таблицы в SQL)
	collection := client.Database("example").Collection("notes")

	// Создаем индекс для поля "title"
	// Индексы ускоряют поиск по определенным полям документов
	// Здесь мы создаем простой индекс для поля "title" в порядке возрастания (1)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}}, // 1 означает индекс по возрастанию, -1 был бы по убыванию
		Options: options.Index().SetUnique(false), // Индекс не уникальный, могут быть записи с одинаковым title
	}

	// Создаем индекс если он еще не существует
	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Ошибка создания индекса: %v\n", err)
		return
	}
	log.Printf("Создан индекс: %s\n", indexName)

	// Операция вставки (Create)
	// --------------------------------

	// Создаем новую заметку с тестовыми данными
	// Для ID мы не указываем значение - MongoDB автоматически сгенерирует ObjectID
	// Для полей title и body используем случайные данные из библиотеки gofakeit
	note := model.Note{
		Title:     gofakeit.City(),           // Случайный город в качестве заголовка
		Body:      gofakeit.Address().Street, // Случайная улица в качестве содержания
		CreatedAt: time.Now(),                // Текущее время как время создания
		// UpdatedAt не указываем, будет nil, так как это новая запись
	}

	// Вставляем документ в коллекцию
	// InsertOne вставляет один документ и возвращает его ID
	insertResult, err := collection.InsertOne(ctx, note)
	if err != nil {
		log.Printf("Ошибка вставки заметки: %v\n", err)
		return
	}

	// Преобразуем возвращенный ID в ObjectID для дальнейшего использования
	// MongoDB использует ObjectID как стандартный тип для уникальных идентификаторов
	insertedID := insertResult.InsertedID.(primitive.ObjectID)
	log.Printf("Заметка успешно добавлена, ID: %s\n", insertedID.Hex())

	// Операция чтения всех документов (Read All)
	// --------------------------------

	// Получаем все заметки из коллекции
	// Find с пустым фильтром bson.M{} вернет все документы
	// bson.M - это удобный способ создания BSON-документов в Go (аналог JSON)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v\n", err)
		return
	}

	// Закрываем курсор при завершении, чтобы освободить ресурсы
	// Курсор - это итератор по результатам запроса
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Ошибка при закрытии курсора: %v\n", cerr)
		}
	}()

	// Декодируем все документы из курсора в срез заметок
	// Метод All удобен, когда нужно получить все результаты сразу
	var notes []model.Note
	if err = cursor.All(ctx, &notes); err != nil {
		log.Printf("Ошибка декодирования заметок: %v\n", err)
		return
	}

	log.Printf("Найдено заметок: %d\n", len(notes))

	// Выводим информацию о каждой заметке
	for _, n := range notes {
		// Проверяем, имеет ли поле UpdatedAt значение
		// Так как UpdatedAt - это указатель, он может быть nil
		updatedAtStr := "не обновлялась"
		if n.UpdatedAt != nil {
			updatedAtStr = n.UpdatedAt.String()
		}

		log.Printf("ID: %s, Заголовок: %s, Содержание: %s, Создана: %v, Обновлена: %s\n",
			n.ID.Hex(), n.Title, n.Body, n.CreatedAt.Format(time.RFC3339), updatedAtStr)
	}

	// Операция обновления (Update)
	// --------------------------------

	// Обновляем созданную ранее заметку
	now := time.Now()

	// UpdateOne обновляет первый документ, соответствующий фильтру
	// Первый параметр - фильтр для поиска документа (по ID)
	// Второй параметр - операции обновления:
	//   $set - устанавливает новые значения для указанных полей
	updateResult, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": insertedID}, // Фильтр: ищем документ с нашим ID
		bson.M{
			"$set": bson.M{ // Оператор $set меняет значения указанных полей
				"title":      gofakeit.City(),           // Новый случайный заголовок
				"body":       gofakeit.Address().Street, // Новое случайное содержание
				"updated_at": now,                       // Устанавливаем текущее время как время обновления
			},
		},
	)
	if err != nil {
		log.Printf("Ошибка обновления заметки: %v\n", err)
		return
	}
	log.Printf("Обновлено заметок: %d\n", updateResult.ModifiedCount)

	// Операция чтения одного документа (Read One)
	// --------------------------------

	// Получаем обновленную заметку для проверки изменений
	// FindOne находит и возвращает первый документ, соответствующий фильтру
	var updatedNote model.Note
	err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&updatedNote)
	if err != nil {
		log.Printf("Ошибка получения обновленной заметки: %v\n", err)
		return
	}

	// Проверяем, есть ли значение в поле UpdatedAt
	updatedAtStr := "не обновлялась"
	if updatedNote.UpdatedAt != nil {
		updatedAtStr = updatedNote.UpdatedAt.Format(time.RFC3339)
	}

	log.Printf("Обновленная заметка - ID: %s, Заголовок: %s, Содержание: %s, Создана: %v, Обновлена: %s\n",
		updatedNote.ID.Hex(), updatedNote.Title, updatedNote.Body,
		updatedNote.CreatedAt.Format(time.RFC3339), updatedAtStr)

	// Операция удаления (Delete)
	// --------------------------------

	// Создаем еще одну заметку специально для демонстрации удаления
	noteForDelete := model.Note{
		Title:     gofakeit.City(),
		Body:      gofakeit.Address().Street,
		CreatedAt: time.Now(),
	}

	// Вставляем заметку, предназначенную для удаления
	insertResult, err = collection.InsertOne(ctx, noteForDelete)
	if err != nil {
		log.Printf("Ошибка вставки заметки для удаления: %v\n", err)
		return
	}

	// Получаем ID вставленной заметки
	deleteID := insertResult.InsertedID.(primitive.ObjectID)
	log.Printf("Добавлена заметка для удаления, ID: %s\n", deleteID.Hex())

	// Удаляем заметку по её ID
	// DeleteOne удаляет первый документ, соответствующий фильтру
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": deleteID})
	if err != nil {
		log.Printf("Ошибка удаления заметки: %v\n", err)
		return
	}
	log.Printf("Удалено заметок: %d, ID удаленной заметки: %s\n",
		deleteResult.DeletedCount, deleteID.Hex())

	log.Println("Демонстрация операций с MongoDB успешно завершена!")
}
