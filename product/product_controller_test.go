package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"task-one/category"
	"task-one/category/model"
	"task-one/configs/database"
	"task-one/exception"
	product_model "task-one/product/model"
	"testing"
)

func setupRouter(db *sql.DB) http.Handler {
	router := httprouter.New()
	router.PanicHandler = exception.ErrorHandler
	RegisterRoute(router, db)

	return router

}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
	db.Exec("TRUNCATE product")

}

func TestMain(m *testing.M) {
	m.Run()
}

func TestGetListProduct(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)

	req := httptest.NewRequest("GET", "http://localhost:3001/products", nil)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)
}

func TestCreateProduct(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)
	truncateCategory(db)
	tx, _ := db.Begin()

	categoryRepository := category.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, model.Category{
		Name: "Furniture",
	})
	tx.Commit()

	t.Run("Test Create Product Success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Table","category_id":` + strconv.Itoa(category.Id) + "}")
		req := httptest.NewRequest("POST", "http://localhost:3001/products", reqBody)
		req.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "Table", responseBody["data"].(map[string]interface{})["name"])
		assert.Equal(t, "Furniture", responseBody["data"].(map[string]interface{})["category_name"])
	})

	t.Run("Test Create Product Failed", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Table","category_id":404}`)
		req := httptest.NewRequest("POST", "http://localhost:3001/products", reqBody)
		req.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, res.StatusCode)
	})
}

func TestUpdateProduct(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)
	truncateCategory(db)
	tx, _ := db.Begin()

	ctx := context.Background()

	categoryRepository := category.NewCategoryRepository()
	category := categoryRepository.Save(ctx, tx, model.Category{
		Name: "Furniture",
	})
	categoryUpdate := categoryRepository.Save(ctx, tx, model.Category{
		Name: "Alat Rumah",
	})
	productRepository := NewProductRepository()
	product := productRepository.Save(ctx, tx, product_model.Product{
		Name:       "Table",
		CategoryId: category.Id,
	})
	tx.Commit()

	t.Run("Test Update Product Success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Meja","category_id":` + strconv.Itoa(categoryUpdate.Id) + "}")
		req := httptest.NewRequest("PATCH", "http://localhost:3001/products/"+strconv.Itoa(product.Id), reqBody)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "Meja", responseBody["data"].(map[string]interface{})["name"])
		assert.Equal(t, "Alat Rumah", responseBody["data"].(map[string]interface{})["category_name"])
	})

	t.Run("Test Update Product Failed", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Table","category_id":404}`)
		req := httptest.NewRequest("PATCH", "http://localhost:3001/products/"+strconv.Itoa(product.Id), reqBody)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, res.StatusCode)
	})
}

func TestGetProductById(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)
	truncateCategory(db)
	tx, _ := db.Begin()

	ctx := context.Background()

	categoryRepository := category.NewCategoryRepository()
	category := categoryRepository.Save(ctx, tx, model.Category{
		Name: "Furniture",
	})
	productRepository := NewProductRepository()
	product := productRepository.Save(ctx, tx, product_model.Product{
		Name:       "Table",
		CategoryId: category.Id,
	})
	tx.Commit()

	t.Run("Test Get Product By Id Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3001/products/"+strconv.Itoa(product.Id), nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "Table", responseBody["data"].(map[string]interface{})["name"])
		assert.Equal(t, "Furniture", responseBody["data"].(map[string]interface{})["category_name"])
	})

	t.Run("Test Get Product By Id Failed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3001/products/404", nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, res.StatusCode)
	})
}

func TestDeleteProduct(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)
	truncateCategory(db)
	tx, _ := db.Begin()

	ctx := context.Background()

	categoryRepository := category.NewCategoryRepository()
	category := categoryRepository.Save(ctx, tx, model.Category{
		Name: "Furniture",
	})
	productRepository := NewProductRepository()
	product := productRepository.Save(ctx, tx, product_model.Product{
		Name:       "Table",
		CategoryId: category.Id,
	})
	tx.Commit()

	t.Run("Test Delete Product Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://localhost:3001/products/"+strconv.Itoa(product.Id), nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("Test Delete Product Failed", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://localhost:3001/products/404", nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, res.StatusCode)
	})
}
