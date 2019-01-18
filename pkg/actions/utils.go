package actions

import (
	"net/url"
	"net/http"
	"log"
	"encoding/json"
	"bytes"
	"strings"
)

// Data - general interface for storing json body answers
type Data []interface{}

// All items are contained of data property of json answer
type resourceConfig struct {
	Data Data `json:"data"`
}

// Get url and path and return concatenation
// e.g http://localhost:8001, services will return http://localhost:8001/services
func getFullPath(adminURL string, pathElements []string) string {
	uri, _ := url.Parse(adminURL)
	path := strings.Join(append([]string{uri.Path}, pathElements...), "/")
	uri.Path = path
	return uri.String()
}

func getResourceList(client *http.Client, fullPath string, authKey string) resourceConfig {
	request, _ := http.NewRequest("GET", fullPath, nil)
	if authKey != "" {
		request.Header.Set("apikey", authKey)
	}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal("Request to Kong admin failed")
	}

	defer response.Body.Close()

	var body resourceConfig
	json.NewDecoder(response.Body).Decode(&body)

	return body
}

// Get list of resources by http and pass it to the channel where it will handled further
func getResourceListToChan(client *http.Client, writeData chan *resourceAnswer, fullPath string,
	resource string, authKey string) {
	body := getResourceList(client, fullPath, authKey)

	// send only data field for writing in order to write { "service": [items...] } instead of
	// { "service": {"data": [items...] }}
	writeData <- &resourceAnswer{resource, body.Data}
}

func requestNewResource(client *http.Client, resource interface{}, url string, authKey string) (string, error) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(resource)

	// Create services first, as routes are nested resources
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/json")
	if authKey != "" {
		request.Header.Set("apikey", authKey)
	}
	log.Println("Making new resource request", url)

	response, err := client.Do(request)
	defer response.Body.Close()

	if err != nil {
		log.Fatal("Request to Kong admin failed")
		return "", err
	}

	if response.StatusCode != 201 {
		message := Message{}
		json.NewDecoder(response.Body).Decode(&message)

		log.Println(message.Message)
		log.Fatal("Was not able to create resource - ", response.StatusCode, " - ", response.Body)
		return "", err
	}

	createdResource := ResourceInstance{}

	json.NewDecoder(response.Body).Decode(&createdResource)

	return createdResource.Id, nil
}

func addResource(connectionBundle *ConnectionBundle, resource interface{}, Id string, idMap *ConcurrentStringMap) {
	defer func() { <-connectionBundle.ReqLimitChan}()

	id, err := requestNewResource(connectionBundle.Client, resource, connectionBundle.URL, connectionBundle.AuthKey)

	if err != nil {
		log.Fatalf("Failed to create resource, %v\n", err)
	}

	idMap.Add(Id, id)
}

func isJSONString(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil

}
