package client

import "net/url"

func EncodeURLForm(values map[string]string) string {
	formData := url.Values{}
	for k, v := range values {
		formData.Set(k, v)
	}
	return formData.Encode()
}

func EncodeQuery(values map[string]string) string {
	formData := url.Values{}
	for k, v := range values {
		formData.Set(k, v)
	}
	return "?" + formData.Encode()
}
