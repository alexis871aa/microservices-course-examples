package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/shared/pkg/proto/ufo/v1"
)

var _ = Describe("UFOService", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		ufoClient ufoV1.UFOServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		ufoClient = ufoV1.NewUFOServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearSightingsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции sightings")

		cancel()
	})

	Describe("Create", func() {
		It("должен успешно создавать новое наблюдение НЛО", func() {
			info := env.GetTestSightingInfo()

			resp, err := ufoClient.Create(ctx, &ufoV1.CreateRequest{
				Info: info,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetUuid()).ToNot(BeEmpty())
			Expect(resp.GetUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})
	})

	Describe("Get", func() {
		var sightingUUID string

		BeforeEach(func() {
			// Вставляем тестовое наблюдение
			var err error
			sightingUUID, err = env.InsertTestSighting(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового наблюдения в MongoDB")
		})

		It("должен успешно возвращать наблюдение по UUID", func() {
			resp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: sightingUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetSighting()).ToNot(BeNil())
			Expect(resp.GetSighting().Uuid).To(Equal(sightingUUID))
			Expect(resp.GetSighting().GetInfo()).ToNot(BeNil())
			Expect(resp.GetSighting().GetInfo().Location).ToNot(BeEmpty())
			Expect(resp.GetSighting().GetInfo().Description).ToNot(BeEmpty())
			Expect(resp.GetSighting().GetCreatedAt()).ToNot(BeNil())
		})
	})

	Describe("Update", func() {
		var sightingUUID string

		BeforeEach(func() {
			// Вставляем тестовое наблюдение
			var err error
			sightingUUID, err = env.InsertTestSighting(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового наблюдения в MongoDB")
		})

		It("должен успешно обновлять наблюдение", func() {
			updateInfo := env.GetUpdatedSightingInfo()

			_, err := ufoClient.Update(ctx, &ufoV1.UpdateRequest{
				Uuid:       sightingUUID,
				UpdateInfo: updateInfo,
			})

			Expect(err).ToNot(HaveOccurred())

			// Проверяем, что наблюдение действительно обновилось
			resp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: sightingUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetSighting().GetInfo().Location).To(Equal(updateInfo.Location.GetValue()))
			Expect(resp.GetSighting().GetInfo().Description).To(Equal(updateInfo.Description.GetValue()))
			Expect(resp.GetSighting().GetInfo().Color.GetValue()).To(Equal(updateInfo.Color.GetValue()))
			Expect(resp.GetSighting().GetInfo().DurationSeconds.GetValue()).To(Equal(updateInfo.DurationSeconds.GetValue()))
			Expect(resp.GetSighting().GetUpdatedAt()).ToNot(BeNil())
		})
	})

	Describe("Delete", func() {
		var sightingUUID string

		BeforeEach(func() {
			// Вставляем тестовое наблюдение
			var err error
			sightingUUID, err = env.InsertTestSighting(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового наблюдения в MongoDB")
		})

		It("должен успешно выполнять мягкое удаление наблюдения", func() {
			_, err := ufoClient.Delete(ctx, &ufoV1.DeleteRequest{
				Uuid: sightingUUID,
			})

			Expect(err).ToNot(HaveOccurred())

			// Проверяем, что наблюдение помечено как удаленное
			resp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: sightingUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetSighting().GetDeletedAt()).ToNot(BeNil())
		})
	})

	Describe("Полный жизненный цикл", func() {
		It("должен поддерживать полный CRUD цикл", func() {
			// 1. Создаем наблюдение
			info := env.GetTestSightingInfo()
			createResp, err := ufoClient.Create(ctx, &ufoV1.CreateRequest{
				Info: info,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(createResp.GetUuid()).ToNot(BeEmpty())
			uuid := createResp.GetUuid()

			// 2. Получаем созданное наблюдение
			getResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getResp.GetSighting().Uuid).To(Equal(uuid))
			Expect(getResp.GetSighting().GetInfo().Location).To(Equal(info.Location))
			Expect(getResp.GetSighting().GetInfo().Description).To(Equal(info.Description))

			// 3. Обновляем наблюдение
			updateInfo := env.GetUpdatedSightingInfo()
			_, err = ufoClient.Update(ctx, &ufoV1.UpdateRequest{
				Uuid:       uuid,
				UpdateInfo: updateInfo,
			})

			Expect(err).ToNot(HaveOccurred())

			// 4. Проверяем обновление
			getUpdatedResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getUpdatedResp.GetSighting().GetInfo().Location).To(Equal(updateInfo.Location.GetValue()))
			Expect(getUpdatedResp.GetSighting().GetInfo().Description).To(Equal(updateInfo.Description.GetValue()))

			// 5. Удаляем наблюдение
			_, err = ufoClient.Delete(ctx, &ufoV1.DeleteRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())

			// 6. Проверяем, что наблюдение помечено как удаленное
			getDeletedResp, err := ufoClient.Get(ctx, &ufoV1.GetRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getDeletedResp.GetSighting().GetDeletedAt()).ToNot(BeNil())
		})
	})
})
