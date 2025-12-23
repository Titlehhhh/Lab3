package services

import (
	"backApp/models"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

type CartService struct {
	path  string
	mu    sync.Mutex
	carts map[string][]models.CartItem
}

func NewCartService(path string) (*CartService, error) {
	carts, err := loadCarts(path)
	if err != nil {
		return nil, err
	}

	return &CartService{
		path:  path,
		carts: carts,
	}, nil
}

func (s *CartService) GetCart(userID int) []models.CartItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := s.carts[strconv.Itoa(userID)]
	return append([]models.CartItem{}, items...)
}

func (s *CartService) AddItem(userID int, productID int, quantity int) error {
	if userID <= 0 || productID <= 0 {
		return errors.New("invalid ids")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := strconv.Itoa(userID)
	items := s.carts[key]
	for index, item := range items {
		if item.ProductID == productID {
			items[index].Quantity += quantity
			s.carts[key] = items
			return saveCarts(s.path, s.carts)
		}
	}

	items = append(items, models.CartItem{ProductID: productID, Quantity: quantity})
	s.carts[key] = items
	return saveCarts(s.path, s.carts)
}

func (s *CartService) RemoveItem(userID int, productID int) error {
	if userID <= 0 || productID <= 0 {
		return errors.New("invalid ids")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := strconv.Itoa(userID)
	items := s.carts[key]
	for index, item := range items {
		if item.ProductID == productID {
			items = append(items[:index], items[index+1:]...)
			if len(items) == 0 {
				delete(s.carts, key)
			} else {
				s.carts[key] = items
			}
			return saveCarts(s.path, s.carts)
		}
	}
	return errors.New("product not found in cart")
}

func loadCarts(path string) (map[string][]models.CartItem, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]models.CartItem{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var carts map[string][]models.CartItem
	if err := json.NewDecoder(file).Decode(&carts); err != nil {
		return nil, err
	}
	if carts == nil {
		carts = map[string][]models.CartItem{}
	}
	return carts, nil
}

func saveCarts(path string, carts map[string][]models.CartItem) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(carts)
}
