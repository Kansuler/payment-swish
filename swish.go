package swish

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"net/http"
	"time"
)

type webServiceURL string

const (
	testURL webServiceURL = "https://mss.cpc.getswish.net"
	prodURL webServiceURL = "https://cpc.getswish.net"
)

const (
	Certificate = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlEcnpDQ0FwZWdBd0lCQWdJUUNEdmdWcEJDUnJHaGRXckpXWkhIU2pBTkJna3Foa2lHOXcwQkFRVUZBREJoDQpNUXN3Q1FZRFZRUUdFd0pWVXpFVk1CTUdBMVVFQ2hNTVJHbG5hVU5sY25RZ1NXNWpNUmt3RndZRFZRUUxFeEIzDQpkM2N1WkdsbmFXTmxjblF1WTI5dE1TQXdIZ1lEVlFRREV4ZEVhV2RwUTJWeWRDQkhiRzlpWVd3Z1VtOXZkQ0JEDQpRVEFlRncwd05qRXhNVEF3TURBd01EQmFGdzB6TVRFeE1UQXdNREF3TURCYU1HRXhDekFKQmdOVkJBWVRBbFZUDQpNUlV3RXdZRFZRUUtFd3hFYVdkcFEyVnlkQ0JKYm1NeEdUQVhCZ05WQkFzVEVIZDNkeTVrYVdkcFkyVnlkQzVqDQpiMjB4SURBZUJnTlZCQU1URjBScFoybERaWEowSUVkc2IySmhiQ0JTYjI5MElFTkJNSUlCSWpBTkJna3Foa2lHDQo5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNGp2aEVYTGVxS1RUbzFlcVVLS1BDM2VReWFLbDdoTE9sbHNCDQpDU0RNQVpPblRqQzNVL2REeEdrQVY1M2lqU0xkaHdaQUFJRUp6czRiZzcvZnpUdHhSdUxXWnNjRnMzWW5Gbzk3DQpuaDZWZmU2M1NLTUkydGF2ZWd3NUJtVi9TbDBmdkJmNHE3N3VLTmQwZjNwNG1WbUZhRzVjSXpKTHYwN0E2RnB0DQo0M0MvZHhDLy9BSDJoZG1vUkJCWU1xbDFHTlhSb3I1SDRpZHE5Sm96K0VrSVlJdlVYN1E2aEwraHFrcE1mVDdQDQpUMTlzZGw2Z1N6ZVJudHdpNW0zT0ZCcU9hc3YremJNVVpCZkhXeW1lTXIveTd2clRDMExVcTdkQk10b00xTy80DQpnZFc3alZnL3RSdm9TU2lpY05veEJOMzNzaGJ5VEFwT0I2anRTajFldFgramtNT3ZKd0lEQVFBQm8yTXdZVEFPDQpCZ05WSFE4QkFmOEVCQU1DQVlZd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVUE5NVFOVmJSDQpUTHRtOEtQaUd4dkRsN0k5MFZVd0h3WURWUjBqQkJnd0ZvQVVBOTVRTlZiUlRMdG04S1BpR3h2RGw3STkwVlV3DQpEUVlKS29aSWh2Y05BUUVGQlFBRGdnRUJBTXVjTjZwSUV4SUsrdDFFbkU5U3NQVGZyZ1QxZVhrSW95UVkvRXNyDQpoTUF0dWRYSC92VEJIMWpMdUcyY2VuVG5tQ21yRWJYamNLQ2h6VXlJbVpPTWtYRGlxdzhjdnBPcC8yUFY1QWRnDQowNk8vblZzSjhkV080MVAwam1QNlA2ZmJ0R2JmWW1iVzBXNUJqZkl0dGVwM1NwK2RXT0lyV2NCQUkrMHRLSUpGDQpQbmxVa2lhWTRJQklxRGZ2OE5aNVlCYmVyT2dPelc2c1JCYzRMMG5hNFVVK0tyazJVODg2VUFiM0x1akVWMGxzDQpZU0VZMVFTdGVEd3NPb0JycCt1dkZSVHAySW5CdVRoczRwRnNpdjlrdVhjbFZ6REFHeVNqNGR6cDMwZDh0YlFrDQpDQVV3N0MyOUM3OUZ2MUM1cWZQcm1BRVNyY2lJeHBnMFg0MEtQTWJwMVpXVmJkND0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0NCg=="
)

// Options are settings that is used by the http client
type Options struct {
	// Passphrase is the password for the p12 encoded SSL certificate
	Passphrase string

	// SSLCertificate is a byte encoded array with the SSL certificate content
	SSLCertificate []byte

	// Test indicates whether the http client will use the test environment endpoint and CA certificate
	Test bool // enable test environment

	// CA is base64 encoded string with your certificate authority
	CA string

	// Timeout in seconds for the http client
	Timeout int // Client timeout in seconds
}

// Swish holds settings for this session
type Swish struct {
	client *http.Client
	test   bool

	// URL is the endpoint which we use to talk with BankID and can be replaced.
	URL string
}

// New creates a new client
func New(opts Options) (*Swish, error) {
	url := string(prodURL)
	if opts.Test {
		url = string(testURL)
	}

	blocks, err := pkcs12.ToPEM(opts.SSLCertificate, opts.Passphrase)
	if err != nil {
		return nil, err
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, err
	}

	ca, err := base64.StdEncoding.DecodeString(opts.CA)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(opts.Timeout),
	}

	return &Swish{
		client: client,
		URL:    url,
		test:   opts.Test,
	}, nil
}

type errorResponse struct {
	// ErrorCode is the short code for the error
	ErrorCode string `json:"errorCode"`
	// ErrorMessage is a more in-depth message about what went wrong
	ErrorMessage string `json:"errorMessage"`
	// AdditionalInformation
	AdditionalInformation string `json:"additionalInformation"`
}

// CreatePaymentRequestOptions for the create payment request
type CreatePaymentRequestOptions struct {
	// Required: The identifier of the payment request to be saved. Example 11A86BE70EA346E4B1C39C874173F088 or d2eb91f4-f3a7-4088-970f-a108b58bf8d9
	// The endpoint will format the string to fit Swish specification.
	InstructionUUID string `json:"-"`

	// Required: The endpoint Swish will call on with payment status updates, you need to receive data on this endpoint
	CallbackURL string `json:"callbackUrl"`

	// Required: The phone number that will receive the payment. Format E.164 except the plus ("+") symbol.
	PayeeAlias string `json:"payeeAlias"`

	// Required: The amount that is charged with a float value. Example "100.01"
	Amount string `json:"amount"`

	// Required: Currency code according to ISO 4217
	Currency string `json:"currency"`

	// Optional: Payment reference of the payee, which is the merchant that receives the payment. This reference could
	// be order id or similar. Allowed characters are a-z A-Z 0-9 -_.+*/ and length must be between 1 and 36 characters.
	PayeePaymentReference string `json:"payeePaymentReference,omitempty"`

	// Optional: The registered cellphone number of the person that makes the payment. It can only contain numbers and
	// has to be at least 8 and at most 15 numbers. It also needs to match the following format in order to be found in
	// Swish: country code + cellphone number (without leading zero). E.g.: 46712345678
	PayerAlias string `json:"payerAlias,omitempty"`

	// Optional: The social security number of the individual making the payment, should match the registered value for
	// payerAlias or the payment will not be accepted. The value should be a proper Swedish social security number
	// (personnummer or sammordningsnummer).
	PayerSSN string `json:"payerSSN,omitempty"`

	// Optional: Minimum age (in years) that the individual connected to the payerAlias has to be in order for the
	// payment to be accepted. Value has to be in the range of 1 to 99.
	PayerAgeLimit string `json:"payerAgeLimit,omitempty"`

	// Optional: Merchant supplied message about the payment/order. Max 50 chars. Allowed characters are the letters
	// a-ö, A-Ö, the numbers 0-9 and the special characters :;.,?!()”.
	Message string `json:"message,omitempty"`
}

type createPaymentRequestResponse struct {
	// Location is an URL that you use as GET to retrieve the status of the payment request
	Location string
	// PaymentRequestToken is returned when creating an m-commerce payment request. The token to use when opening the
	// Swish app.
	PaymentRequestToken string
	// ErrorCodes returns error codes
	ErrorCodes []errorResponse
}

// CreatePaymentRequest sends a v2 payment request to Swish to create a payment
func (s *Swish) CreatePaymentRequest(ctx context.Context, opts CreatePaymentRequestOptions) (result createPaymentRequestResponse, err error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/swish-cpcapi/api/v2/paymentrequests/%s", s.URL, opts.InstructionUUID), bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		err = json.NewDecoder(resp.Body).Decode(&result.ErrorCodes)
		if err != nil {
			return
		}

		var errs string
		for _, errCode := range result.ErrorCodes {
			if len(errs) > 0 {
				errs += " | "
			}
			errs += fmt.Sprintf("[%s] %s", errCode.ErrorCode, errCode.ErrorMessage)
		}

		return result, errors.New(errs)
	}

	if resp.StatusCode == http.StatusForbidden {
		result.ErrorCodes = append(result.ErrorCodes, errorResponse{
			ErrorCode:             "PA01",
			ErrorMessage:          "The payeeAlias in the payment request object is not the same as merchant’s Swish number",
			AdditionalInformation: "",
		})
		return result, errors.New("[PA01] The payeeAlias in the payment request object is not the same as merchant’s Swish number")
	}

	result.Location = resp.Header.Get("Location")
	result.PaymentRequestToken = resp.Header.Get("Paymentrequesttoken")

	return
}

type statusResponse struct {
	// InstructionUUID is the ID that the request was created with
	InstructionUUID string `json:"id"`

	// PayeePaymentReference Payment reference of the payee, which is the merchant that receives the payment. This
	// reference could be order id or similar. Allowed characters are a-z A-Z 0-9 -_.+*/ and length must be between 1
	// and 36 characters.
	PayeePaymentReference string `json:"payeePaymentReference"`

	// PaymentReference Payment reference, from the bank, of the payment that occurred based on the Payment request.
	// Only available if status is PAID.
	PaymentReference string `json:"paymentReference"`

	// CallbackURL URL that Swish will use to notify caller about the outcome of the Payment request. The URL has to
	// use HTTPS.
	CallbackURL string `json:"callbackUrl"`

	// PayerAlias The registered cellphone number of the person that makes the payment. It can only contain numbers and
	// has to be at least 8 and at most 15 numbers. It also needs to match the following format in order to be found in
	// Swish: country code + cellphone number (without leading zero). E.g.: 46712345678
	PayerAlias string `json:"payerAlias"`

	// PayerSSN The social security number of the individual making the payment, should match the registered value for
	// payerAlias or the payment will not be accepted. The value should be a proper Swedish social security number
	// (personnummer or sammordningsnummer).
	PayerSSN string `json:"payerSSN"`

	// PayeeAlias The Swish number of the payee.
	PayeeAlias string `json:"payeeAlias"`

	// Amount The amount of money to pay. The amount cannot be less than 0.01 SEK and not more than 999999999999.99 SEK.
	// Valid value has to be all numbers or with 2-digit decimal separated by a period.
	Amount float64 `json:"amount"`

	// Currency The currency to use. The only currently supported value is SEK
	Currency string `json:"currency"`

	// Message Merchant supplied message about the payment/order. Max 50 chars. Allowed characters are the letters a-ö, A-Ö,
	// the numbers 0-9 and the special characters :;.,?!()”.
	Message string `json:"message"`

	// Status The status of the transaction. Possible values: CREATED, PAID, DECLINED, ERROR.
	Status string `json:"status"`

	// DateCreated The time and date that the payment request was created.
	DateCreated time.Time `json:"dateCreated"`

	// DatePaid The time and date that the payment request was paid. Only applicable if status was PAID.
	DatePaid time.Time `json:"datePaid"`

	// ErrorCode A code indicating what type of error occurred. Only applicable if status is ERROR.
	ErrorCode string `json:"errorCode"`

	// ErrorMessage A descriptive error message (in English) indicating what type of error occurred. Only applicable if status is
	// ERROR.
	ErrorMessage string `json:"errorMessage"`

	// AdditionalInformation Additional information about the error. Only applicable if status is ERROR.
	AdditionalInformation string `json:"additionalInformation"`
}

// Status use the location header from other endpoints to get status from Swish
func (s *Swish) Status(ctx context.Context, Location string) (result statusResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", Location, nil)
	if err != nil {
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		var errCodes []errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errCodes)
		if err != nil {
			return
		}

		var errs string
		for _, errCode := range errCodes {
			if len(errs) > 0 {
				errs += " | "
			}
			errs += fmt.Sprintf("[%s] %s", errCode.ErrorCode, errCode.ErrorMessage)
			result.ErrorCode = errCode.ErrorCode
			result.ErrorMessage = errCode.ErrorMessage
		}

		return result, errors.New(errs)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

// CreateRefundOptions for create refund
type CreateRefundOptions struct {
	// Required: InstructionUUID The ID for this refund, should be different from create payment request InstructionUUID
	InstructionUUID string `json:"-"`

	// Required: OriginalPaymentReference Reference of the original payment that this refund is for.
	OriginalPaymentReference string `json:"originalPaymentReference"`

	// Required: CallbackURL URL that Swish will use to notify caller about the outcome of the refund. The URL has to
	// use HTTPS.
	CallbackURL string `json:"callbackUrl"`

	// Required: PayerAlias The Swish number of the merchant that makes the refund payment.
	PayerAlias string `json:"payerAlias"`

	// Required: Amount The amount of money to refund. The amount cannot be less than 0.01 SEK and not more than
	// 999999999999.99 SEK. Moreover, the amount cannot exceed the remaining amount of the original payment that the
	// refund is for.
	Amount string `json:"amount"`

	// Required: Currency The currency to use. The only currently supported value is SEK.
	Currency string `json:"currency"`

	// Optional: PayerPaymentReference Payment reference supplied by the merchant. This could be order id or similar.
	PayerPaymentReference string `json:"payerPaymentReference"`

	// Optional: Merchant supplied message about the refund. Max 50 chars. Allowed characters are the letters a-ö, A-Ö,
	// the numbers 0-9 and the special characters :;.,?!()”.
	Message string `json:"message"`
}

type createRefundResponse struct {
	// Location is an URL that you use as GET to retrieve the status of the payment request
	Location string
	// ErrorCodes returns error codes
	ErrorCodes []errorResponse
}

// CreateRefund A merchant that has received a Swish payment can refund the whole or part of the original transaction
// amount to the consumer.
func (s *Swish) CreateRefund(ctx context.Context, opts CreateRefundOptions) (result createRefundResponse, err error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/swish-cpcapi/api/v2/refunds/%s", s.URL, opts.InstructionUUID), bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		err = json.NewDecoder(resp.Body).Decode(&result.ErrorCodes)
		if err != nil {
			return
		}

		var errs string
		for _, errCode := range result.ErrorCodes {
			if len(errs) > 0 {
				errs += " | "
			}
			errs += fmt.Sprintf("[%s] %s", errCode.ErrorCode, errCode.ErrorMessage)
		}

		return result, errors.New(errs)
	}

	if resp.StatusCode == http.StatusForbidden {
		result.ErrorCodes = append(result.ErrorCodes, errorResponse{
			ErrorCode:             "PA01",
			ErrorMessage:          "The payeeAlias in the payment request object is not the same as merchant’s Swish number",
			AdditionalInformation: "",
		})
		return result, errors.New("[PA01] The payeeAlias in the payment request object is not the same as merchant’s Swish number")
	}

	result.Location = resp.Header.Get("Location")

	return
}
