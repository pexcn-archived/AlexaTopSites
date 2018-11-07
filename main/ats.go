package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	AtsServiceHost     = "ats.amazonaws.com"
	AtsServiceEndpoint = "ats.us-west-1.amazonaws.com"
	AtsUri             = "/api"
	AtsApiUrl          = "https://" + AtsServiceHost + AtsUri

	AtsServiceName   = "AlexaTopSites"
	AtsServiceRegion = "us-west-1"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Println("Usage:")
		fmt.Println("\t./<application_name> access_key secret_key [country_code] [start_number] [count]")
		return
	}

	accessKey := os.Args[1]
	secretKey := os.Args[2]
	country := os.Args[3]
	start := os.Args[4]
	count := os.Args[5]

	response := SendRequest(accessKey, secretKey, country, start, count)
	fmt.Println(response)
}

func SendRequest(accessKey, secretKey, country, start, count string) string {
	dateTz, query, authorization := CreateHeaders(accessKey, secretKey, country, start, count)
	url := AtsApiUrl + "?" + query

	// create request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Amz-Date", dateTz)
	req.Header.Set("Authorization", authorization)

	// send request
	client := &http.Client{
		// set 5s timeout
		Timeout: time.Duration(5 * time.Second),
	}
	response, _ := client.Do(req)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

func CreateHeaders(accessKey, secretKey, country, start, count string) (string, string, string) {
	// TODO: 优化下这里
	location, _ := time.LoadLocation("UTC")
	now := time.Now().In(location)
	date := now.Format("20060102")
	dateTz := now.Format("20060102T150405Z")

	query := "Action=TopSites&Count=" + count + "&CountryCode=" + country + "&ResponseGroup=Country" + "&Start=" + start
	headers := "host:" + AtsServiceEndpoint + "\n" + "x-amz-date:" + dateTz + "\n"
	signedHeaders := "host;x-amz-date"
	payload := GetSha256("")
	request := "GET" + "\n" + AtsUri + "\n" + query + "\n" + headers + "\n" + signedHeaders + "\n" + payload

	algorithm := "AWS4-HMAC-SHA256"
	scope := date + "/" + AtsServiceRegion + "/" + AtsServiceName + "/aws4_request"
	stringToSign := algorithm + "\n" + dateTz + "\n" + scope + "\n" + GetSha256(request)
	signingKey := GetSignatureKey_2(secretKey, date)
	signature := GetHmacSha256(stringToSign, string(signingKey))

	authorization := algorithm + " " + "Credential=" + accessKey + "/" + scope + ", " + "SignedHeaders=" + signedHeaders + ", " + "Signature=" + signature
	return dateTz, query, authorization
}

func GetSignatureKey(secretKey, date string) string {
	kSecret := "AWS4" + secretKey
	kDate := GetHmacSha256(date, kSecret)
	kRegion := GetHmacSha256(AtsServiceRegion, kDate)
	kService := GetHmacSha256(AtsServiceName, kRegion)
	return GetHmacSha256("aws4_request", kService)
}

func GetSignatureKey_2(secretKey, date string) []byte {
	kSecret := []byte("AWS4" + secretKey)
	kDate := GetHmacSha256_2(date, kSecret)
	kRegion := GetHmacSha256_2(AtsServiceRegion, kDate)
	kService := GetHmacSha256_2(AtsServiceName, kRegion)
	return GetHmacSha256_2("aws4_request", kService)
}
