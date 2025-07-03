package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	ufoV1 "github.com/alexis871aa/microservices-course-examples/week_1/grpc/pkg/proto/ufo/v1"
)

const serverAddress = "localhost:50051"

// createSighting создает новое наблюдение НЛО с рандомными данными
func createSighting(ctx context.Context, client ufoV1.UFOServiceClient) (string, error) {
	// Генерируем случайные данные с помощью gofakeit
	observedAt := gofakeit.DateRange(
		time.Now().AddDate(-3, 0, 0), // за последние 3 года
		time.Now(),
	)
	location := gofakeit.City() + ", " + gofakeit.StreetName()
	description := gofakeit.Sentence(gofakeit.Number(5, 15))

	// Создаем базовую информацию о наблюдении
	info := &ufoV1.SightingInfo{
		ObservedAt:  timestamppb.New(observedAt),
		Location:    location,
		Description: description,
	}

	// Иногда добавляем дополнительные поля (с вероятностью 70%)
	if gofakeit.Bool() {
		info.Color = wrapperspb.String(gofakeit.Color())
	}

	if gofakeit.Bool() {
		info.Sound = wrapperspb.Bool(gofakeit.Bool())
	}

	if gofakeit.Bool() {
		info.DurationSeconds = wrapperspb.Int32(gofakeit.Int32())
	}

	// Вызываем gRPC метод Create
	resp, err := client.Create(ctx, &ufoV1.CreateRequest{Info: info})
	if err != nil {
		return "", err
	}

	return resp.Uuid, nil
}

// getSighting получает информацию о наблюдении по UUID
func getSighting(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) (*ufoV1.Sighting, error) {
	resp, err := client.Get(ctx, &ufoV1.GetRequest{Uuid: uuid})
	if err != nil {
		return nil, err
	}

	return resp.Sighting, nil
}

// updateSighting обновляет наблюдение НЛО
func updateSighting(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) error {
	// Генерируем рандомные данные для обновления
	updateInfo := &ufoV1.SightingUpdateInfo{}

	// Обновляем часть полей случайным образом
	if gofakeit.Bool() {
		updateInfo.ObservedAt = timestamppb.New(gofakeit.DateRange(
			time.Now().AddDate(-3, 0, 0),
			time.Now(),
		))
	}

	if gofakeit.Bool() {
		location := gofakeit.City() + ", " + gofakeit.StreetName()
		updateInfo.Location = wrapperspb.String(location)
	}

	if gofakeit.Bool() {
		description := gofakeit.Sentence(gofakeit.Number(5, 15))
		updateInfo.Description = wrapperspb.String(description)
	}

	if gofakeit.Bool() {
		updateInfo.Color = wrapperspb.String(gofakeit.Color())
	}

	if gofakeit.Bool() {
		updateInfo.Sound = wrapperspb.Bool(gofakeit.Bool())
	}

	if gofakeit.Bool() {
		updateInfo.DurationSeconds = wrapperspb.Int32(gofakeit.Int32())
	}

	// Вызываем gRPC метод Update
	_, err := client.Update(ctx, &ufoV1.UpdateRequest{
		Uuid:       uuid,
		UpdateInfo: updateInfo,
	})
	if err != nil {
		return err
	}

	return nil
}

// deleteSighting удаляет наблюдение НЛО
func deleteSighting(ctx context.Context, client ufoV1.UFOServiceClient, uuid string) error {
	_, err := client.Delete(ctx, &ufoV1.DeleteRequest{Uuid: uuid})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// создаём контекст, потому что он понадобится
	ctx := context.Background()

	// создаём коннекшн, то есть подключение. Передаём адрес, а вторым параметром идёт опции (если посмотреть на этот метод).
	conn, err := grpc.NewClient(
		serverAddress,
		// нас интересует одна опция WithTransportCredentials, передаём в неё insecure.NewCredentials() из пакета insecure
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	// выше получается мы создаём коннекшн, но не используем транспортную безопасность (TSL, известный как Secure Socket Layer (SSL) протокол транспортной безопасности). Чтобы себе жизнь не усложнять мы делаем как описано выше, явно показывая grpc что мы не собираемся шифровать, что это будет несекьюрное подключение
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	// закрываем за собой коннекшн, когда клиент отработает
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент, берём из сгенерированного файла, передаём туда conn наш коннекшн
	client := ufoV1.NewUFOServiceClient(conn) // после этого у клиента теперь можно вызывать методы

	log.Println("=== Тестирование API для работы с наблюдениями НЛО ===")
	log.Println()

	// 1. Создаем несколько наблюдений
	log.Println("🛸 Создание наблюдений НЛО")
	log.Println("===========================")
	uuid, err := createSighting(ctx, client)
	if err != nil {
		log.Printf("Ошибка при создании наблюдения: %v\n", err)
		return
	}

	// Выводим информацию о созданном наблюдении
	log.Printf("Создано наблюдение НЛО: UUID=%s\n", uuid)

	// 2. Получаем информацию о наблюдении
	log.Println("🔍 Получение информации о наблюдении")
	log.Println("==================================")
	sighting, err := getSighting(ctx, client, uuid)
	if err != nil {
		log.Printf("Ошибка при получении наблюдения: %v\n", err)
		return
	}

	// Выводим информацию о полученном наблюдении
	log.Printf("Получено наблюдение НЛО: UUID=%s", uuid)
	log.Printf("%v\n", sighting)

	// 3. Обновляем наблюдение
	log.Println("✏️ Обновление наблюдение")
	log.Println("=======================")

	err = updateSighting(ctx, client, uuid)
	if err != nil {
		log.Printf("Ошибка при обновлении наблюдения: %v", err)
		return
	}

	// 4. Проверяем обновленное наблюдение
	log.Println("🔍 Проверка обновленного наблюдения")
	log.Println("=================================")
	updatedSighting, err := getSighting(ctx, client, uuid)
	if err != nil {
		log.Printf("Ошибка при получении обновленного наблюдения: %v", err)
		return
	}

	// Выводим информацию об обновленном наблюдении
	log.Printf("Получено наблюдение НЛО: UUID=%s", uuid)
	log.Printf("%v\n", updatedSighting)

	// 6. Удаляем наблюдение
	err = deleteSighting(ctx, client, uuid)
	if err != nil {
		log.Printf("Ошибка при удалении наблюдения: %v", err)
	}

	// 7. Проверяем удаленное наблюдение
	log.Println("🔍 Проверка удаленного наблюдения")
	log.Println("=================================")
	deletedSighting, err := getSighting(ctx, client, uuid)
	if err != nil {
		log.Printf("Ошибка при получении удаленного наблюдения: %v", err)
		return
	}

	// Выводим информацию об удаленном наблюдении
	log.Printf("Получено удаленное наблюдение НЛО: UUID=%s", uuid)
	log.Printf("%v\n", deletedSighting)

	log.Println("Тестирование завершено!")
}
