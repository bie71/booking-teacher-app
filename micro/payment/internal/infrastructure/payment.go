package infrastructure

import (
	"os"
	"payment/internal/config"

	"github.com/spf13/cast"
	midtrans "github.com/veritrans/go-midtrans"
)

type PaymentService struct {
	midtransClient *midtrans.Client
	midtransConfig config.Midtrans
}

func NewPaymentService(midtransClient *midtrans.Client, midtransConfig config.Midtrans) *PaymentService {
	return &PaymentService{
		midtransClient: midtransClient,
		midtransConfig: midtransConfig,
	}
}

func (p *PaymentService) CreatePayment(orderID string, bookingID uint, amount int64) (string, error) {

	if cast.ToBool(os.Getenv("IS_NFT")) {
		return "link-mock-nft", nil
	}

	snapGateway := midtrans.SnapGateway{Client: *p.midtransClient}

	req := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	resp, err := snapGateway.GetToken(req)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}
