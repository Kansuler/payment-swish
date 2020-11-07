# Kansuler/payment-swish

![License](https://img.shields.io/github/license/Kansuler/payment-swish) ![Version](https://img.shields.io/github/go-mod/go-version/Kansuler/payment-swish) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/8be4b74cbfa74fb1bcc6b38bdba52aed)](https://www.codacy.com/gh/Kansuler/payment-swish/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Kansuler/payment-swish&amp;utm_campaign=Badge_Grade)

A package to simplify integrations against the Swedish payment solution [Swish](https://www.swish.nu/).

It is recommended to read through the [integration guide](https://www.swish.nu/developer#swish-for-merchants) thoroughly to understand the the process, and what responses that can occur.

API and detailed documentation can be found at [https://godoc.org/github.com/Kansuler/payment-swish](https://godoc.org/github.com/Kansuler/payment-swish)

This library implements version 2 of create payment request and create refund.

## Installation

`go get github.com/Kansuler/payment-swish`

## Functions

```go

```

## Usage

```go
package main

import (
	"context"
	swish "github.com/Kansuler/payment-swish"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
)

func main() {
	cert, err := ioutil.ReadFile("certificates/Swish_Merchant_TestCertificate_1234679304.p12")
	if err != nil {
		log.Fatalf("could not load test certificate: %s", err.Error())
	}

	s, err := swish.New(swish.Options{
		Passphrase:     "swish",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})
	if err != nil {
		log.Fatalf("could not create swish instance: %s", err.Error())
	}

	paymentRequest := swish.CreatePaymentRequestOptions{
		InstructionUUID:       uuid.NewV4().String(),
		CallbackURL:           "https://localhost:8080/callback",
		PayeeAlias:            "1234679304",
		Amount:                "100.01",
		Currency:              "SEK",
		PayeePaymentReference: "",
		PayerAlias:            "",
		PayerSSN:              "",
		PayerAgeLimit:         "",
		Message:               "",
	}

	payment, err := s.CreatePaymentRequest(context.Background(), paymentRequest)
	if err != nil {
		log.Fatalf("could not create payment paymentRequest: %s", err.Error())
	}

	_, err = s.Status(context.Background(), payment.Location)
	if err != nil {
		log.Fatalf("could not get status of payment: %s", err.Error())
	}

	refundRequest := swish.CreateRefundOptions{
		InstructionUUID:          uuid.NewV4().String(),
		OriginalPaymentReference: paymentRequest.InstructionUUID,
		CallbackURL:              paymentRequest.CallbackURL,
		PayerAlias:               paymentRequest.PayeeAlias,
		Amount:                   paymentRequest.Amount,
		Currency:                 paymentRequest.Currency,
		PayerPaymentReference:    "123",
		Message:                  "Refund",
	}

	refund, err := s.CreateRefund(context.Background(), refundRequest)
	if err != nil {
		log.Fatalf("could not create refund: %s", err.Error())
	}

	_, err = s.Status(context.Background(), refund.Location)
	if err != nil {
		log.Fatalf("could not get status of refund: %s", err.Error())
	}
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
