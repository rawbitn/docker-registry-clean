package main

import (
	"crypto/tls"
	"docker_registry/clean/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/magiconair/properties"
)

type RepoList struct {
	Repositories []string
}

type TagList struct {
	Name string
	Tags []string
}

func main() {

	// global struct variables
	var repolist RepoList
	var tagslist TagList

	// global string array variables
	var dev_arr []string
	var hot_arr []string
	var mas_arr []string
	var rel_arr []string

	// reads application.properties file and assigns environment related properties to variables
	envPropertyValues := properties.MustLoadFile("./application.properties", properties.UTF8)

	registry_url := envPropertyValues.MustGetString("DOCKER_REGISTRY_URL")
	tags_to_keep := envPropertyValues.GetInt("KEEP_TAGS", 5)
	min_tag_array_length := envPropertyValues.GetInt("MIN_NO_OF_TAGS", 10)

	// disables SSL verification globally
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	getRepoListUrl := registry_url + "/v2/_catalog?n=5000"

	// gets the available repository lis from the registry
	repoResponse, err := http.Get(getRepoListUrl)
	if err != nil {
		log.Fatalln(err)
	}

	defer repoResponse.Body.Close()

	reporesponsebody, err := io.ReadAll(repoResponse.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// converts the response body to a json object
	json.Unmarshal([]byte(reporesponsebody), &repolist)

	// checks whether the returned repolist has any content. if not exits from the programme
	if (len(repolist.Repositories) == 0) || (repolist.Repositories == nil) {
		log.Fatalln("Received repository list has no content")
	}

	// loops through the each repository
	for i := range repolist.Repositories {

		dev_index := 0
		hot_index := 0
		mas_index := 0
		rel_index := 0

		dev_arr = nil
		hot_arr = nil
		mas_arr = nil
		rel_arr = nil

		// checks whether the repository name is not null
		if len(repolist.Repositories[i]) > 0 {

			getTagListRequest := registry_url + "/v2/" + repolist.Repositories[i] + "/tags/list"

			// gets the tag list related to the repository
			tagsResponse, err := http.Get(getTagListRequest)
			if err != nil {
				log.Fatalln(err)
			}

			defer tagsResponse.Body.Close()

			tagsResponseBody, err := io.ReadAll(tagsResponse.Body)
			if err != nil {
				log.Fatalln(err)
			}

			// converts the response body to a json object
			json.Unmarshal([]byte(tagsResponseBody), &tagslist)

			// gets the tag list as a string array
			tagListArray := tagslist.Tags

			// checks whether the number of tags is greater than the specific value define in application.properties
			if len(tagListArray) > min_tag_array_length {

				sort.Strings(tagListArray)
				taglistarrayindexes := utils.Getindexes(tagListArray)

				dev_index = taglistarrayindexes.Devindex
				hot_index = taglistarrayindexes.Hotindex
				mas_index = taglistarrayindexes.Masindex
				rel_index = taglistarrayindexes.Relindex

				// creates new arrays according to available branch tags
				switch {
				case (dev_index == 0) && (hot_index == 0) && (mas_index == 0) && (rel_index > 0):
					rel_arr = tagListArray[0 : rel_index+1]

				case (dev_index == 0) && (hot_index == 0) && (mas_index > 0) && (rel_index == 0):
					mas_arr = tagListArray[0 : mas_index+1]

				case (dev_index == 0) && (hot_index == 0) && (mas_index > 0) && (rel_index > 0):
					mas_arr = tagListArray[0 : mas_index+1]
					rel_arr = tagListArray[mas_index+1 : rel_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index == 0) && (rel_index == 0):
					hot_arr = tagListArray[0 : hot_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index == 0) && (rel_index > 0):
					hot_arr = tagListArray[0 : hot_index+1]
					rel_arr = tagListArray[hot_index+1 : rel_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index > 0) && (rel_index == 0):
					hot_arr = tagListArray[0 : hot_index+1]
					mas_arr = tagListArray[hot_index+1 : mas_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index > 0) && (rel_index > 0):
					hot_arr = tagListArray[0 : hot_index+1]
					mas_arr = tagListArray[hot_index+1 : mas_index+1]
					rel_arr = tagListArray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index == 0) && (rel_index == 0):
					dev_arr = tagListArray[0 : dev_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index == 0) && (rel_index > 0):
					dev_arr = tagListArray[0 : dev_index+1]
					rel_arr = tagListArray[dev_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index > 0) && (rel_index == 0):
					dev_arr = tagListArray[0 : dev_index+1]
					mas_arr = tagListArray[dev_index+1 : mas_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index > 0) && (rel_index > 0):
					dev_arr = tagListArray[0 : dev_index+1]
					mas_arr = tagListArray[dev_index+1 : mas_index+1]
					rel_arr = tagListArray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index == 0) && (rel_index == 0):
					dev_arr = tagListArray[0 : dev_index+1]
					hot_arr = tagListArray[dev_index+1 : hot_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index == 0) && (rel_index > 0):
					dev_arr = tagListArray[0 : dev_index+1]
					hot_arr = tagListArray[dev_index+1 : hot_index+1]
					rel_arr = tagListArray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index > 0) && (rel_index == 0):
					dev_arr = tagListArray[0 : dev_index+1]
					hot_arr = tagListArray[dev_index+1 : hot_index+1]
					mas_arr = tagListArray[hot_index+1 : mas_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index > 0) && (rel_index > 0):
					dev_arr = tagListArray[0 : dev_index+1]
					hot_arr = tagListArray[dev_index+1 : hot_index+1]
					mas_arr = tagListArray[hot_index+1 : mas_index+1]
					rel_arr = tagListArray[mas_index+1 : rel_index+1]

				}

				// sorts arrays according to the tag version values
				utils.CustomVersionSort(dev_arr)
				utils.CustomVersionSort(hot_arr)
				utils.CustomVersionSort(mas_arr)
				utils.CustomVersionSort(rel_arr)

				// keeps the tags which need to be deleted
				delete_dev_array := utils.Tagstodelete(dev_arr, tags_to_keep)
				delete_hot_array := utils.Tagstodelete(hot_arr, tags_to_keep)
				delete_mas_array := utils.Tagstodelete(mas_arr, tags_to_keep)
				delete_rel_array := utils.Tagstodelete(rel_arr, tags_to_keep)

				var deleteTagsArray []string

				// creates a single arrary which includes all tags to be deleted
				if delete_dev_array != nil {
					deleteTagsArray = append(deleteTagsArray, delete_dev_array...)
				}
				if delete_hot_array != nil {
					deleteTagsArray = append(deleteTagsArray, delete_hot_array...)
				}
				if delete_mas_array != nil {
					deleteTagsArray = append(deleteTagsArray, delete_mas_array...)
				}
				if delete_rel_array != nil {
					deleteTagsArray = append(deleteTagsArray, delete_rel_array...)
				}

				// prints repositories and tags to be deleted to console
				printline := ""

				if deleteTagsArray != nil {
					fmt.Println(repolist.Repositories[i])
					for i := range deleteTagsArray {
						if i == 0 {
							printline = deleteTagsArray[0]
						} else {
							printline = printline + ", " + deleteTagsArray[i]
						}
					}
					fmt.Println(printline)

				} else {
					fmt.Println("NOTHING TO DELETE IN REPOSITORY: ", repolist.Repositories[i])
				}

				// initializes a http client
				httpclient := &http.Client{}

				for j := range deleteTagsArray {

					getDockerDigestUrl := registry_url + "/v2/" + repolist.Repositories[i] + "/manifests/" + deleteTagsArray[j]

					// gets "Docker-Content-Digest"
					tagdigestequest, err := http.NewRequest("GET", getDockerDigestUrl, nil)
					if err != nil {
						log.Fatalln(err)
					}

					tagdigestequest.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

					tagdigestresponse, err := httpclient.Do(tagdigestequest)
					if err != nil {
						fmt.Println("Error sending HTTP request:", err)
						return
					}

					defer tagdigestresponse.Body.Close()

					if tagdigestresponse.StatusCode == 200 {

						tagdigestresponseheaders := tagdigestresponse.Header

						digestValue := tagdigestresponseheaders["Docker-Content-Digest"][0]

						tagDeleteRequestURL := registry_url + "/v2/" + repolist.Repositories[i] + "/manifests/" + digestValue

						// deletes the manifest identified by the received digest
						tagdeleterequest, err := http.NewRequest("DELETE", tagDeleteRequestURL, nil)
						if err != nil {
							log.Fatalln(err)
						}

						tagdeleterequest.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

						tagdeleteresponse, err := httpclient.Do(tagdeleterequest)
						if err != nil {
							fmt.Println("Error sending HTTP request:", err)
							return
						}

						defer tagdeleteresponse.Body.Close()

						deletestatuscode := tagdeleteresponse.StatusCode

						if deletestatuscode == 202 {
							fmt.Println("Successfully deleted Image: ", repolist.Repositories[i], " Tag :", deleteTagsArray[j])
						} else {
							fmt.Println("Unable to delete Image: ", repolist.Repositories[i], " Tag :", deleteTagsArray[j], "   ", deletestatuscode)
						}

					} else {
						fmt.Println("Unable to get digest for Image: ", repolist.Repositories[i], " Tag :", deleteTagsArray[j], "   ", tagdigestresponse.StatusCode)
					}
				}

			}

		}

	}

}
