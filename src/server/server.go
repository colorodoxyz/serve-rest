package main

import (
	"log"
	"net/http"

	"github.com/colorodoxyz/serve-rest/src/helper"
	"github.com/gin-gonic/gin"
)

var store = make(map[string]helper.KeyValue)

/*
func keyValueApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	} else if r.Method == "POST" {
	} else if r.Method == "DELETE" {
	} else {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

}
*/

/**
 * Store key-value json struct in the store map
 * {
 *	 "key": "{key}",
 *	 "value": "{value}"
 * }
 */
func storeKeyValue(ctxt *gin.Context) {
	var keyValuePair KeyValue

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
}

/**
 * Get specific key-value pair by key path variable
 */
func getKeyValueByKey(ctxt *gin.Context) {
	key := ctxt.Param("key")
	if keyVal, ok := store[key]; ok {
		log.Printf("Retrieiving Key-Value pair at key: %s\n", key)
		log.printf("Found Key-Value pair: %s\n", keyVal)
		ctxt.IndentedJSON(http.StatusOk, keyVal)
	} else {
		log.Printf("No Key-Value pair found for key: %s\n", key)
		ctxt.IndentedJSON(http.StatusNotFound, gin.H{"errMessage": "Key not found"})
	}
}

func deleteKeyValueBykey(ctxt *gin.Context) {
}

func main() {
	router := gin.Default()

	// KV store
	router.POST("/api/key", storeKeyValue)

	// KV retrieve
	router.GET("/api/key", getAllKeyValuePairs)
	router.GET("/api/key/:key", getKeyValueByKey)

	// KV delete
	router.DELETE("/api/key/:key", deleteKeyValueByKey)

	log.Printf("About to listen on port 5001, go to %s", helper.ServerUrl)

	router.Run(helper.Url)
}
