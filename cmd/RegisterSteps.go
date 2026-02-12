package main

import(
	"golang.org/x/crypto/bcrypt"
)

// Hash the password 
func PasswordHash(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password,bcrypt.DefaultCost) 
	if err != nil {
		return "",err
	}
	return string(hash), nil
}


// Check the size of email and pseudo 
func validateRegister(email string,pseudo string,password string) string{
	msg := ""
	if !(len(email) >= 8 ){
		msg = msg + "8 caractères minimum pour l'email requis, "
	}
	if !(len(pseudo) >= 3){
		msg = msg + "Pseudo doit contenir au moins 3 caractères "
	}

	msg = msg + StrongPassword(password)

	return msg 
}


// Check if the password has at least a MAJ and a size of 8
func StrongPassword(password string) string {
	msg := ""
	VerMaj := " Votre mot de passe doit contenir au moins une majuscule."
	
	if len(password) < 8 {
		msg = msg + "Votre mot de passe doit contenir au moins 8 caractères"
	}

	for i := 0; i < len(password); i++ {
		if password[i] >= 'A' && password[i] <= 'Z' {
			VerMaj = ""
		}
	}

	msg = msg + VerMaj

	return msg

}