package product

type ProductService struct {
	ProductRepository *ProductRepository
}

func NewProductService(productRepository *ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}

func (ps *ProductService) GetAll() string {
	// Implement logic to get all product from the repository
	return "ok"
}
