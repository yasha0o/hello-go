package controller

import (
	"net/http"

	"hello-go/internal/dto"
	"hello-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArchiveController struct {
	service *service.ArchiveService
}

func NewArchiveController(s *service.ArchiveService) *ArchiveController {
	return &ArchiveController{
		service: s,
	}
}

func (c *ArchiveController) Init(enginie *gin.Engine) {
	enginie.GET("/archive/documents", ErrorWrapper(c.GetDocuments))
	enginie.GET("/archive/documents/:id", ErrorWrapper(c.GetDocument))
	enginie.POST("/archive/documents/:id", ErrorWrapper(c.LoadDocument))
}

// GetDocuments получение списка документов с пагинацией
// @Summary Получить список документов
// @Description Возвращает список документов с поддержкой пагинации
// @Tags archive
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param size query int false "Размер страницы" default(10)
// @Success 200 {object} dto.RequestedDocumentsPage
// "Успешный ответ со списком документов"
// @Failure 400 {object} ApiError "Ошибка валидации"
// @Failure 500 {object} ApiError "Внутренняя ошибка сервера"
// @Router /archive/documents [get]
func (c *ArchiveController) GetDocuments(ctx *gin.Context) error {
	var req dto.PageParams

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return &ValidationError{Err: err}
	}
	archives, err := c.service.GetDocuments(ctx.Request.Context(), &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, archives)
	return nil
}

// GetDocument получение документа по ID
// @Summary Получить документ по ID
// @Description Возвращает документ по указанному ID
// @Tags archive
// @Accept json
// @Produce json
// @Param id path string true "ID документа" format(uuid)
// @Success 200 {object} dto.DocumentDto "Успешный ответ с документом"
// @Failure 404 {object} ApiError "Документ не найден"
// @Failure 500 {object} ApiError "Внутренняя ошибка сервера"
// @Router /archive/documents/{id} [get]
func (c *ArchiveController) GetDocument(ctx *gin.Context) error {
	var id uuid.UUID
	if err := ctx.ShouldBindUri(&id); err != nil {
		return &ValidationError{Err: err}
	}

	archives, err := c.service.GetDocument(ctx.Request.Context(), &id)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, archives)
	return nil
}

// LoadDocument загрузка документа по ID
// @Summary Загрузить документ по ID
// @Description Загружает документ с указанным ID из архивов
// @Tags archive
// @Accept json
// @Produce json
// @Param id path string true "ID документа" format(uuid)
// @Param sender header string true "Отправитель запроса"
// @Param region query string false "Регион"
// Enums(MOSCOW,SAINT_PETERSBURG,NOVOSIBIRSK,EKATERINBURG,KAZAN,
// NIZHNY_NOVGOROD,CHELYABINSK,SAMARA,OMSK,ROSTOV_ON_DON)
// @Success 200 {object} dto.CreatedRequest "Успешный ответ c идентификатором созданного запроса"
// @Failure 400 {object} ApiError "Ошибка валидации"
// @Failure 500 {object} ApiError "Внутренняя ошибка сервера" 
// @Router /archive/documents/{id} [post]
func (c *ArchiveController) LoadDocument(ctx *gin.Context) error {
	var req dto.LoadDocumentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return &ValidationError{Err: err}
	}
	if err := ctx.ShouldBindHeader(&req); err != nil {
		return &ValidationError{Err: err}
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return &ValidationError{Err: err}
	}

	response, err := c.service.LoadDocument(ctx.Request.Context(), &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, response)
	return nil
}
