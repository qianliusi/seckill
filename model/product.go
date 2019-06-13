package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ProductDao struct {
}

type Product struct {
	Id        int    `db:"id" json:"id" form:"id"`
	Name      string `db:"name" json:"name" form:"name"`
	ShortName string `db:"shortName" json:"shortName" form:"shortName"`
	Area      string `db:"area" json:"area" form:"area"`
	Total     int    `db:"total" json:"total" form:"total"`
}

type ProductListParam struct {
	PageParam
	Id   int           `form:"id"`
	Ids  []interface{} `form:"ids"`
	Name string        `form:"name"`
}

func NewProductDao() *ProductDao {
	return &ProductDao{}
}

func (p *ProductDao) QueryProductDetail(productId int) (pro Product, err error) {
	sql := "select id, name,shortName,area, total from product where id=?"
	err = Db.Get(&pro, sql, productId)
	if err != nil {
		logs.Warn("exec sql err:%v sql:%v", err, sql)
		return
	}
	return
}

func (p *ProductDao) QueryProductList(param *ProductListParam) (list []Product, total int, err error) {
	hasLike := param.Name != ""
	hasIds := len(param.Ids) > 0
	sqlc := "select count(*) from product where 1=1 "
	sqlq := "select id, name,shortName,area, total from product where 1=1 "
	var argsT []interface{}
	var args []interface{}
	if hasIds {
		sqlq += "and id in (?) "
		sqlq, _, err = sqlx.In(sqlq, param.Ids)
		args = append(args, param.Ids...)
	}
	if hasLike {
		sqlc += "and name like ? "
		sqlq += "and name like ? "
		argsT = append(argsT, "%"+param.Name+"%")
		args = append(args, "%"+param.Name+"%")
	}
	err = Db.Get(&total, sqlc, argsT...)
	if err != nil {
		logs.Error("exec sql err:%v sql:%v", err, sqlc)
		return
	}
	if total == 0 {
		return
	}
	sqlq += "order by id limit ?,?"
	args = append(args, param.Start, param.Limit)
	err = Db.Select(&list, sqlq, args...)
	if err != nil {
		logs.Error("exec sql err:%v sql:%v", err, sqlq)
		return
	}
	return
}

func (p *ProductDao) AddProduct(product *Product) (err error) {
	sql := "insert into product(name,shortName,area, total)values(?,?,?,?)"
	_, err = Db.Exec(sql, product.Name, product.ShortName, product.Area, product.Total)
	if err != nil {
		logs.Warn("CreateProduct failed, err:%v sql:%v", err, sql)
		return
	}
	return
}
func (p *ProductDao) EditProduct(product *Product) (err error) {
	sql := "update product set name=?,shortName=?,area=?,total=? where id=?"
	_, err = Db.Exec(sql, product.Name, product.ShortName, product.Area, product.Total, product.Id)
	if err != nil {
		logs.Warn("EditProduct failed, err:%v sql:%v", err, sql)
		return
	}
	return
}
