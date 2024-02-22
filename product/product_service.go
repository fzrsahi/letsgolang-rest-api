package product

import (
	"context"
	"database/sql"
	"sync"
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
	Wg                 *sync.WaitGroup
}

func NewProductService(repository ProductRepository, DB *sql.DB, categoryRepository category.CategoryRepository, wg *sync.WaitGroup) ProductService {
	return &ProductServiceImpl{Repository: repository, DB: DB, CategoryRepository: categoryRepository, Wg: wg}
}

func (service *ProductServiceImpl) Create(ctx context.Context, request *dto.ProductCreateDto) response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	_, err = service.CategoryRepository.FindById(ctx, tx, request.CategoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	productChannel := make(chan model.Product)
	defer close(productChannel)
	product := model.Product{
		Name:       request.Name,
		CategoryId: request.CategoryId,
	}

	service.Wg.Add(1)
	go func() {
		defer service.Wg.Done()
		product = service.Repository.Save(ctx, tx, product)
		productChannel <- product
	}()

	product = <-productChannel
	defer service.Wg.Wait()
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
			Id:         request.Id,
			Name:       request.Name,
			CategoryId: request.CategoryId,
		}
	} else {
		product = model.Product{
			Id:   request.Id,
			Name: request.Name,
		}
	}

	productChannel := make(chan model.Product)
	defer close(productChannel)
	service.Wg.Add(1)
	go func() {
		defer service.Wg.Done()
		product = service.Repository.Update(ctx, tx, product)
		productChannel <- product
	}()

	product = <-productChannel
	defer service.Wg.Wait()
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

	service.Wg.Add(1)
	go func() {
		defer service.Wg.Done()
		service.Repository.Delete(ctx, tx, product.Id)
	}()

	service.Wg.Wait()
}

func (service *ProductServiceImpl) FindById(ctx context.Context, productId int) response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	productChannel := make(chan struct {
		Product model.Product
		Error   error
	})
	service.Wg.Add(1)
	defer close(productChannel)

	go func() {
		defer service.Wg.Done()
		product, err := service.Repository.FindById(ctx, tx, productId)
		productChannel <- struct {
			Product model.Product
			Error   error
		}{product, err}
	}()

	productResult := <-productChannel
	if productResult.Error != nil {
		panic(exception.NewNotFoundError(productResult.Error.Error()))
	}

	defer service.Wg.Wait()
	return model.ToProductResponse(productResult.Product)
}

func (service *ProductServiceImpl) FindAll(ctx context.Context) []response.ProductResponse {
	tx, err := service.DB.Begin()
	helpers.PanicIfError(err)
	defer helpers.CommitOrRollback(tx)

	productsChannel := make(chan []model.Product)
	defer close(productsChannel)
	service.Wg.Add(1)

	go func() {
		defer service.Wg.Done()
		products := service.Repository.FindAll(ctx, tx)
		productsChannel <- products
	}()

	products := <-productsChannel
	defer service.Wg.Wait()
	return model.ToProductResponses(products)
}
