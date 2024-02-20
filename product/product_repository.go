package product

import (
	"context"
	"database/sql"
	"errors"
	"task-one/helpers"
	"task-one/product/model"
)

type ProductRepository interface {
	Save(ctx context.Context, tx *sql.Tx, product model.Product) model.Product
	Update(ctx context.Context, tx *sql.Tx, product model.Product) model.Product
	Delete(ctx context.Context, tx *sql.Tx, productId int)
	FindAll(ctx context.Context, tx *sql.Tx) []model.Product
	FindById(ctx context.Context, tx *sql.Tx, productId int) (model.Product, error)
}

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (p *ProductRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product model.Product) model.Product {
	query := "INSERT INTO category(name,category_id) values($1,$2) RETURNING id"
	row := tx.QueryRowContext(ctx, query, product.Name, product.CategoryId)
	err := row.Scan(&product.Id)
	helpers.PanicIfError(err)

	return product

}

func (p *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product model.Product) model.Product {
	query := "UPDATE product set name = $1 where id = $2"
	_, err := tx.ExecContext(ctx, query, product.Name, product.Id)
	helpers.PanicIfError(err)

	return product
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, productId int) {
	query := "DELETE FROM category where id = $1"
	_, err := tx.ExecContext(ctx, query, productId)
	helpers.PanicIfError(err)
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.Product {
	query := "SELECT product.id,product.name,category.name FROM product INNER JOIN category ON product.category_id = category.id"
	rows, err := tx.QueryContext(ctx, query)
	helpers.PanicIfError(err)
	defer rows.Close()
	var products []model.Product

	for rows.Next() {
		product := model.Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.CategoryName)
		helpers.PanicIfError(err)

		products = append(products, product)
	}

	return products
}

func (p *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, productId int) (model.Product, error) {
	query := "SELECT product.id,product.name,category.name FROM product INNER JOIN category ON product.category_id = category.id WHERE product.id = $1"
	rows, err := tx.QueryContext(ctx, query, productId)
	helpers.PanicIfError(err)
	defer rows.Close()

	product := model.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.Name, &product.CategoryName)
		helpers.PanicIfError(err)
		return product, nil
	} else {
		return product, errors.New("product Not Found")
	}
}
