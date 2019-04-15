package api

import "fmt"

type Pagination struct {
	limit  int
	page   int
	search string
}

func (p *Pagination) Offset() int {
	return (p.Page() - 1) * p.limit
}

func (p *Pagination) Limit() int {
	if p.limit == 0 {
		return 2147483647
	}
	return p.limit
}

func (p *Pagination) NoLimit() bool {
	return p.limit == 0
}

func (p *Pagination) Page() int {
	if p.page == 0 {
		return 1
	}
	return p.page
}

func (p *Pagination) Search() string {
	return fmt.Sprintf("%%%s%%", p.search)
}

func (p *Pagination) ShouldSearch() bool {
	return len(p.search) > 0
}

func NewPagination(limit int, page int, search string) *Pagination {
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	return &Pagination{
		limit,
		page,
		search,
	}
}
