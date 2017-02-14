/* The goal of this is going to be pulling a simple webpage off of a website that handles ASL definitions, possible at sompoint implementing multiple manners of managing that. */

package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

// main ... The main function?
func main() {
	//Setting up the default layout.
	address := "http://www.lifeprint.com/asl101/pages-signs/" //Example URL for apple --> http://www.lifeprint.com/asl101/pages-signs/a/apple.htm
	//Getting the url from the getUrl function.
	url, word, err := getURL()
	//Checking the error if exists panic.
	check(err)
	//Parse the url.
	resp, err := http.Get(address + url)
	//Checks the error.
	check(err)
	//Creating the file and writing the data from the HTTP request to it.
	file := writePage(resp, word)
	//Modify the file so that its resources point to the correct directory.
	modifyPage(file)

}

// modifyPage ... Modifies the page so that it points to the correct resources.
func modifyPage(file *os.File) {
	// //Creates the goquery document.
	// doc, err := goquery.NewDocumentFromReader(file)
	// //Checks the doc error.
	// check(err)
	// //Looking for the images and where they're trying to point to.
	// doc.Find("img").Each(func(i int, s *goquery.Selection) {
	// 	str, isExist := s.Attr("src")
	// 	//Checks if that is an attribute of the image.
	// 	if isExist {
	// 		//If it exists it's going to add the static address to the new directory.
	// 		str = breakdownURL(str)
	// 		//Sets the name as the file name.
	// 		name := file.Name()
	// 		//Creating the static path to add to.
	// 		path := filepath.Join(filepath.Join("Resources", name), str)
	// 		s.SetAttr("src", path)
	// 	}
	// })
}

// downloadResources ... Reads through the html looking for specific contents and downloads those resources.
func downloadResources(resp *http.Response, word string) {
	staticAddress := "http://www.lifeprint.com/asl101/"
	//Creating the document from the http response.
	doc, err := goquery.NewDocumentFromResponse(resp)
	//Checks if theres an error.
	check(err)
	//Looking for the image elements in the page.
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		//Checks to see if the string exists, if it does it grabs the value.
		str, isExist := s.Attr("src")
		if isExist {
			//Removes the ../../ and adds the original address.
			str = str[5:]

			//Gets the data by braking down the URL
			data := breakdownURL(str)

			//Creating the url by joining the string with the staticAddress
			url := staticAddress + str
			//Downloads the image using an HTTP request.
			res, err := http.Get(url)
			//Check if there is an error.
			check(err)
			//Write the resource using writeResource()
			writeResource(res, url, data, word)
		}
	})
}

// breakdownURL ... Takes a url and breaks it down by looking for the 3rd "/" and then returning the values after that.
func breakdownURL(str string) string {

	data := ""

	//Itterates through the string recursively finding the / and removes the last bit of data until it find no more and the finally appends data + 1 creating the final title.
	for i, v := range []byte(str) {
		if string(v) == "/" {
			data = str[i+1:]
		}
	}

	return data
}

// writeResource ... Using the response it downloads the source into a folder named resources with another folder being the webpage name.
func writeResource(resp *http.Response, url string, data string, word string) {
	//Makes the required directory for each correct Operating System.
	directory := filepath.Join("Resources", word)
	os.MkdirAll(directory, 0777)
	//Change the directory into the directory.
	err := os.Chdir(directory)
	//Check the error for changing the directory.
	check(err)
	//Create the resource with the correct resource name.
	file, err := os.Create(data)
	//Check the file's error.
	check(err)
	//Copies the datastream into the file.
	io.Copy(file, resp.Body)
	//When done switch back to the original directory.
	//Create a .. directory thing.
	back := filepath.Join("..", "..")
	os.Chdir(back)
	defer resp.Body.Close()

}

// writePage ... Creates the file with the name of word and copies the body of the HTTP request to it.
func writePage(resp *http.Response, word string) *os.File {
	//Creates the file with the name word.
	file, err := os.Create(word + ".html") //Of the type .html
	//Checks the error.
	check(err)
	//Writes to the file by copying the data.
	io.Copy(file, resp.Body)
	//Closes the file when done.
	defer file.Close()
	defer resp.Body.Close()
	//Download resources using goQuery.
	downloadResources(resp, word)

	return file
}

// getURL ... Gets the url by parsing the terminal input using Os.Args[], returns the url.
func getURL() (string, string, error) {
	//Creating the error to be returned.
	err := errors.New("Unable to retrieve URL")
	//Getting the URL by accessing an element of an array from os.Args[] @ position 1.
	word := os.Args[1]

	//Checking if os.Args[1] exists.
	if os.Args[1] != "" {
		err = nil
	}

	//Format the url --> example URL for apple --> http://www.lifeprint.com/asl101/pages-signs/a/apple.htm
	//What we want /a/apple.htm

	/*

	   1. Get the first letter --> Using string slicing.
	   2. Take the entire word --> Already gotten from array slicing the Os.Args[]
	   3. Add .htm

	*/

	//Creating an array to hold all the values which at the end we will pull together to create the final string.
	urlVars := make([]string, 3)

	//Getting the first letter
	urlVars[0] = word[0:1] + "/" //Accessing the first element of the string.
	//We already know the word so just set that directly.
	urlVars[1] = word
	//Just set [2] to .htm
	urlVars[2] = ".htm"

	//Creating the final string for the URL
	url := ""

	//Set the original URL value to what we want.
	//Creats a for each loop and itterates through the entire thing.
	for _, s := range urlVars {
		url += s
	}

	//Returns the values.
	return url, word, err
}

// check ... Checks if there is an error and if there is panic the error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
