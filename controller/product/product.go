package product

import (
	"seckill/controller/base"
	"seckill/model"
)

var productDao = model.ProDao

type ProductController struct {
	base.Controller
}

func (p *ProductController) Index() {
	p.HomeViewer("product/index.html", nil)
}

func (p *ProductController) List() {
	var total int
	var rows interface{}
	var err error
	defer func() {
		p.JsonPage(total, rows, err)
	}()
	param := &model.ProductListParam{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	param.Parse()
	rows, total, err = productDao.QueryProductList(param)
}

func (p *ProductController) Add() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.Product{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	err = productDao.AddProduct(param)
}

func (p *ProductController) Edit() {
	var data interface{}
	var err error
	defer func() {
		p.Json(data, err)
	}()
	param := &model.Product{}
	if err = p.ParseForm(param); err != nil {
		return
	}
	err = productDao.EditProduct(param)
}
