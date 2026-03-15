package delhivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/C2NOfficial/C2NGCShared/models"
)

// Sends the request to the cloud function to get the estimated shipping cost. 
func EstimateShippingPrice(request *ShippingEstimateRequest) (float64, error) {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		GetEstimatedShippingCostURL,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	//cloud function returns a float64 value which is the estimated shipping cost
	var successResponse models.SuccessPayload[float64]

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&successResponse)
	if err != nil {
		return 0, err
	}

	return successResponse.Data, nil
}

// Sends the request directly to delhivery official api to create an order
//
// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/order-creation
func CreateOrder(payload *OrderRequest) (*OrderCreationResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("format", "json")
	formData.Set("data", string(payloadBytes))
	url := "https://track.delhivery.com/api/cmu/create.json"
	newReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Add("Accept", "application/json")
	newReq.Header.Set("Authorization", "Token "+API_TOKEN)

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var shipmentJSONResponse OrderCreationResponse
	err = json.NewDecoder(resp.Body).Decode(&shipmentJSONResponse)

	if err != nil {
		return nil, err
	}
	return &shipmentJSONResponse, nil
}

// Sends the request directly to delhivery official api to create pickup request
// 
// The parameter pickupRequestURLValues can be made by calling the method 
// ToPickupRequestUrlValues of struct CreateShipmentPayload which is 
// sent in the body of the request made by the admin from frontend
//
// No need to create a separate struct and take the response. Just read the response bytes 
// and print them as string if an error occurs. Since a pickup request will only be made when 
// order is created in delhivery so we do not care about the response body if successful. 
//
// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/pickup-scheduling
func CreatePickupRequest(pickupRequetURLValues url.Values) error {
	url := "https://track.delhivery.com/fm/request/new/"
	pickupReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(pickupRequetURLValues.Encode()))
	if err != nil {
		return err
	}
	pickupReq.Header.Set("Authorization", "Token "+API_TOKEN)
	pickupReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Send the request
	pickupResp, err := http.DefaultClient.Do(pickupReq)
	if err != nil {
		return err
	}

	//Read the response body
	defer pickupResp.Body.Close()
	body, err := io.ReadAll(pickupResp.Body)
	if err != nil {
		return err
	}

	if pickupResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d. Response body: %s", pickupResp.StatusCode, string(body))
	}

	return nil
}

// Sends the request directly to delhivery official api to track an order
func TrackOrder(wayBill string) (*TrackingResponse, error) {
	apiURL := "https://track.delhivery.com/api/v1/packages/json/?waybill=" + wayBill
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+API_TOKEN)

	//Send the request
	pickupResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//Read the response body
	defer pickupResp.Body.Close()
	body, err := io.ReadAll(pickupResp.Body)
	if err != nil {
		return nil, err
	}
	if pickupResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d. Response body: %s", pickupResp.StatusCode, string(body))
	}

	var trackingResponse TrackingResponse
	err = json.Unmarshal(body, &trackingResponse)
	if err != nil{
		return nil, err
	}

	return &trackingResponse, nil
}