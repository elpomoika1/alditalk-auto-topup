package main

import (
	"aldi_renewal/v2/browser"
	"aldi_renewal/v2/cli"
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
	loginUrl = "https://www.alditalk-kundenbetreuung.de/de"
)

var (
	err error

	offerId               = flag.String("offerId", "", "uuid of offer")
	subscriptionId        = flag.String("subscriptionId", "", "id of subscribition")
	updateOfferResourceID = flag.String("updateOfferResourceID", "", "update offer resource id")
)

func main() {
	client := &http.Client{}
	flag.Parse()

	cli.CheckArguments(*offerId, *subscriptionId, *updateOfferResourceID)

	for {
		fmt.Println("Checking...", time.Now())

		err := checkAndTopUp(client)
		if err != nil {
			log.Println("Check failed:", err)
		}

		time.Sleep(1 * time.Minute)
	}
}

func checkAndTopUp(client *http.Client) error {
	if err := refreshSession(client); err != nil {
		return err
	}

	contractId, err := browser.GetContractID(client)

	req, err := browser.NewClientRequest(
		"GET",
		fmt.Sprintf(offerUrl, contractId),
		"www.alditalk-kundenportal.de",
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

	fmt.Print(resp.Body)
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
				"\n%s: remain %.2f GB\n",
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
		"www.alditalk-kundenportal.de",
		jsonData,
	)

	if err != nil {
		fmt.Println("\nTop up failed:", err)
		return err
	} else {
		fmt.Println("\nSuccess top up")
	}

	return nil
}

func refreshSession(client *http.Client) error {
	req, err := browser.NewClientRequest(
		"GET",
		loginUrl,
		"www.alditalk-kundenbetreuung.de",
		nil)
	if err != nil {
		log.Fatalf("Error while requesting login url: %v", err)
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	browser.ApplyCookiesToReq(req, "www.alditalk-kundenbetreuung.de")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("final url:", resp.Request.URL)
	fmt.Println("refresh status:", resp.Status)

	return nil
}
