// Getindexes and Tagstodelete functions and ArryIndexes struct exports to utils package

package utils

import (
	"strings"
)

type ArrIndexes struct {
	Devindex int
	Hotindex int
	Masindex int
	Relindex int
}

//Seperates the sorted arrary according to the tag prefix and returns the last indexes

func Getindexes(original_array []string) ArrIndexes {

	var array_indexes ArrIndexes

	for i := range original_array {
		switch {
		case strings.HasPrefix(original_array[i], "develop_"):
			array_indexes.Devindex = i

		case strings.HasPrefix(original_array[i], "hotfix_"):
			array_indexes.Hotindex = i

		case strings.HasPrefix(original_array[i], "master_"):
			array_indexes.Masindex = i

		case strings.HasPrefix(original_array[i], "release_"):
			array_indexes.Relindex = i
		}

	}

	return array_indexes

}

// creates a new array excluding the tags which need to be keep on the registry and returns that array
func Tagstodelete(input_array []string, keep_tags int) []string {

	var delete_tag_array []string
	array_length := len(input_array)

	if array_length > keep_tags {
		delete_tag_array = input_array[0:(array_length - keep_tags)]
	}

	return delete_tag_array
}
