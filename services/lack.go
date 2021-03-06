package services

import (
	"github.com/astaxie/beego/orm"
	"github.com/scmo/apayment-backend/models"
)

func CreateLack(l *models.Lack) error {
	o := orm.NewOrm()
	_, err := o.Insert(l)
	return err
}

func CreateMultiLacks(lacks []models.Lack) error {
	o := orm.NewOrm()
	_, err := o.InsertMulti(50, lacks)
	return err
}

func GetAllLacks() []*models.Lack {
	o := orm.NewOrm()
	var lacks []*models.Lack
	o.QueryTable(new(models.Lack)).All(&lacks)
	return lacks
}

func CountLacks() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(new(models.Lack)).Count() // SELECT COUNT(*) FROM USE
	return cnt, err
}
