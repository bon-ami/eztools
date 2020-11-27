package eztools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rkl-/digest"
)

const (
	AUTH_NONE = iota
	AUTH_PLAIN
	AUTH_BASIC
	AUTH_DIGEST
	METHOD_GET  = "GET"
	METHOD_PUT  = "PUT"
	METHOD_POST = "POST"
)

type AuthInfo struct {
	Type       int
	User, Pass string
}

func digestAuth(user, pass string, req *http.Request) (*http.Response, error) {
	return digest.NewTransport(user, pass).RoundTrip(req)
}

func genReq(method, url string, bodyReq io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, bodyReq)
	if err != nil {
		if Debugging {
			ShowStrln("failed to create " + method)
		}
		return
	}
	if bodyReq != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		/*if Debugging && Verbose > 2 {
			ShowSthln("body ", bodyReq)
		}*/
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	return
}

func parseResp(resp *http.Response, magic []byte) (bodyMap interface{}, errNo int, err error) {
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		if Debugging {
			ShowStr(strconv.Itoa(resp.StatusCode))
			ShowStrln(" failure response")
		}
		var b []byte
		if resp.ContentLength > 0 {
			b = make([]byte, resp.ContentLength)
			_, err = resp.Body.Read(b)
		}
		return b, resp.StatusCode, errors.New(resp.Status)
	}

	if Debugging && Verbose > 2 {
		ShowStrln("resp code=" + strconv.Itoa(resp.StatusCode))
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		err = errors.New(ct)
		return
	}
	if cl := resp.Header.Get("Content-Length"); cl == "0" {
		if Debugging && Verbose > 0 {
			ShowStrln("no body in response")
		}
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || body == nil || len(body) < 1 {
		if Debugging {
			ShowStrln("failed to read body")
		}
		return
	}
	if len(magic) > 0 {
		if Debugging && Verbose > 1 {
			ShowStrln("stripping magic")
		}
		if bytes.HasPrefix(body, magic) {
			body = bytes.TrimLeft(bytes.TrimPrefix(body, magic), "\n\r")
		} else {
			err = errors.New("Magic not matched")
			return
		}
	}
	//var bodyMp map[string]interface{}
	//bodyMap = &bodyMp
	err = json.Unmarshal(body, &bodyMap)
	//fullResponse := successResponse{
	//Data: v,
	//}
	//if err = json.NewDecoder(resp.Body).Decode(&fullResponse); err != nil {
	//return
	//}
	if Debugging && Verbose > 2 {
		ShowSthln(bodyMap)
	}
	return
}

// RestGetPlain sends with timeout, Restful API request and returns the result.
func restGetPlain(method, url, pass string, to time.Duration, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	if Debugging {
		ShowStrln("REST 2 " + url)
	}
	cli := &http.Client{
		Timeout: to,
	}
	req, err := genReq(method, url, bodyReq)
	if err != nil {
		return
	}
	if len(pass) > 0 {
		req.Header.Set("authorization", "Basic "+pass)
	}
	resp, err := cli.Do(req)
	if err != nil {
		if Debugging {
			ShowStrln("failed to send/get")
		}
		return
	}

	return parseResp(resp, magic)
}

// RestGetBasic sends with timeout, Restful API request and returns the result.
func restGetBasic(method, url, user, pass string, to time.Duration, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	if Debugging {
		ShowStrln("REST 2 " + url)
	}
	cli := &http.Client{
		Timeout: to,
	}
	req, err := genReq(method, url, bodyReq)
	if err != nil {
		return
	}
	if len(pass) > 0 {
		req.SetBasicAuth(user, pass)
	}
	resp, err := cli.Do(req)
	if err != nil {
		if Debugging {
			ShowStrln("failed to send/get")
		}
		return
	}

	return parseResp(resp, magic)
}

const defRestGetTO = 60 * time.Second

// RestGet sends Restful API request and returns the result.
func restGetDigest(method, url, user, pass string, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	if Debugging {
		ShowStrln("REST 2 " + url)
	}
	req, err := genReq(method, url, bodyReq)
	if err != nil {
		return
	}
	resp, err := digestAuth(user, pass, req)
	if err != nil {
		if Debugging {
			ShowStrln("failed to send/get")
		}
		return
	}

	return parseResp(resp, magic)
}

// RestGet sends Restful API request and returns the result.
func RestGet(url string, authInfo AuthInfo, bodyReq io.Reader) (bodyMap interface{}, errNo int, err error) {
	return RestGetWtMagic(url, authInfo, bodyReq, nil)
}

func RestGetOrPostWtMagic(method, url string, authInfo AuthInfo, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	switch authInfo.Type {
	case AUTH_DIGEST:
		return restGetDigest(method, url, authInfo.User, authInfo.Pass, bodyReq, magic)
	case AUTH_NONE:
		return restGetPlain(method, url, "", defRestGetTO, bodyReq, magic)
	case AUTH_BASIC:
		return restGetBasic(method, url, authInfo.User, authInfo.Pass, defRestGetTO, bodyReq, magic)
	}
	// AUTH_PLAIN, default
	return restGetPlain(method, url, authInfo.Pass, defRestGetTO, bodyReq, magic)
}

/* RestPostWtMagic sends Restful API request and returns the result.
A magic string can be listed to be stripped in the beginning of
the result. If magic string is assigned but not found in the result,
the result will not be parsed. */
func RestPostWtMagic(url string, authInfo AuthInfo, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	return RestGetOrPostWtMagic(METHOD_POST, url, authInfo, bodyReq, magic)
}

/* RestGetWtMagic sends Restful API request and returns the result.
A magic string can be listed to be stripped in the beginning of
the result. If magic string is assigned but not found in the result,
the result will not be parsed. */
func RestGetWtMagic(url string, authInfo AuthInfo, bodyReq io.Reader, magic []byte) (bodyMap interface{}, errNo int, err error) {
	return RestGetOrPostWtMagic(METHOD_GET, url, authInfo, bodyReq, magic)
}

/* RangeStrMap iterate through map[string]interface{} obj, calling fun for
each element recursively. When fun returns true, it stops.
false is returned if no element found. */
func RangeStrMap(obj interface{}, fun func(k string, v interface{}) bool) bool {
	//if the argument is not a map, ignore it
	mobj, ok := obj.(map[string]interface{})
	if !ok {
		return false
	}

	for k, v := range mobj {
		//key match, return value
		if fun(k, v) {
			return true
		}

		//if the value is a map, search recursively
		if m, ok := v.(map[string]interface{}); ok {
			if RangeStrMap(m, fun) {
				return true
			}
		}
		//if the value is an array, search recursively
		//from each element
		if va, ok := v.([]interface{}); ok {
			for _, a := range va {
				if RangeStrMap(a, fun) {
					return true
				}
			}
		}
	}

	//element not found
	return false
}

/* FindStrMap find string key in map[string]interface{} obj,
returning the value and true or nil and false. */
func FindStrMap(obj interface{}, key string) (interface{}, bool) {
	//if the argument is not a map, ignore it
	mobj, ok := obj.(map[string]interface{})
	if !ok {
		return nil, false
	}

	for k, v := range mobj {
		//key match, return value
		if k == key {
			return v, true
		}

		//if the value is a map, search recursively
		if m, ok := v.(map[string]interface{}); ok {
			if res, ok := FindStrMap(m, key); ok {
				return res, true
			}
		}
		//if the value is an array, search recursively
		//from each element
		if va, ok := v.([]interface{}); ok {
			for _, a := range va {
				if res, ok := FindStrMap(a, key); ok {
					return res, true
				}
			}
		}
	}

	//element not found
	return nil, false
}
