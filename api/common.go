package api

import (
	"encoding/json"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"strings"
)

func ApiPost(url string, m map[string]string) (utils.ResultParam, error) {
	//token, _ := dao.QueryLoginToken()
	nm := map[string]string{"token": global.Acc.Token}
	if m != nil {
		for k, v := range m {
			nm[k] = v
		}
	}
	b, _ := json.Marshal(nm)
	var newURL string
	if strings.Index(url, "http://") == 0 || strings.Index(url, "https://") == 0 {
		newURL = url
	} else {
		newURL = global.API_SERVICE_DOMAIN + url
	}
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    newURL,
		Body:   b,
	})
	if err != nil {
		return utils.ResultParam{}, err
	}
	//global.Log.Info(string(body))
	return utils.ToResultParam(body), nil
}
