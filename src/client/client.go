package main

/**
 * client.go is used in a testing capacity on the server (server.go)
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/colorodoxyz/serve-rest/src/helper"
)

var AdminAccount helper.Account = helper.Account{helper.AdminUser, helper.AdminPassword}

var tokenString string

/**
 * Accesses the login api and gets the JWT Token
 */
func performLogin(cli http.Client) error {

	body, err := json.Marshal(AdminAccount)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, helper.ServerUrl+helper.LoginApi, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	res, err := cli.Do(req)

	var bearerToken helper.BearerToken
	err = json.NewDecoder(res.Body).Decode(&bearerToken)
	res.Body.Close()
	tokenString = "Bearer " + bearerToken.Token
	return nil
}

/**
 * Reads json key-value files from build testingJson directory for test_set and overwrite
 */
func readFile(fileName string) ([]helper.KeyValue, error) {
	file, err := os.Open("testingJsons/" + fileName)
	var KVArray []helper.KeyValue

	if err != nil {
		return KVArray, err
	}
	err = json.NewDecoder(file).Decode(&KVArray)
	return KVArray, nil
}

/**
 * This function is added due to test_delete.json, and is used to read a json file into
 * an array of strings
 */
func readFileStringArray(fileName string) ([]string, error) {
	file, err := os.Open("testingJsons/" + fileName)
	var jsonStrs []string

	if err != nil {
		return jsonStrs, err
	}
	err = json.NewDecoder(file).Decode(&jsonStrs)
	return jsonStrs, err
}

/**
 * Read test_set.json file and POST individual key pairs to the server api
 * Then checks said kv json objects using the /api/keys/{key} GET endpoint
 */
func test_SetAndGet(cli http.Client) error {
	log.Println("Running test_SetAndGet")

	KVArray, err := readFile("test_set.json")

	if err != nil {
		return err
	}

	for _, keyVal := range KVArray {
		err := postKeyVal(keyVal, cli)
		if err != nil {
			return err
		}

		keyPair, err := getKeyVal(keyVal.Key, cli)
		if err != nil {
			return err
		}

		if keyPair != keyVal {
			return fmt.Errorf("Key-Value does not match new Key-Value expected")
		}
	}
	log.Println("test_SetAndGet succeeded")
	return nil
}

/**
 * A test on the overwrite functionality at /api/keys/{key}
 * First checks that the key to be overwritten is there, then checks that the
 * key-value objects are distinct.
 * After the checks, the POST is performed to do the overwrite, and one final GET
 * is used to confirm the change
 */
func test_overwrite(cli http.Client) error {
	log.Println("Running test_overwrite")

	KVArray, err := readFile("test_overwrite.json")

	if err != nil {
		return err
	}

	for _, keyVal := range KVArray {
		keyPair, err := getKeyVal(keyVal.Key, cli)
		// Reason we use if against keyPair first: validate that key
		if len(keyPair.Key) == 0 {
			return fmt.Errorf("No key-value at key (%s) to overwrite", keyVal.Key)
		}
		if err != nil {
			return err
		}

		// Check if distinct kv is present at key
		if keyVal == keyPair {
			return fmt.Errorf("Key-Value stored at key (%s) is same as one used to overwrite", keyVal.Key)
		}

		// Perform overwrite
		err = postKeyVal(keyVal, cli)
		if err != nil {
			return err
		}

		// Check that overwrite succeeded
		keyPair, err = getKeyVal(keyVal.Key, cli)
		if err != nil {
			return err
		}

		if keyPair != keyVal {
			return fmt.Errorf("Key-Value does not match new Key-Value expected")
		}
	}
	log.Println("test_overwrite succeeded")
	return nil
}

func test_delete(cli http.Client) error {
	log.Println("Running test_delete")

	jsonStr, err := readFileStringArray("test_delete.json")

	if err != nil {
		return err
	}

	for _, key := range jsonStr {
		keyPair, err := getKeyVal(key, cli)
		// Reason we use if against keyPair first: validate that key
		if len(keyPair.Key) == 0 {
			return fmt.Errorf("No key-value at key (%s) to delete", key)
		}
		if err != nil {
			return err
		}

		// perform delete operation
		err = deleteKeyVal(key, cli)
		if err != nil {
			return err
		}

		keyPair, err = getKeyVal(key, cli)
		if err != nil {
			return err
		}

		if len(keyPair.Key) != 0 {
			return fmt.Errorf("Key-Value was not deleted")
		}
	}
	log.Println("test_delete succeeded")
	return nil
}

/**
 * Sends a POST to the /api/keys endpoint with a K-V pair to be updated or overwritten
 */
func postKeyVal(keyVal helper.KeyValue, cli http.Client) error {
	body, err := json.Marshal(keyVal)
	req, err := http.NewRequest(http.MethodPost, helper.ServerUrl+helper.KeyValueApi, bytes.NewBuffer(body))

	// Error building request
	if err != nil {
		return err
	}

	req.Header.Add(helper.Auth, tokenString)
	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error: %s", res.StatusCode)
	}

	res.Body.Close()

	return nil
}

/**
 * Sends a GET to the /api/keys/{key} endpoint and returns the K-V pair, if any
 */
func getKeyVal(key string, cli http.Client) (helper.KeyValue, error) {
	var keyPair helper.KeyValue

	req, err := http.NewRequest(http.MethodGet, helper.ServerUrl+helper.KeyValueApi+"/"+key, nil)
	if err != nil {
		return keyPair, err
	}

	req.Header.Add(helper.Auth, tokenString)
	res, err := cli.Do(req)
	if err != nil {
		return keyPair, err
	}

	if res.StatusCode != http.StatusOK {
		// We want to know if a 404 is present
		if res.StatusCode == http.StatusNotFound {
			return keyPair, nil
		}
		return keyPair, fmt.Errorf("Error: %s", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&keyPair)
	res.Body.Close()

	return keyPair, nil
}

/**
 * Sends a DELETE to the /api/keys/{key} endpoint to delete the K-V pair, if any
 * Returns 202 if delete is successful; 204 if no key was found
 */
func deleteKeyVal(key string, cli http.Client) error {
	req, err := http.NewRequest(http.MethodDelete, helper.ServerUrl+helper.KeyValueApi+"/"+key, nil)
	if err != nil {
		return err
	}

	req.Header.Add(helper.Auth, tokenString)
	res, err := cli.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error: %s", res.StatusCode)
	}

	res.Body.Close()

	return nil
}

func main() {
	time.Sleep(5 * time.Second)
	client := http.Client{}

	err := performLogin(client)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = test_SetAndGet(client)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = test_overwrite(client)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = test_delete(client)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
