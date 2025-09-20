package orders

import "errors"

var (
    ErrEmptyItems        = errors.New("items cannot be empty")
    ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
    ErrProductNotFound   = errors.New("product not found")
    ErrInsufficientStock = errors.New("insufficient stock for product")
)
