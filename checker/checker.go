package main

import (
	"aldi_renewal/v2/browser"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Responce struct {
	SubscribedOffers []SubscribedOffers `json:"subscribedOffers"`
}

type SubscribedOffers struct {
	OfferId        string `json:"offerId"`
	SubscriptionId string `json:"subscriptionId"`
	ResourceId     string `json:"resourceId"`
}

const (
	offerUrl = "https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers?contractId=%v&productType=Mobile_Product_Offer"
)

func main() {
	client := &http.Client{}

	contractId, err := browser.GetContractID(client)

	req, err := browser.NewClientRequest(
		"GET",
		fmt.Sprintf(offerUrl, contractId),
		nil,
	)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("resp: %v", err)
	}
	defer resp.Body.Close()

	var responce Responce

	err = json.NewDecoder(resp.Body).Decode(&responce)
	if err != nil {
		log.Fatalf("decode: %v", err)
	}

	for _, offer := range responce.SubscribedOffers {
		fmt.Printf("OfferId: %v\n", offer.OfferId)
		fmt.Printf("SubscriptionId: %v\n", offer.SubscriptionId)
		fmt.Printf("ResourceId: %v\n", offer.ResourceId)
	}
}
