# ALDI Auto Top-Up

Automated tool that monitors remaining mobile data in ALDI TALK and automatically triggers a 1 GB refill when the threshold is reached.

---

## Features

- Automatically requests 1 GB refill when low
- Works with ALDI TALK customer portal API
- Lightweight Go-based tool
- Grab Cookies from your browser

---

## How it works

1. Fetches contract and subscription data from ALDI TALK portal
2. Calculates remaining mobile data
3. If threshold is reached → sends refill request (1 GB)
4. Runs continuously (recommended every 10 minutes)

---

## How to use

### 1. Download release

Go to **Releases** and download aldi_checker.exe

---

### 2. Get required parameters

You need to extract the following values from aldi_checker:

- `offerId`
- `subscriptionId`
- `updateOfferResourceID`

These values are required for refill requests.

---

### 3. Run the tool

#### Windows
cmd
start /B aldi_auto.exe --offerId "..." --subscriptionId "..." --updateOfferResourceID "..."

#### Linux
git clone https://github.com/elpomoika1/alditalk-auto-topup.git
cd alditalk-auto-topup
go run main.go --offerId "..." --subscriptionId "..." --updateOfferResourceID "..."
or
go build -o main.go and run it

##### OR

#### Windows/Linux
run aldi_auto.exe as default file (or in Linux build it) and enter your credits from aldi_checker.exe

---

## ⚠️ Troubleshooting

If you encounter the following error: `decode: json: cannot unmarshal array into Go value of type main.Response`

this usually means that your session is not fully initialized or valid.

#### Solution

Simply open the ALDI TALK customer portal in your browser and log in manually:

[alditalk](https://www.alditalk-kundenportal.de/portal/auth/uebersicht/)

This will refresh and initialize the required session cookies.  
After that, restart the tool.

---
