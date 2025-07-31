package utils

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/veritrans/go-midtrans"
)

func MidtransPayment(order_id string, amount int, name string, email string) midtrans.SnapResponse {
	godotenv.Load()
	serverKey := os.Getenv("MIDTRANS_SERVER_KEYS")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEYS")

	midclient := midtrans.NewClient()
	midclient.ClientKey = clientKey
	midclient.ServerKey = serverKey
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  string(order_id),
			GrossAmt: int64(amount),
		}, CustomerDetail: &midtrans.CustDetail{
			FName: name,
		},
	}

	snapUrl, _ := snapGateway.GetToken(snapReq)
	return snapUrl
}

func GetStatus(order_id string) midtrans.Response {
	godotenv.Load()
	serverKey := os.Getenv("MIDTRANS_SERVER_KEYS")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEYS")

	midclient := midtrans.NewClient()
	midclient.ClientKey = clientKey
	midclient.ServerKey = serverKey
	midclient.APIEnvType = midtrans.Sandbox

	gateway := midtrans.CoreGateway{
		Client: midclient,
	}

	transactionResp, _ := gateway.Status(order_id)

	return transactionResp

}
