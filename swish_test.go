package swish_test

import (
	"context"
	swish "github.com/Kansuler/payment-swish"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNew(t *testing.T) {
	cert, err := ioutil.ReadFile("certificates/Swish_Merchant_TestCertificate_1234679304.p12")
	if err != nil {
		t.Fatalf("could not load test certificate: %s", err.Error())
	}

	s, err := swish.New(swish.Options{
		Passphrase:     "swish",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})

	assert.NoError(t, err)
	assert.NotNil(t, s)

	s, err = swish.New(swish.Options{
		Passphrase:     "hsiws",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})

	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestSwish_CreatePaymentRequest(t *testing.T) {
	cert, err := ioutil.ReadFile("certificates/Swish_Merchant_TestCertificate_1234679304.p12")
	if err != nil {
		t.Fatalf("could not load test certificate: %s", err.Error())
	}

	s, err := swish.New(swish.Options{
		Passphrase:     "swish",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})

	response, err := s.CreatePaymentRequest(context.Background(), swish.CreatePaymentRequestOptions{
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
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, response.PaymentRequestToken)
	assert.NotEmpty(t, response.Location)
	assert.Empty(t, response.ErrorCodes)

	response, err = s.CreatePaymentRequest(context.Background(), swish.CreatePaymentRequestOptions{
		InstructionUUID:       uuid.NewV4().String(),
		CallbackURL:           "",
		PayeeAlias:            "",
		Amount:                "",
		Currency:              "",
		PayeePaymentReference: "",
		PayerAlias:            "",
		PayerSSN:              "",
		PayerAgeLimit:         "",
		Message:               "",
	})

	assert.Error(t, err)
	assert.Empty(t, response.PaymentRequestToken)
	assert.Empty(t, response.Location)
	assert.NotEmpty(t, response.ErrorCodes)
}

func TestSwish_Status(t *testing.T) {
	cert, err := ioutil.ReadFile("certificates/Swish_Merchant_TestCertificate_1234679304.p12")
	if err != nil {
		t.Fatalf("could not load test certificate: %s", err.Error())
	}

	s, err := swish.New(swish.Options{
		Passphrase:     "swish",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})

	request := swish.CreatePaymentRequestOptions{
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

	response, err := s.CreatePaymentRequest(context.Background(), request)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.PaymentRequestToken)
	assert.NotEmpty(t, response.Location)
	assert.Empty(t, response.ErrorCodes)

	status, err := s.Status(context.Background(), response.Location)
	assert.Equal(t, request.InstructionUUID, status.InstructionUUID)

	// Incomplete UUID
	status, err = s.Status(context.Background(), "https://mss.cpc.getswish.net/swish-cpcapi/api/v1/paymentrequests/771178D5BF45450882F1B53681D")
	assert.Error(t, err)
	assert.Equal(t, status.ErrorCode, "RP04")
}

func TestSwish_CreateRefund(t *testing.T) {
	cert, err := ioutil.ReadFile("certificates/Swish_Merchant_TestCertificate_1234679304.p12")
	if err != nil {
		t.Fatalf("could not load test certificate: %s", err.Error())
	}

	s, err := swish.New(swish.Options{
		Passphrase:     "swish",
		CA:             swish.TestCertificate,
		SSLCertificate: cert,
		Test:           true,
		Timeout:        5,
	})

	request := swish.CreatePaymentRequestOptions{
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

	response, err := s.CreatePaymentRequest(context.Background(), request)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.PaymentRequestToken)
	assert.NotEmpty(t, response.Location)
	assert.Empty(t, response.ErrorCodes)

	refund, err := s.CreateRefund(context.Background(), swish.CreateRefundOptions{
		InstructionUUID:          uuid.NewV4().String(),
		OriginalPaymentReference: request.InstructionUUID,
		CallbackURL:              request.CallbackURL,
		PayerAlias:               request.PayeeAlias,
		Amount:                   request.Amount,
		Currency:                 request.Currency,
		PayerPaymentReference:    "123",
		Message:                  "Refund",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, refund.Location)
	assert.Empty(t, refund.ErrorCodes)

	refund, err = s.CreateRefund(context.Background(), swish.CreateRefundOptions{
		InstructionUUID:          uuid.NewV4().String(),
		OriginalPaymentReference: request.InstructionUUID,
		CallbackURL:              request.CallbackURL,
		PayerAlias:               request.PayerAlias, // Invalid reference
		Amount:                   request.Amount,
		Currency:                 request.Currency,
		PayerPaymentReference:    "123",
		Message:                  "Refund",
	})

	assert.Error(t, err)
	assert.Empty(t, refund.Location)
	assert.NotEmpty(t, refund.ErrorCodes)
}
