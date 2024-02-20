package category

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
	"task-one/category/model"
	"task-one/configs/database"
	"task-one/exception"
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
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestGetListCategory(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)

	req := httptest.NewRequest("GET", "http://localhost:3001/categories", nil)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)
}

func TestCreateCategory(t *testing.T) {
	db := database.ConnectToDbTest()
	router := setupRouter(db)
	truncateCategory(db)

	reqBody := strings.NewReader(`{"name" : "Website"}`)
	req := httptest.NewRequest("POST", "http://localhost:3001/categories", reqBody)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	res := recorder.Result()

	body, _ := ioutil.ReadAll(res.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "Website", responseBody["data"].(map[string]interface{})["name"])

}

func TestUpdateCategory(t *testing.T) {
	db := database.ConnectToDbTest()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()

	repository := NewCategoryRepository()
	category := repository.Save(context.Background(), tx, model.Category{
		Name: "Handphone",
	})
	tx.Commit()

	t.Run("Test Update Category Success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Not Handphone"}`)
		req := httptest.NewRequest("PATCH", "http://localhost:3001/categories/"+strconv.Itoa(category.Id), reqBody)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "Not Handphone", responseBody["data"].(map[string]interface{})["name"])

	})

	t.Run("Test Update Category Failed", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name" : "Not Handphone"}`)
		req := httptest.NewRequest("PATCH", "http://localhost:3001/categories/"+strconv.Itoa(404), reqBody)
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

func TestDeleteCategory(t *testing.T) {
	db := database.ConnectToDbTest()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()

	repository := NewCategoryRepository()
	category := repository.Save(context.Background(), tx, model.Category{
		Name: "Delete",
	})
	tx.Commit()

	t.Run("Test Delete Category Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://localhost:3001/categories/"+strconv.Itoa(category.Id), nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, res.StatusCode)

	})

	t.Run("Test Update Category Failed", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://localhost:3001/categories/"+strconv.Itoa(404), nil)
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

func TestGetCategoryById(t *testing.T) {
	db := database.ConnectToDbTest()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()

	repository := NewCategoryRepository()
	category := repository.Save(context.Background(), tx, model.Category{
		Name: "Delete",
	})
	tx.Commit()

	t.Run("Test Delete Category Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3001/categories/"+strconv.Itoa(category.Id), nil)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		res := recorder.Result()

		body, _ := ioutil.ReadAll(res.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, res.StatusCode)

	})

	t.Run("Test Update Category Failed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3001/categories/"+strconv.Itoa(404), nil)
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
