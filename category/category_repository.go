package category

import (
	"context"
	"database/sql"
	"errors"
	"task-one/category/model"
	"task-one/helpers"
)

type CategoryRepository interface {
	Save(ctx context.Context, tx *sql.Tx, category model.Category) model.Category
	Update(ctx context.Context, tx *sql.Tx, category model.Category) model.Category
	Delete(ctx context.Context, tx *sql.Tx, categoryId int)
	FindAll(ctx context.Context, tx *sql.Tx) []model.Category
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (model.Category, error)
}

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category model.Category) model.Category {
	query := "INSERT INTO category(name) values ($1) RETURNING id"
	row := tx.QueryRowContext(ctx, query, category.Name)
	err := row.Scan(&category.Id)
	helpers.PanicIfError(err)

	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category model.Category) model.Category {
	query := "UPDATE category set name = $1 where id = $2"
	_, err := tx.ExecContext(ctx, query, category.Name, category.Id)
	helpers.PanicIfError(err)

	return category

}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	query := "DELETE FROM category where id = $1"
	_, err := tx.ExecContext(ctx, query, categoryId)
	helpers.PanicIfError(err)
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []model.Category {
	query := "SELECT id,name FROM category"
	rows, err := tx.QueryContext(ctx, query)
	helpers.PanicIfError(err)
	defer rows.Close()
	var categories []model.Category

	for rows.Next() {
		category := model.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		helpers.PanicIfError(err)

		categories = append(categories, category)

	}

	return categories
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (model.Category, error) {
	SQL := "SELECT id, name FROM category WHERE id = $1"
	rows, err := tx.QueryContext(ctx, SQL, categoryId)
	helpers.PanicIfError(err)
	defer rows.Close()

	category := model.Category{}
	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)
		helpers.PanicIfError(err)
		return category, nil
	} else {
		return category, errors.New("category Not Found")
	}
}
