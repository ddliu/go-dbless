// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"math"
)

type Pagination struct {
	PageSize  uint `json:"page_size"`
	Page      uint `json:"page"`
	PageTotal uint `json:"page_total"`
	Total     uint `json:"total"`
}

func (p *Pagination) GetOffsetLimit() (uint, uint) {
	offset := p.PageSize * (p.Page - 1)
	limit := p.PageSize

	return offset, limit
}

func (p *Pagination) Valid() {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}

func (p *Pagination) SetTotal(total uint) {
	p.Valid()
	p.Total = total
	p.PageTotal = uint(math.Ceil(float64(p.Total) / float64(p.PageSize)))
}

func NewPagination(pageSize, page uint) *Pagination {
	p := &Pagination{
		PageSize: pageSize,
		Page:     page,
	}

	p.Valid()

	return p
}
