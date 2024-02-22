package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"task-one/configs/redis"
	"task-one/helpers"
	"task-one/product/model"
)

type ProductRepository interface {
	Save(ctx context.Context, tx *sql.Tx, product model.Product) model.Product
	Update(ctx context.Context, tx *sql.Tx, product model.Product) model.Product
	Delete(ctx context.Context, tx *sql.Tx, productId int)
	FindAll(ctx context.Context, tx *sql.Tx) []model.Product
	FindById(ctx context.Context, tx *sql.Tx, productId int) (model.Product, error)
	UpdateCache(ctx context.Context, tx *sql.Tx)
}

type ProductRepositoryImpl struct {
	rdb *redis.RedisClient
}

func (p *ProductRepositoryImpl) UpdateCache(ctx context.Context, tx *sql.Tx) {
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

	key := "list:products"
	err = p.rdb.Set(ctx, key, products)
	helpers.PanicIfError(err)
}

func NewProductRepository(rdb *redis.RedisClient) ProductRepository {
	return &ProductRepositoryImpl{
		rdb: rdb,
	}
}

func (p *ProductRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product model.Product) model.Product {
	query := `
		WITH product AS (
			INSERT INTO product(name, category_id)
			VALUES($1, $2)
			RETURNING id, name ,category_id
		)
		SELECT product.id,product.name, category.name
		FROM product
		INNER JOIN category ON product.category_id = category.id
	`

	row := tx.QueryRowContext(ctx, query, product.Name, product.CategoryId)
	err := row.Scan(&product.Id, &product.Name, &product.CategoryName)
	helpers.PanicIfError(err)

	p.UpdateCache(ctx, tx)

	return product
}

func (p *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product model.Product) model.Product {
	var query string
	var err error

	if product.CategoryId != 0 && product.Name != "" {
		query = "UPDATE product SET name = $1, category_id = $2 WHERE id = $3 RETURNING id, name"
		err = tx.QueryRowContext(ctx, query, product.Name, product.CategoryId, product.Id).Scan(&product.Id, &product.Name)
	} else if product.CategoryId != 0 && product.Name == "" {
		query = "UPDATE product SET category_id = $1 WHERE id = $2 RETURNING id"
		err = tx.QueryRowContext(ctx, query, product.CategoryId, product.Id).Scan(&product.Id)
	} else if product.CategoryId == 0 && product.Name != "" {
		query = "UPDATE product SET name = $1 WHERE id = $2 RETURNING id, name"
		err = tx.QueryRowContext(ctx, query, product.Name, product.Id).Scan(&product.Id, &product.Name)
	} else {
		return product
	}

	helpers.PanicIfError(err)

	selectQuery := `
		SELECT p.id, p.name, c.name
		FROM product p
		INNER JOIN category c ON p.category_id = c.id
		WHERE p.id = $1
	`
	row := tx.QueryRowContext(ctx, selectQuery, product.Id)
	err = row.Scan(&product.Id, &product.Name, &product.CategoryName)
	helpers.PanicIfError(err)

	return product
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, productId int) {
	query := "DELETE FROM product where id = $1"
	_, err := tx.ExecContext(ctx, query, productId)
	helpers.PanicIfError(err)
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.Product {
	var products []model.Product

	key := "list:products"
	productsCache, err := p.rdb.Get(ctx, key)
	if err != nil {
		query := "SELECT product.id,product.name,category.name FROM product INNER JOIN category ON product.category_id = category.id"
		rows, err := tx.QueryContext(ctx, query)
		helpers.PanicIfError(err)
		defer rows.Close()

		for rows.Next() {
			product := model.Product{}
			err := rows.Scan(&product.Id, &product.Name, &product.CategoryName)
			helpers.PanicIfError(err)

			products = append(products, product)
		}

		err = p.rdb.Set(ctx, key, products)
		if err != nil {
			helpers.PanicIfError(err)
		}
		return products
	}

	err = json.Unmarshal([]byte(productsCache), &products)
	if err != nil {
		helpers.PanicIfError(err)
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
