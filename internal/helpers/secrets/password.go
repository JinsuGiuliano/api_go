package secrets

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const pepperKey = "4234kxzjcjj3nxnxbcvsjfj"

func Encrypt(password string) (hashedPwd string, err error) {

	pwdMacStr := applyPepper(password, pepperKey)

	// pwdHashed contains the hashed pepperd password and the salt hash
	pwdHashed, err := bcrypt.GenerateFromPassword([]byte(pwdMacStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(pwdHashed), nil
}

func Check(password, hashedPassword string) (err error) {

	pwdMacStr := applyPepper(password, pepperKey)

	pwdHashedOut, err := base64.URLEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword(pwdHashedOut, []byte(pwdMacStr)); err != nil {
		fmt.Println("Not Equal passwords! :( ", err)
		return err
	}

	fmt.Println("Equal passwords! :) ")
	return nil
}

func applyPepper(password, pepper string) (pepperedPwd string) {
	// Hash password with pepper
	key := []byte(pepper)
	pwd := []byte(password)

	mac := hmac.New(sha512.New, key)
	mac.Write(pwd)
	pwdMac := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString(pwdMac)
}
