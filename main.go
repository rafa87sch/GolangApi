package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

//Item is into a single unit
type Post struct {
	Password string `json:"Password"`
}

func main() {
	router := mux.NewRouter()
	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes)
	fmt.Printf("encrypted : %s\n", key)
	router.HandleFunc("/api/encrypt", encrypt).Methods("POST")
	router.HandleFunc("/api/decrypt", decrypt).Methods("POST")
	http.ListenAndServe(":5000", router)
	fmt.Println("Server is up")
}

//Api to encrypt the passsword
func encrypt(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	//checking error before the Encrypt
	if r.Body == nil {
		w.WriteHeader(400)
		w.Write([]byte("Need to send a json file"))
		return
	}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Need to send a json file"))
		return
	}
	fmt.Println("Password before Encrypt", newPost.Password)
	if len(newPost.Password) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("The format is not correct"))
		return
	}
	//encrypt the information with the key that should be on the .env file but I have to delivery on time
	ciphertext := encryptMethod(newPost.Password, "a50ddcb14f5ea66f40acd1dc30b022c0baf8ae7c9a1d0407e3682c83a39b41b9")
	fmt.Printf("encrypted : %s\n", ciphertext)
	posts := &ciphertext
	w.Header().Set("Content-Type", "application/json")
	//returning the encrypt password
	json.NewEncoder(w).Encode(posts)
}

func decrypt(w http.ResponseWriter, r *http.Request) {
	//checking the err
	var newPost Post
	if r.Body == nil {
		w.WriteHeader(400)
		w.Write([]byte("Need to send a json file"))
		return
	}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Need to send a json file"))
		return
	}
	fmt.Println("test encrypt", newPost.Password)
	if len(newPost.Password) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("The format is not correct"))
		return
	}
	//decoding to take the value from json
	json.NewDecoder(r.Body).Decode(&newPost)
	//sending  a key number I will put on a .env to get but for get the delivery
	decrypted := decryptMethod(newPost.Password, "a50ddcb14f5ea66f40acd1dc30b022c0baf8ae7c9a1d0407e3682c83a39b41b9")
	fmt.Printf("decrypted : %s\n", decrypted)
	posts := decrypted
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func decryptMethod(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func encryptMethod(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}
