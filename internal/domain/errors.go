package domain

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrInsufficientIngredients = errors.New("insufficient ingredients")
