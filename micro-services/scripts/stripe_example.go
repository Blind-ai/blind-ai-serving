package main

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/source"
)

func main() {
	stripe.Key = "sk_test_361C3CRl7SawiARmWJxyTqmn00PhA1s0Cl"
	var amount int64 = 1
	p1 := &stripe.SourceObjectParams{
		Type:		stripe.String("wechat"),
		Amount:		&amount,
		Currency:	stripe.String(string(stripe.CurrencyEUR)),
		StatementDescriptor: stripe.String("ORDER AT11990"),
		Owner: 		&stripe.SourceOwnerParams{
			Email: stripe.String("fabien.lababe.desvignes@gmail.com"),
		},
	}
	session, err := source.New(p1)
	if err != nil {
		fmt.Println("%v", err)
		return
	}
	//png, err := qrcode.Encode(session.TypeData["qr_code_url"], qrcode.Medium, 256)
	str := fmt.Sprintf("%v", session.TypeData["qr_code_url"])
	fmt.Println("[", str, "]")
	err = qrcode.WriteFile(str, qrcode.Medium, 256, "qr.png")
}
