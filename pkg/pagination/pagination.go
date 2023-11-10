package pagination

import (
	"math"

	"gorm.io/gorm"
)

type WhereConditionMap map[string]func(db *gorm.DB, value any)

type Paginator struct {
	db                *gorm.DB
	whereConditionMap WhereConditionMap
	orders            map[string]string
	page              int
	perPage           int
}

type PaginationRequest struct {
	Page    int    `form:"page" validate:"min=1"`
	PerPage int    `form:"per_page" validate:"min=10"`
	OrderBy string `form:"order_by" validate:"omitempty"`
}

type PaginationResponse struct {
	LastPage int   `json:"last_page,omitempty"`
	Total    int64 `json:"total,omitempty"`
	Data     any   `json:"data"`
}

func NewPaginator(
	db *gorm.DB,
	whereConditionMap WhereConditionMap,
	orders map[string]string,
	page int,
	perPage int,
) Paginator {
	paginator := Paginator{
		db:                db,
		whereConditionMap: whereConditionMap,
		orders:            orders,
		page:              page,
		perPage:           perPage,
	}

	return paginator
}

func (p *Paginator) AddWhereConditions(wm map[string]any) *Paginator {
	for column, value := range wm {
		if whereFunc, ok := p.whereConditionMap[column]; ok {
			whereFunc(p.db, value)
		}
	}

	return p
}

func (p *Paginator) OrderBy(orderBy string) *Paginator {
	order, ok := p.orders[orderBy]
	if !ok {
		return p
	}
	p.db.Order(order)

	return p
}

func (p *Paginator) GetTotalAndLastPage() (int64, int) {
	var total int64
	p.db.Session(&gorm.Session{}).Count(&total)

	return total, int(math.Ceil(float64(total) / float64(p.perPage)))
}

func (p *Paginator) Execute(data any, query any, args ...any) *gorm.DB {
	return p.db.Session(&gorm.Session{}).
		Offset((p.page-1)*p.perPage).
		Limit(p.perPage).
		Select(query, args...).
		Find(data)
}

func (p *Paginator) CloneDB() *gorm.DB {
	return p.db.Session(&gorm.Session{})
}
