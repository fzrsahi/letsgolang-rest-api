package product

import (
	"context"
	"database/sql"
	"task-one/category"
	"task-one/exception"
	"task-one/helpers"
	"task-one/product/dto"
	"task-one/product/model"
	"task-one/product/response"
)

type ProductService interface {
	Create(ctx context.Context, request *dto.ProductCreateDto) response.ProductResponse
	Update(ctx context.Context, request *dto.ProductUpdateDto) response.ProductResponse
	Delete(ctx context.Context, productId int)
	FindById(ctx context.Context, productId int) response.ProductResponse
	FindAll(ctx context.Context) []response.ProductResponse
}

type ProductServiceImpl struct {
	Repository         ProductRepository
	DB                 *sql.DB
	CategoryRepository category.CategoryRepository
}

func NewProductService(repository ProductRepository, DB *sql.DB, categoryRepository category.CategoryRepository) ProductService {
	return &ProductServiceImpl{Repository: repository, DB: DB, CategoryRepository: categoryRepository}
}

func (service *ProductServiceImpl) Create(ctx context.Context, request *dto.ProductCreateDto) response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)

	defer helpers.CommitOrRollback(tx)

	product := model.Product{
		Name:       request.Name,
		CategoryId: request.CategoryId,
	}

	product = service.Repository.Save(ctx, tx, product)
	return model.ToProductResponse(product)

}

func (service *ProductServiceImpl) Update(ctx context.Context, request *dto.ProductUpdateDto) response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	var product model.Product

	product, err = service.Repository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if request.CategoryId != 0 {
		_, err := service.CategoryRepository.FindById(ctx, tx, request.CategoryId)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}
		product = model.Product{
			Name:       request.Name,
			CategoryId: request.Id,
		}
	} else {
		product = model.Product{
			Name:       request.Name,
			CategoryId: request.Id,
		}
	}

	product = service.Repository.Update(ctx, tx, product)
	return model.ToProductResponse(product)

}

func (service *ProductServiceImpl) Delete(ctx context.Context, productId int) {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	product, err := service.Repository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.Repository.Delete(ctx, tx, product.Id)
}

func (service *ProductServiceImpl) FindById(ctx context.Context, productId int) response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	product, err := service.Repository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return model.ToProductResponse(product)
}

func (service *ProductServiceImpl) FindAll(ctx context.Context) []response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	products := service.Repository.FindAll(ctx, tx)
	return model.ToProductResponses(products)
}
