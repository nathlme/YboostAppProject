package main

import(
	"golang.org/x/crypto/bcrypt"
)


func PasswordHash(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password,bcrypt.DefaultCost) 
	if err != nil {
		return "",err
	}
	return string(hash), nil
}



func validateRegister(email string,pseudo string,password string) string{
	msg := ""
	if !(len(email) >= 8 ){
		msg = msg + "8 caractères minimum pour l'email requis, "
	}
	if !(len(pseudo) >= 3){
		msg = msg + "Pseudo doit contenir au moins 3 caractères, "
	}
	if !(len(password) >= 8){
		msg = msg + "8 caractères minimum pour le mot de passe. "
	}
	return msg 
}