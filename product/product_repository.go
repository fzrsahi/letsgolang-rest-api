package product

type ProductRepository struct {
	Id string
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (pr *ProductRepository) GetAll() string {
	// Implement logic to fetch product from the database
	return "Get All Products"
}
