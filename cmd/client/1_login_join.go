package main

import (
	"errors"
)

func progressLoginOrJoin() (msg string, ok bool) {
	var err error
	switch SEQ_PROGRESS {
	case 0:
		return "your id", true
	case 1:
		STORED_PARAM["id"] = LAST_INPUT
		return "your password", true
	case 2:
		STORED_PARAM["password"] = LAST_INPUT
		if IN_PROGRESS == PROGRESS_JOIN {
			id, pass, _ := getIdPassword()
			if err = join(id, pass); err == nil {
				err = createAccount()
			}
		} else {
			err = login()
		}
		DONE_PROGRESS = true
	}
	if err != nil {
		return err.Error(), false
	}
	return "success", true
}

func getIdPassword() (string, string, bool) {
	id := STORED_PARAM["id"]
	password := STORED_PARAM["password"]
	return id, password, id != "" && password != ""
}

func join(user, pass string) error {
	_, _, err := post("/svc/v1/users/join", map[string]interface{}{
		"userId":   user,
		"password": pass,
	})
	return err
}

func createAccount() error {
	_, _, err := post("/svc/v1/users/accounts/create", map[string]interface{}{
		"owner": ADDRESS_EOA,
	})
	return err
}

func login() error {
	_, code, err := get("/svc/v1/users/accounts/state", nil)
	if code == 1003 || err == nil {
		LOGGED_IN = true
		return nil
	}
	return errors.New("check id, password")
}
