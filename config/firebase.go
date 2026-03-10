package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

)

//instance Firebase Auth yang dipakai untuk verify token
var FirebaseAuth *auth.Client

func initFirebase(){
	credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")

	//Inisialisasi firebase App dengan services account credentials
	opt := option.WithAuthCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil{
		log.Fatalf("Gagal init firebase: %v" err)
	}
	// Mendapatkan Firebase Auth Client 
	FirebaseAuth, err = app.Auth(context.Background())
	if err != nil{
		log.Fatalf("Gaga; mendapatkan Firebase Auth Client: %v", err)
	}
}
