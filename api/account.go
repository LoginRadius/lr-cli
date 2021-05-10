package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/loginradius/lr-cli/request"
)

type SitesToken struct {
	APIVersion    string `json:"ApiVersion"`
	AppID         int64  `json:"AppId"`
	AppName       string `json:"AppName"`
	Authenticated bool   `json:"authenticated"`
	XSign         string `json:"xsign"`
	XToken        string `json:"xtoken"`
}

func SetSites(appid int64) (*SitesToken, error) {
	switchapp := conf.AdminConsoleAPIDomain + "/account/switchapp?appid=" + strconv.FormatInt(appid, 10)
	switchResp, err := request.Rest(http.MethodGet, switchapp, nil, "")
	var switchRespObj SitesToken
	err = json.Unmarshal(switchResp, &switchRespObj)
	if err != nil {
		return nil, err
	}
	return &switchRespObj, nil
}

type AccountPayment struct {
	Data struct {
		Order []struct {
			Totalamount         int         `json:"TotalAmount"`
			Recurringprofileid  interface{} `json:"RecurringProfileId"`
			Basediscount        int         `json:"BaseDiscount"`
			Promotionaldiscount int         `json:"PromotionalDiscount"`
			Initialamount       int         `json:"InitialAmount"`
			Tax                 int         `json:"Tax"`
			UUID                interface{} `json:"Uuid"`
			Invoiceno           int         `json:"InvoiceNo"`
			Createddate         time.Time   `json:"CreatedDate"`
			Lastmodifieddate    time.Time   `json:"LastModifiedDate"`
			Isactive            bool        `json:"IsActive"`
			Isdeleted           bool        `json:"IsDeleted"`
			Orderid             int         `json:"OrderId"`
			Paymentdetail       struct {
				Stripecustomerid      string `json:"StripeCustomerId"`
				Stripepaymentmethodid string `json:"StripePaymentMethodId"`
			} `json:"PaymentDetail"`
			Orderdetails []interface{} `json:"OrderDetails"`
		} `json:"Order"`
		Carddetails struct {
			Expmonth int    `json:"expMonth"`
			Expyear  int    `json:"expYear"`
			Last4    string `json:"last4"`
		} `json:"cardDetails"`
	} `json:"data"`
	Sharedsiteownerdata interface{} `json:"sharedSiteOwnerData"`
}

func PaymentInfo() (*AccountPayment, error) {
	payment := conf.AdminConsoleAPIDomain + "/account/accountpaymentdetail?"
	paymentResp, err := request.Rest(http.MethodGet, payment, nil, "")
	if err != nil {
		return nil, err
	}
	var paymentRespObj AccountPayment
	err = json.Unmarshal(paymentResp, &paymentRespObj)
	if err != nil {
		return nil, err
	}
	return &paymentRespObj, nil
}
