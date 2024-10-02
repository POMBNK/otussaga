package order

import "github.com/POMBNK/orderservice/internal/entity"

func (r *CreateOrderJSONBody) ToOrder() (entity.Order, error) {
	return entity.Order{
		Goods: r.Goods,
	}, nil
}
