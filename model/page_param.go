package model

type DataParam struct {
	Data string `json:"data" form:"data"`
}

type PageParam struct {
	Page    int    `form:"page"`
	Rows    int    `form:"rows"`
	Order   string `form:"order"`
	Sort    string `form:"sort"`
	Start   int    `form:"-"`
	Limit   int    `form:"-"`
	OrderBy string `form:"-"`
}

func (p *PageParam) Parse() {
	if p.Rows != 0 && p.Page != 0 {
		p.Start = p.Rows * (p.Page - 1)
		p.Limit = p.Rows
	} else {
		p.Start = 0
		p.Limit = 0
	}
	if p.Sort != "" {
		p.OrderBy += p.Sort
		if p.Order != "" {
			p.OrderBy += " " + p.Order
		}
	}
}
