package category

import (
	"context"
	"database/sql"
	"task-one/category/dto"
	"task-one/category/model"
	"task-one/category/response"
	"task-one/exception"
	"task-one/helpers"
)

type CategoryService interface {
	Create(ctx context.Context, request *dto.CategoryCreateDto) response.CategoryResponse
	Update(ctx context.Context, request *dto.CategoryUpdateDto) response.CategoryResponse
	Delete(ctx context.Context, categoryId int)
	FindById(ctx context.Context, categoryId int) response.CategoryResponse
	FindAll(ctx context.Context) []response.CategoryResponse
}

type CategoryServiceImpl struct {
	Repository CategoryRepository
	DB         *sql.DB
}

func NewCategoryService(repository CategoryRepository, DB *sql.DB) CategoryService {
	return &CategoryServiceImpl{
		Repository: repository,
		DB:         DB,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request *dto.CategoryCreateDto) response.CategoryResponse {

	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)

	defer helpers.CommitOrRollback(tx)

	category := model.Category{
		Name: request.Name,
	}

	category = service.Repository.Save(ctx, tx, category)
	return helpers.ToCategoryResponse(category)

}

func (service *CategoryServiceImpl) Update(ctx context.Context, request *dto.CategoryUpdateDto) response.CategoryResponse {

	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	category, err := service.Repository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	category.Name = request.Name
	category = service.Repository.Update(ctx, tx, category)

	return helpers.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	category, err := service.Repository.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.Repository.Delete(ctx, tx, category.Id)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) response.CategoryResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	category, err := service.Repository.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helpers.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []response.CategoryResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	categories := service.Repository.FindAll(ctx, tx)

	return helpers.ToCategoryResponses(categories)
}
