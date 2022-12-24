package util

import (
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type WhereCondition map[string]func(db *gorm.DB, value any)

type Paginator struct {
	db              *gorm.DB
	whereConditions WhereCondition
	orders          map[string]string
	Page            int
	PerPage         int
	Total           int64
	LastPage        uint
}

type PaginationRequest struct {
	Page    int    `form:"page" mod:"default=1" validate:"omitempty,min=1"`
	PerPage int    `form:"per_page" mod:"default=10" validate:"omitempty,min=1"`
	OrderBy string `form:"order_by" validate:"omitempty"`
}

type PaginationResponse struct {
	Page     int   `json:"page,omitempty"`
	PerPage  int   `json:"per_page,omitempty"`
	LastPage uint  `json:"last_page,omitempty"`
	Total    int64 `json:"total,omitempty"`
}

func NewPaginator(
	db *gorm.DB,
	whereConditions WhereCondition,
	orders map[string]string,
) Paginator {
	paginator := Paginator{
		db:              db,
		whereConditions: whereConditions,
		orders:          orders,
	}
	paginator.SetPagination(1, 10)

	return paginator
}

func (p *Paginator) AddWhereConditions(filter reflect.Value) *Paginator {
	if filter.IsNil() {
		return p
	}

	v := filter.Elem()
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			key := v.Type().Field(i).Tag.Get("key")
			if whereCondition, ok := p.whereConditions[key]; ok {
				whereCondition(p.db, v.Field(i).Interface())
			}
		}
	}

	return p
}

func (p *Paginator) SetPagination(page int, perPage int) *Paginator {
	p.Page = page
	p.PerPage = perPage

	return p
}

func (p *Paginator) SetOrderBy(orderBy string) *Paginator {
	order := strings.Trim(orderBy, "-")
	order, ok := p.orders[order]
	if !ok {
		return p
	}
	if strings.HasPrefix(orderBy, "-") {
		order = order + " DESC"
	}
	p.db.Order(order)

	return p
}

func (p *Paginator) CalculateTotalAndLastPage() *Paginator {
	p.db.Session(&gorm.Session{}).Count(&p.Total)
	lastPage := math.Ceil(float64(p.Total) / float64(p.PerPage))
	p.LastPage = uint(lastPage)

	return p
}

func (p *Paginator) GetData(data any, query any, args ...any) *Paginator {
	p.db.Session(&gorm.Session{}).
		Offset((p.Page-1)*p.PerPage).
		Limit(p.PerPage).
		Select(query, args...).
		Find(data)

	return p
}
