package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/scmo/apayment-backend/services/tvd"
	"math/big"
)

type User struct {
	Id                          int64                               `json:"id"`
	Username                    string                              `orm:"unique" json:"username"`
	Password                    string                              `json:"password"`
	Email                       string                              `orm:"unique" json:"email"`
	Roles                       []*Role                             `orm:"rel(m2m)" json:"roles"`
	JwtToken                    string                              `orm:"-" json:"token"`
	EtherumAddress              string                              `json:"etherumAddress"`
	EthereumBalance             *big.Int                            `orm:"-" json:"ethereumBalance"`
	APaymentTokenBalance        *big.Int                            `orm:"-" json:"apaymentTokenBalance"`
	Firstname                   string                              `json:"firstname"`
	Lastname                    string                              `json:"lastname"`
	PersonAddressResult         *tvd.PersonAddressResult            `orm:"-" json:"agateDetails"`
	AnimalHusbandryDetailResult *tvd.GetAnimalHusbandryDetailResult `orm:"-" json:"AnimalHusbandryDetailResult"`
	TVD                         int32                               `json:"tvd"`
}

func init() {
	// Register model
	orm.RegisterModel(new(User))
}

func (user *User) HasRole(roleName string) bool {
	for _, role := range user.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
