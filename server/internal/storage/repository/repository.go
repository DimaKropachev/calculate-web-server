package repository

import (
	"fmt"
	"sync"

	"github.com/DimaKropachev/calculate-web-server/server/internal/entities"
	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
)

type ExpressionsStorage struct {
	Exprs map[string]*models.Expression
	mu    sync.RWMutex
	id    *entities.SafeId
}

func NewExpressinsStorage() *ExpressionsStorage {
	return &ExpressionsStorage{
		Exprs: make(map[string]*models.Expression),
		mu:    sync.RWMutex{},
		id:    entities.NewSafeId(),
	}
}

func (ec *ExpressionsStorage) Add(expression string) string {
	Expr := &models.Expression{
		ID:     ec.id.Get(),
		Value:  expression,
		Status: "",
	}

	ec.mu.Lock()
	ec.Exprs[Expr.ID] = Expr
	ec.mu.Unlock()

	return Expr.ID
}

func (ec *ExpressionsStorage) Get(id string) (*models.Expression, error) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	if expr, ok := ec.Exprs[id]; ok {
		return expr, nil
	}
	return nil, fmt.Errorf("invalid ID expression")
}

func (ec *ExpressionsStorage) GetAll() []*models.Expression {
	var result []*models.Expression

	ec.mu.RLock()
	defer ec.mu.RUnlock()

	for _, expr := range ec.Exprs {
		result = append(result, expr)
	}

	return result
}