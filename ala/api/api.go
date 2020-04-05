package api

import (
	"encoding/json"
	"github.com/GreatGodApollo/ala/models"
	"github.com/go-resty/resty/v2"
)

func GetAll(c *resty.Client, baseurl, username, password, product string) (interface{}, error) {
	resp, err := c.R().
		SetBasicAuth(username, password).
		SetHeader("Accept", "application/json").
		Get(baseurl + "/api/v1/all/" + product)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode()/100 == 2 {
		var respBody models.Licenses
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	} else {
		var respBody models.BasicResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	}
}

func CreateLicense(c *resty.Client, baseurl, username, password, email, product string) (interface{}, error) {
	resp, err := c.R().
		SetBasicAuth(username, password).
		SetBody(models.LicenseRequest{Email: email, Product: product}).
		SetHeader("Accept", "application/json").
		Post(baseurl + "/api/v1/create")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode()/100 == 2 {
		var respBody models.LicenseResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	} else {
		var respBody models.BasicResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	}
}

func GetSpecific(c *resty.Client, baseurl, username, password, key string) (interface{}, error) {
	resp, err := c.R().
		SetBasicAuth(username, password).
		SetHeader("Accept", "application/json").
		SetBody(models.BasicRequest{Key: key}).
		Post(baseurl + "/api/v1/specific")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode()/100 == 2 {
		var respBody models.License
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	} else {
		var respBody models.BasicResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	}
}

func InvalidateLicense(c *resty.Client, baseurl, username, password, key string) (interface{}, error) {
	var result models.BasicResponse
	resp, err := c.R().
		SetBasicAuth(username, password).
		SetBody(models.BasicRequest{Key: key}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Post(baseurl + "/api/v1/invalidate")

	if err != nil {
		return models.BasicResponse{}, err
	}

	if resp.StatusCode()/100 == 2 {
		var respBody models.LicenseResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	} else {
		var respBody models.BasicResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return nil, err
		}
		return respBody, nil
	}

}

func CheckValidity(c *resty.Client, baseurl, key, product string) bool {
	resp, err := c.R().
		SetHeader("Accept", "application/json").
		SetBody(models.CheckRequest{Key: key, Product: product}).
		Post(baseurl + "/license/check")

	if err != nil {
		return false
	}

	if resp.StatusCode()/100 == 2 {
		var respBody models.LicenseResponse
		err = json.Unmarshal(resp.Body(), &respBody)
		if err != nil {
			return false
		}
		if respBody.Status == "valid" {
			return true
		}
		return false
	} else {
		return false
	}
}
