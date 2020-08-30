package rate

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dougkirkley/usps"
)

// Request handles rate request
type Request struct {
	XMLName        xml.Name `xml:Package" json:"package"`
	ID             string   `xml:"ID,attr" json:"id"`
	Service        string   `xml:"Service" json:"service"`
	ZipOrigination string   `xml:"ZipOrigination" json:"zip_origin"`
	ZipDestination string   `xml:"ZipDestination" json:"zip_dest"`
	Pounds         string   `xml:"Pounds" json:"pounds"`
	Ounces         string   `xml:"Ounces" json:"ounces"`
	Container      string   `xml:"Container" json:"container"`
}

// Response handles Response
type Response struct {
	XMLName        xml.Name `xml:RateV4Response"`
	Package        Package  `xml:"Package"`
}

// Package handles xml Package data
type Package struct {
	XMLName        xml.Name `xml:Package"`
	ID             string   `xml:"ID,attr"`
	ZipOrigination string   `xml:"ZipOrigination"`
	ZipDestination string   `xml:"ZipDestination"`
	Pounds         string   `xml:"Pounds"`
	Ounces         string   `xml:"Ounces"`
	Container      string   `xml:"Container"`
	Postage        Postage  `xml:"Postage"`
}

// Postage handles xml Postage data
type Postage struct {
	XMLName     xml.Name `xml:"Postage"`
	ClassID     string   `xml:"CLASSID,attr"`
	MailService string   `xml:"MailService"`
	Rate        string   `xml:"Rate"`
}

// Interface is implemented by Client
type Interface interface {
	Calculate(rates []Request) (string, error)
}

// Client is a USPS API client.
type Client struct {
	user   string
	url    string
	client *http.Client
}

// NewRate returns a USPS API rate client.
func NewRate(user usps.USPS) Interface {
	return &Client{
		user: user.Username,
		url:    "https://secure.shippingapis.com/RateV4API.dll",
		client: http.DefaultClient,
	}
}

// Calculate retuns shipping rate
func (c *Client) Calculate(rates []Request) (string, error) {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return "", err
	}
	
	var packages string
	for _, rate := range rates {
		xmlOut, err := xml.Marshal(rate)
	    if err != nil {
		    return "", err
		}
		packages += string(xmlOut)
	}

	// Construct the URL encoded query
	query := `<RateV4Request USERID=%q><Revision>%s</Revision>%s</RateV4Request>`
	req.URL.RawQuery = fmt.Sprintf("API=RateV4&XML=%s", url.QueryEscape(fmt.Sprintf(query, c.user, "0", packages)))

	// Get the request
	resp, err := c.client.Do(req)
	if err != nil {
	    return "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var rateResp Response
	err = xml.Unmarshal(body, &rateResp)
	if err != nil {
		return "", err
	}
	
	return rateResp.Package.Postage.Rate, nil
}
