package controllers

import (
	"github.com/astaxie/beego"

	"github.com/scmo/apayment-backend/models"
	"encoding/json"
	"github.com/scmo/apayment-backend/services"
	"github.com/ethereum/go-ethereum/common"
)

// Operations about APaymentToken
type APaymentTokenController struct {
	beego.Controller
}


// @Title Creates a new APayment Token Transfer
// @Description Endpoint to transfer APayment Token from the System Account to the selected account
// @Param	body		body 	models.APaymentTokenTransfer	true		"body for request content"
// @Success 200 {Object} models.APaymentTokenTransfer
// @Failure 403 body is empty
// @router / [post]
func (this *APaymentTokenController) Transfer() {
	var aPaymentTokenTransfer models.APaymentTokenTransfer

	json.Unmarshal(this.Ctx.Input.RequestBody, &aPaymentTokenTransfer)

	claims, _ := services.ParseToken(this.Ctx.Request.Header.Get("Authorization"))
	user, err := services.GetUserByUsername(claims.Subject)
	if err != nil {
		this.CustomAbort(404, err.Error())
	}

	if ( !user.HasRole("Admin") ) {
		this.CustomAbort(404, err.Error())
	}

	aPaymentTokenTransfer.From = beego.AppConfig.String("systemAccountAddress")
	beego.Debug(aPaymentTokenTransfer)

	err = services.Transfer(common.HexToAddress(aPaymentTokenTransfer.From), common.HexToAddress(aPaymentTokenTransfer.To), aPaymentTokenTransfer.Amount)
	if (err != nil ) {
		beego.Error("Error while tranfering tokens. ", err)
		this.CustomAbort(500, err.Error())
	}
	this.ServeJSON()
}


// @Title Get Transactions
// @Description get all transactions
// @Success 200 {Object} models.APaymentTokenTransaction
// @router /transactions [get]
func (this *APaymentTokenController) GetAllTransactions() {

	claims, _ := services.ParseToken(this.Ctx.Request.Header.Get("Authorization"))
	user, err := services.GetUserByUsername(claims.Subject)
	if err != nil {
		this.CustomAbort(404, err.Error())
	}

	if ( (user.HasRole("Admin") || user.HasRole("Canton") ) == false ) {
		this.CustomAbort(401, "Unauthorized")
	}

	transactions, err := services.GetTransactions()
	if (err != nil) {
		beego.Error("Error while getting transactions. ", err)
		this.CustomAbort(500, err.Error())
	}
	this.Data["json"] = transactions
	this.ServeJSON()
}