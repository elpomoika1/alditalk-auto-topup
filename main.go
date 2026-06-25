package main

import (
	"aldi_renewal/v2/browser"
	"aldi_renewal/v2/math"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	SubscribedOffers []Offer `json:"subscribedOffers"`
}

type Offer struct {
	Pack []Pack `json:"pack"`
}

type Pack struct {
	Unit                      string `json:"unit"`
	NextExpirationDate        string `json:"nextExpirationDate"`
	Tariff                    string `json:"tariff"`
	Used                      string `json:"used"`
	Type                      string `json:"type"`
	Allocated                 string `json:"allocated"`
	BalanceAttributeReference string `json:"balanceAttributeReference"`
}

type UpdateRequest struct {
	OfferID               string `json:"offerId"`
	SubscriptionID        string `json:"subscriptionId"`
	UpdateOfferResourceID string `json:"updateOfferResourceID"`
	Amount                string `json:"amount"`
	RefillThresholdValue  string `json:"refillThresholdValue"`
}

const (
	offerUrl = "https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offers?contractId=%v&productType=Mobile_Product_Offer"
	postUrl  = "https://www.alditalk-kundenportal.de/scs/bff/scs-209-selfcare-dashboard-bff/selfcare-dashboard/v1/offer/updateUnlimited"
)

var (
	err error

	offerId               = flag.String("offerId", "some-32-hex", "uuid of offer")
	subscriptionId        = flag.String("subscriptionId", "12345678", "id of subscribition")
	updateOfferResourceID = flag.String("updateOfferResourceID", "30", "update offer resource id")
)

func main() {
	flag.Parse()

	for {
		fmt.Println("Checking...", time.Now())

		err := checkAndTopUp()
		if err != nil {
			log.Println("Check failed:", err)
		}

		time.Sleep(10 * time.Minute)
	}
}

func checkAndTopUp() error {
	client := &http.Client{}

	contractId, err := browser.GetContractID(client)

	req, err := browser.NewClientRequest(
		"GET",
		fmt.Sprintf(offerUrl, contractId),
		nil,
	)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("resp: %v", err)
		return err
	}

	defer resp.Body.Close()

	var packs Response

	err = json.NewDecoder(resp.Body).Decode(&packs)
	if err != nil {
		log.Fatalf("decode: %v", err)
		return err
	}

	remaining := 0.0

	for _, offer := range packs.SubscribedOffers {
		for _, pack := range offer.Pack {
			ballanceName := pack.BalanceAttributeReference
			if !strings.EqualFold(ballanceName, "dataGrantAmount") {
				continue
			}

			remaining, err = math.CalculateRemaining(pack.Allocated, pack.Used)
			if err != nil {
				log.Fatalf("Could not calculate remaining: %v", err)
				return err
			}

			fmt.Printf(
				"%s: remain %.2f GB\n",
				pack.BalanceAttributeReference,
				remaining,
			)
		}
	}

	if remaining < 0.9 {
		addGB()
	}

	return nil
}

func addGB() error {
	reqBody := UpdateRequest{
		OfferID:               *offerId,
		SubscriptionID:        *subscriptionId,
		UpdateOfferResourceID: *updateOfferResourceID,
		Amount:                "1048576",
		RefillThresholdValue:  "1048576",
	}

	jsonData, err := json.Marshal(reqBody)

	err = browser.PerformPost(
		postUrl,
		jsonData,
	)

	if err != nil {
		fmt.Println("Top up failed:", err)
		return err
	} else {
		fmt.Println("Success top up")
	}

	return nil
}
