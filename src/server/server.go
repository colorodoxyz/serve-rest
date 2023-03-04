package main

/**
 * A TLS enabled server that maintains a basic REST API at /api/keys, enabling
 * the read, write, and delete of Key-Value pairs
 *
 * The endpoint is secured using a JSON web token that is provided via the login
 * endpoint at /api/login
 */

import (
	"log"
	"net/http"

	"github.com/colorodoxyz/serve-rest/src/helper"
	"github.com/colorodoxyz/serve-rest/src/jwtMiddleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

var store = make(map[string]helper.KeyValue)

/**
 * Login at /api/login by posting {"Username": "{username}", "Password": "{password}"}
 */
func login(ctxt *gin.Context) {
	var account helper.Account

	var err error
	if err = ctxt.BindJSON(&account); err != nil {
		ctxt.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	if account.Username == helper.AdminUser && account.Password == helper.AdminPassword {
		token, tokenErr := jwtMiddleware.GenerateToken(account.Username)
		if tokenErr != nil {
			ctxt.JSON(http.StatusBadRequest, gin.H{"error": tokenErr.Error()})
			return
		}
		ctxt.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	ctxt.JSON(http.StatusBadRequest, gin.H{"error": "The username or password is incorrect"})
}

/**
 * Store key-value json struct in the store map
 * {
 *	 "key": "{key}",
 *	 "value": "{value}"
 * }
 */
func storeKeyValue(ctxt *gin.Context) {
	jwtMiddleware.ValidateJwt(ctxt)

	var keyValuePair helper.KeyValue

	if err := ctxt.BindJSON(&keyValuePair); err != nil {
		log.Println("Invalid json provided")
		return
	}

	log.Printf("Writing %s\n", keyValuePair)
	store[keyValuePair.Key] = keyValuePair
	ctxt.IndentedJSON(http.StatusCreated, keyValuePair)
}

/**
 * Return all KeyValue pairs as a JSON array of key-value structs
 */
func getAllKeyValuePairs(ctxt *gin.Context) {
	jwtMiddleware.ValidateJwt(ctxt)

	mapVals := maps.Values(store)
	log.Printf("Retrieving all key-value pairs: %s\n", mapVals)
	ctxt.IndentedJSON(http.StatusOK, mapVals)
}

/**
 * Get specific key-value pair by key path variable
 */
func getKeyValueByKey(ctxt *gin.Context) {
	jwtMiddleware.ValidateJwt(ctxt)

	key := ctxt.Param("key")
	if keyVal, ok := store[key]; ok {
		log.Printf("Retrieiving Key-Value pair at key: %s\n", key)
		log.Printf("Found Key-Value pair: %s\n", keyVal)
		ctxt.IndentedJSON(http.StatusOK, keyVal)
	} else {
		log.Printf("No Key-Value pair found for key: %s\n", key)
		ctxt.IndentedJSON(http.StatusNotFound, gin.H{"errMessage": helper.MissingKeyMsg})
	}
}

/**
 * Delete the key found at /api/keys/{key}
 * If no Key-Value is stored under that key, don't error out, but return a 204 status
 */
func deleteKeyValueByKey(ctxt *gin.Context) {
	jwtMiddleware.ValidateJwt(ctxt)

	key := ctxt.Param("key")
	if keyVal, ok := store[key]; ok {
		log.Printf("Deleting value at key: %s\n", key)
		ctxt.IndentedJSON(http.StatusOK, keyVal)
		delete(store, key)
	} else {
		log.Printf("No value found for key: %s\n", key)
		ctxt.IndentedJSON(http.StatusNoContent, gin.H{"message": helper.MissingKeyMsg})
	}
}

func main() {
	router := gin.Default()

	// Login endpoint
	router.POST("/api/login", login)

	// KV store
	router.POST("/api/key", storeKeyValue)

	// KV retrieve
	router.GET("/api/key", getAllKeyValuePairs)
	router.GET("/api/key/:key", getKeyValueByKey)

	// KV delete
	router.DELETE("/api/key/:key", deleteKeyValueByKey)

	log.Printf("About to listen on port 5001, go to %s", helper.ServerUrl)

	// NOTE: Long term, a certificate authority to sign the cert used, rather
	// than generating one in a script would be ideal, but that would require
	// client side changes as well to support the x509 certificate pool
	router.RunTLS(helper.Url, "scripts/certificates/serverCert.crt", "scripts/serverCert.key")
}
