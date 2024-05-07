package googlescraper

import (
	"VLN-backend/config"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	g "github.com/serpapi/google-search-results-golang"
)

type Image struct {
	ImagePath string
	Width     float64
	Height    float64
}

func ImageSearch(query string) ([]*Image, error) {

	type ImageItem = map[string]interface{}

	serpapi_key := config.GetSerapAPIKey()
	parameter := map[string]string{
		"engine":  "google_images",
		"q":       query,
		"api_key": serpapi_key,
	}

	search := g.NewGoogleSearch(parameter, serpapi_key)
	results, err := search.GetJSON()
	images_results := results["images_results"].([]interface{})
	images := make([]*Image, 0)

	for _, val := range images_results {
		image_item := val.(ImageItem)
		image_url, key_ok := image_item["original"]
		if !key_ok {
			log.Println("Image URL does not exist in: ", image_item["link"])
			continue
		}

		if len(image_url.(string)) > 0 {
			image_width := image_item["original_width"].(float64)
			image_height := image_item["original_height"].(float64)
			images = append(images, &Image{ImagePath: image_url.(string), Height: image_height, Width: image_width})
		}
	}

	return images, err
}

func httpGet(request_url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "*/*")

	response, err := client.Do(req)
	if err != nil {
		// Check if the error is due to a timeout
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("HTTP request timed out")
			return nil, nil
		} else {
			log.Fatal(err)
		}
		return nil, err
	}

	return response, err
}

func downloadImage(image *Image, index int) (*string, error) {
	// Get the data from the URL
	url := image.ImagePath

	response, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Handle Unexecpted image state response
	if response.StatusCode != http.StatusOK {
		log.Printf("UNEXPECTED Response Status for image %s: %s", url, response.Status)
		return nil, err
	}

	// handle SEO backlink
	if strings.Contains(url, "seo") || strings.Contains(url, "crawler") || strings.Contains(url, "media_id") {
		log.Printf("\nSKIPPING URL %s: Likely a backlink SEO URL", url)
		return nil, nil
	}

	// Create the "images" directory if it doesn't exist
	err = os.MkdirAll("images", os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create the file
	filenameParts := strings.Split(url, "?")
	ext := filepath.Ext(filenameParts[0])
	if ext == "" || strings.HasSuffix(ext, "cms") || strings.HasSuffix(ext, "- ") {
		ext = ".jpg"
	}

	cwd, _ := os.Getwd()
	filepath := cwd + "/" + fmt.Sprintf("images/%d%s", index, ext)
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Downloaded image %s at %s\n", filepath, url)
	return &filepath, nil
}

func DownloadImages(images []*Image) ([]*Image, error) {
	//returns a list of updated Images with updated path to downloaded file
	updated_images := make([]*Image, 0)
	for index, image := range images {
		result, err := downloadImage(image, index)
		if err != nil || result == nil {
			log.Println(err, " Occurred while downloading ", image.ImagePath)
			continue
		}
		updated_images = append(updated_images, &Image{
			ImagePath: *result, Height: image.Height, Width: image.Width,
		})
	}

	return updated_images, nil
}