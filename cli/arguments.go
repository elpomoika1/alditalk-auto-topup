package cli

import "fmt"

func CheckArguments(offerId, subscriptionId, updateOfferResourceID string) {
	if offerId == "" {
		offerId = mustAsk("Enter offerId")
	}

	if subscriptionId == "" {
		subscriptionId = mustAsk("Enter subscriptionId")
	}

	if updateOfferResourceID == "" {
		updateOfferResourceID = mustAsk("Enter updateOfferResourceID")
	}

	fmt.Println("OfferId:", offerId)
	fmt.Println("SubscriptionId:", subscriptionId)
	fmt.Println("ResourceId:", updateOfferResourceID)
}

func mustAsk(prompt string) string {
	var input string

	fmt.Print(prompt + ": ")
	fmt.Scanln(&input)

	return input
}
