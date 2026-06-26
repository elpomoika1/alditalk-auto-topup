package browser

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NavigationResponse struct {
	UserDetails UserDetails `json:"userDetails"`
}

type UserDetails struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	ContractID string `json:"contractId"`
}

var contractIdUrl = "https://www.alditalk-kundenportal.de/scs/bff/scs-207-customer-master-data-bff/customer-master-data/v1/navigation-list"

func GetContractID(client *http.Client) (string, error) {
	req, err := NewClientRequest(
		"GET",
		contractIdUrl,
		"www.alditalk-kundenportal.de",
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data NavigationResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	if len(data.UserDetails.Subscriptions) == 0 {
		return "", fmt.Errorf("subscriptions is empty")
	}

	return data.UserDetails.Subscriptions[0].ContractID, nil
}
