package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"paymentfc/models"
)

type XenditClient interface {
	CreateInvoice(ctx context.Context, param models.XenditInvoiceRequest) (models.XenditInvoiceResponse, error)
	CheckInvoiceStatus(ctx context.Context, externalID string) (string, error)
}

type xenditClient struct {
	APISecretKey string
}

func NewXenditClient(apiSecretKey string) XenditClient {
	return &xenditClient{
		APISecretKey: apiSecretKey,
	}
}

func (c *xenditClient) CreateInvoice(ctx context.Context, param models.XenditInvoiceRequest) (models.XenditInvoiceResponse, error) {
	var result models.XenditInvoiceResponse
	payload, err := json.Marshal(param)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}
	uri := "https://api.xendit.co/v2/invoices"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(payload))
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}

	httpReq.SetBasicAuth(c.APISecretKey, "")
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		//invalid response
		body, _ := io.ReadAll(resp.Body)
		return models.XenditInvoiceResponse{}, errors.New(fmt.Sprintf("xendit.createinvoice got error : %s", string(body)))

	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}

	return result, nil

}
func (c *xenditClient) CheckInvoiceStatus(ctx context.Context, externalID string) (string, error) {
	url := fmt.Sprintf("https://api.xendit.co/v2/invoices?external_id=%s", externalID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	httpReq.SetBasicAuth(c.APISecretKey, "")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//expected response --> models.XenditInvoiceResponse
	var response []models.XenditInvoiceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return "", err
	}

	return response[0].Status, nil

}
