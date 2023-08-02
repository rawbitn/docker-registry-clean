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

	//reporesponsebody := `{"repositories":["abc-bank-open-api-abc-bank-open-api","account-mangement","account-payment","account-paymentinit","account-testmanagement","account-testpaymntinit","adl_mw-testupgrade-testu-testms","adlapp-ms2"]}`
	//tag_response := `{"name":"wow-superapp-notification-policy-execution-service","tags":["develop_1.0.10_43","develop_1.1.0_50","develop_1.0.10_33","release_1.1.1_19","release_1.1.1_24"]}`

	var repolist RepoList
	var tagslist TagList

	var dev_arr []string
	var hot_arr []string
	var mas_arr []string
	var rel_arr []string

	propertyvalues := properties.MustLoadFile("./application.properties", properties.UTF8)

	registry_url := propertyvalues.MustGetString("DOCKER_REGISTRY_URL")
	tags_to_keep := propertyvalues.GetInt("KEEP_TAGS", 5)

	///*
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	getRepoListUrl := registry_url + "/v2/_catalog?n=5000"

	reporesponse, err := http.Get(getRepoListUrl)
	if err != nil {
		log.Fatalln(err)
	}

	defer reporesponse.Body.Close()

	//We Read the response body on the line below.
	reporesponsebody, err := io.ReadAll(reporesponse.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//*/

	//Convert the body to type string
	json.Unmarshal([]byte(reporesponsebody), &repolist)

	lengthrepolist := len(repolist.Repositories)

	if (lengthrepolist == 0) || (repolist.Repositories == nil) {
		log.Fatalln("Received repository list has no content")
	}

	for i := range repolist.Repositories {

		dev_index := 0
		hot_index := 0
		mas_index := 0
		rel_index := 0

		dev_arr = nil
		hot_arr = nil
		mas_arr = nil
		rel_arr = nil

		if len(repolist.Repositories[i]) > 0 {

			getTagListRequest := registry_url + "/v2/" + repolist.Repositories[i] + "/tags/list"

			tagsresponse, err := http.Get(getTagListRequest)
			if err != nil {
				log.Fatalln(err)
			}

			defer tagsresponse.Body.Close()

			//We Read the response body on the line below.
			tagsresponsebody, err := io.ReadAll(tagsresponse.Body)
			if err != nil {
				log.Fatalln(err)
			}

			json.Unmarshal([]byte(tagsresponsebody), &tagslist)

			taglistarray := tagslist.Tags

			if len(taglistarray) > 10 {
				sort.Strings(taglistarray)

				taglistarrayindexes := utils.Getindexes(taglistarray)

				dev_index = taglistarrayindexes.Devindex
				hot_index = taglistarrayindexes.Hotindex
				mas_index = taglistarrayindexes.Masindex
				rel_index = taglistarrayindexes.Relindex

				switch {
				case (dev_index == 0) && (hot_index == 0) && (mas_index == 0) && (rel_index > 0):
					rel_arr = taglistarray[0 : rel_index+1]

				case (dev_index == 0) && (hot_index == 0) && (mas_index > 0) && (rel_index == 0):
					mas_arr = taglistarray[0 : mas_index+1]

				case (dev_index == 0) && (hot_index == 0) && (mas_index > 0) && (rel_index > 0):
					mas_arr = taglistarray[0 : mas_index+1]
					rel_arr = taglistarray[mas_index+1 : rel_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index == 0) && (rel_index == 0):
					hot_arr = taglistarray[0 : hot_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index == 0) && (rel_index > 0):
					hot_arr = taglistarray[0 : hot_index+1]
					rel_arr = taglistarray[hot_index+1 : rel_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index > 0) && (rel_index == 0):
					hot_arr = taglistarray[0 : hot_index+1]
					mas_arr = taglistarray[hot_index+1 : mas_index+1]

				case (dev_index == 0) && (hot_index > 0) && (mas_index > 0) && (rel_index > 0):
					hot_arr = taglistarray[0 : hot_index+1]
					mas_arr = taglistarray[hot_index+1 : mas_index+1]
					rel_arr = taglistarray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index == 0) && (rel_index == 0):
					dev_arr = taglistarray[0 : dev_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index == 0) && (rel_index > 0):
					dev_arr = taglistarray[0 : dev_index+1]
					rel_arr = taglistarray[dev_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index > 0) && (rel_index == 0):
					dev_arr = taglistarray[0 : dev_index+1]
					mas_arr = taglistarray[dev_index+1 : mas_index+1]

				case (dev_index > 0) && (hot_index == 0) && (mas_index > 0) && (rel_index > 0):
					dev_arr = taglistarray[0 : dev_index+1]
					mas_arr = taglistarray[dev_index+1 : mas_index+1]
					rel_arr = taglistarray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index == 0) && (rel_index == 0):
					dev_arr = taglistarray[0 : dev_index+1]
					hot_arr = taglistarray[dev_index+1 : hot_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index == 0) && (rel_index > 0):
					dev_arr = taglistarray[0 : dev_index+1]
					hot_arr = taglistarray[dev_index+1 : hot_index+1]
					rel_arr = taglistarray[mas_index+1 : rel_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index > 0) && (rel_index == 0):
					dev_arr = taglistarray[0 : dev_index+1]
					hot_arr = taglistarray[dev_index+1 : hot_index+1]
					mas_arr = taglistarray[hot_index+1 : mas_index+1]

				case (dev_index > 0) && (hot_index > 0) && (mas_index > 0) && (rel_index > 0):
					dev_arr = taglistarray[0 : dev_index+1]
					hot_arr = taglistarray[dev_index+1 : hot_index+1]
					mas_arr = taglistarray[hot_index+1 : mas_index+1]
					rel_arr = taglistarray[mas_index+1 : rel_index+1]

				}

				utils.CustomVersionSort(dev_arr)
				utils.CustomVersionSort(hot_arr)
				utils.CustomVersionSort(mas_arr)
				utils.CustomVersionSort(rel_arr)

				delete_dev_array := utils.Tagstodelete(dev_arr, tags_to_keep)
				delete_hot_array := utils.Tagstodelete(hot_arr, tags_to_keep)
				delete_mas_array := utils.Tagstodelete(mas_arr, tags_to_keep)
				delete_rel_array := utils.Tagstodelete(rel_arr, tags_to_keep)

				var delete_tags_array []string

				if delete_dev_array != nil {
					delete_tags_array = append(delete_tags_array, delete_dev_array...)
				}
				if delete_hot_array != nil {
					delete_tags_array = append(delete_tags_array, delete_hot_array...)
				}
				if delete_mas_array != nil {
					delete_tags_array = append(delete_tags_array, delete_mas_array...)
				}
				if delete_rel_array != nil {
					delete_tags_array = append(delete_tags_array, delete_rel_array...)
				}

				printline := ""

				if delete_tags_array != nil {
					printline = printline + repolist.Repositories[i]
					for i := range delete_tags_array {
						printline = printline + ", " + delete_tags_array[i]
					}
					fmt.Println(printline)

				} else {
					fmt.Println(repolist.Repositories[i], " : NOTHING TO DELETE")
				}

				httpclient := &http.Client{}

				getDockerDigestUrl := registry_url + "/v2/" + repolist.Repositories[i] + "/manifests/" + "develop_0.0.1_4"

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

				tagdigestresponseheaders := tagdigestresponse.Header

				digestValue := tagdigestresponseheaders["Docker-Content-Digest"][0]

				tagDeleteRequestURL := registry_url + "/v2/" + repolist.Repositories[i] + "/manifests/" + digestValue

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
					fmt.Println("Sucessfully deleted")
				} else {
					fmt.Println("Error occurd while deleting", deletestatuscode)
				}

			}

		}

	}

}
