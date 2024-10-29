package main

func faucet() (string, bool) {
	if _, _, err := post("/svc/v1/faucet", nil); err != nil {
		return err.Error(), false
	}
	return "faucet requested. type 'state' to check your balance", false
}
