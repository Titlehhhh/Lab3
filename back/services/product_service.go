package services

import (
	"backApp/models"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type ProductService struct {
	path  string
	mu    sync.Mutex
	cache []models.Product
}

func NewProductService(path string) (*ProductService, error) {
	products, err := loadProductsFromFile(path)
	if err != nil {
		return nil, err
	}

	return &ProductService{
		path:  path,
		cache: products,
	}, nil
}

func (s *ProductService) List() []models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]models.Product{}, s.cache...)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for index, product := range s.cache {
		if product.Id == id {
			return &s.cache[index], nil
		}
	}
	return nil, errors.New("product not found")
}

func (s *ProductService) Add(product models.Product) (models.Product, error) {
	if err := validateProduct(product); err != nil {
		return models.Product{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := loadProductsFromFile(s.path)
	if err != nil {
		return models.Product{}, err
	}

	product.Id = nextProductID(products)
	products = append(products, product)

	if err := saveProductsToFile(s.path, products); err != nil {
		return models.Product{}, err
	}

	s.cache = products
	return product, nil
}

func (s *ProductService) Update(id int, updates models.Product) (models.Product, error) {
	if id <= 0 {
		return models.Product{}, errors.New("invalid product id")
	}

	if err := validateProduct(updates); err != nil {
		return models.Product{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := loadProductsFromFile(s.path)
	if err != nil {
		return models.Product{}, err
	}

	for index, product := range products {
		if product.Id == id {
			product.Name = updates.Name
			product.Description = updates.Description
			product.Price = updates.Price
			product.Image = updates.Image
			products[index] = product

			if err := saveProductsToFile(s.path, products); err != nil {
				return models.Product{}, err
			}
			s.cache = products
			return product, nil
		}
	}

	return models.Product{}, errors.New("product not found")
}

func (s *ProductService) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid product id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := loadProductsFromFile(s.path)
	if err != nil {
		return err
	}

	for index, product := range products {
		if product.Id == id {
			products = append(products[:index], products[index+1:]...)
			if err := saveProductsToFile(s.path, products); err != nil {
				return err
			}
			s.cache = products
			return nil
		}
	}
	return errors.New("product not found")
}

func validateProduct(product models.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(product.Description) == "" {
		return errors.New("description is required")
	}
	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if strings.TrimSpace(product.Image) == "" {
		return errors.New("image is required")
	}
	return nil
}

func loadProductsFromFile(path string) ([]models.Product, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var products []models.Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func saveProductsToFile(path string, products []models.Product) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(products)
}

func nextProductID(products []models.Product) int {
	maxID := 0
	for _, product := range products {
		if product.Id > maxID {
			maxID = product.Id
		}
	}
	return maxID + 1
}
