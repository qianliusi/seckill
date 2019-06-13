package soa

import (
	"sync"
)

type ProductCountMgr struct {
	productCount map[int]int
	sync.RWMutex
}

func NewProductCountMgr() (productMgr *ProductCountMgr) {
	productMgr = &ProductCountMgr{
		productCount: make(map[int]int, 128),
	}
	return
}

func (p *ProductCountMgr) Count(productId int) (count int) {
	p.RLock()
	defer p.RUnlock()
	count = p.productCount[productId]
	return

}

func (p *ProductCountMgr) Add(productId, count int) {
	p.Lock()
	defer p.Unlock()
	cur, ok := p.productCount[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.productCount[productId] = cur
}
