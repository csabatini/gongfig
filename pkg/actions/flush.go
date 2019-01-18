package actions

import (
	"fmt"
	"bufio"
	"os"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"strings"
	"log"
	"time"
	"encoding/json"
)

func flushAll(adminURL string, authKey string) {
	client := &http.Client{Timeout: Timeout * time.Second}

	// We obtain resources data concurrently and push them to the channel that
	// will be handled by services and routes deleting logic
	flushData := make(chan *resourceAnswer)

	// Collect representation of all resources
	for _, resource := range Apis {
		fullPath := getFullPath(adminURL, []string{resource})

		go getResourceListToChan(client, flushData, fullPath, resource, authKey)

	}

	resourcesNum := len(Apis)
	config := map[string]Data{}

	for {
		resource := <- flushData
		config[resource.resourceName] = resource.config

		resourcesNum--

		if resourcesNum == 0 {
			flushResources(client, adminURL, authKey, config)
			fmt.Println("Done")
			break
		}
	}
}

func flushResources(client *http.Client, url string, authKey string, config map[string]Data) {
	// Firstly we need delete routes and only then services,
	// as routes are nested resources of services
	for _, resourceType := range Apis {
		// In order to not overload the kong, limit concurrent post requests to 10
		reqLimitChan := make(chan bool, 10)

		for _, item := range config[resourceType] {
			reqLimitChan <- true

			// Convert item to resource object for further deleting it from Kong
			var instance ResourceInstance
			mapstructure.Decode(item, &instance)

			go func(instance ResourceInstance){
				defer func() { <-reqLimitChan}()

				// Compose path to routes
				instancePathElements := []string{resourceType, instance.Id}
				instancePath := strings.Join(instancePathElements, "/")
				instanceURL := getFullPath(url, []string{instancePath})
				log.Println("Making delete request ", instanceURL)

				request, _ := http.NewRequest(http.MethodDelete, instanceURL, nil)
				if authKey != "" {
					request.Header.Set("apikey", authKey)
				}

				response, err := client.Do(request)

				if err != nil {
					log.Fatal("Request to Kong admin api failed: ", resourceType)
				}

				if response.StatusCode != 204 {
					// Plugin is deleted automatically when it relies
					// to some service or route id
					if response.StatusCode == 404 && resourceType == PluginsPath {
						log.Println("Plugin is already deleted")
					} else {
						message := Message{}
						json.NewDecoder(response.Body).Decode(&message)
						log.Println(message.Message)

						log.Fatal("Was not able to Delete item ", instance.Id, " ", response.StatusCode)
					}
				}
			}(instance)
		}

		// Wait till all routes deleting is finished
		for i := 0; i < cap(reqLimitChan); i++ {
			reqLimitChan <- true
		}
	}
}

// Flush - main function that is called by CLI in wipe Kong config
func Flush(adminURL string, authKey string) {
	fmt.Println("All services and routes will be deleted from kong, are you sure? Write yes or no:")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	// Delete \n at the end
	answer = answer[0:len(answer)-1]

	if answer== "yes" {
		flushAll(adminURL, authKey)
	} else {
		fmt.Println("Configuration was not flushed")
	}
}
