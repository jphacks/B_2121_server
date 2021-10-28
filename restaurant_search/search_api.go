package restaurant_search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jphacks/B_2121_server/models"
	"golang.org/x/xerrors"
)

const hotpepperBaseUrl = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"

type Restaurant struct {
	Id       string
	Name     string
	Location models.Location
	ImageUrl string
	PageUrl  string
	Address  string
}

type SearchApi interface {
	Search(keyword string, location *models.Location, count int) (*[]Restaurant, error)
	SearchNext(keyword string, location *models.Location, startCount int, count int) (*[]Restaurant, error)
}

func NewSearchApi(apiKey string) SearchApi {
	return &hotpepperSearch{token: apiKey}
}

type hotpepperSearch struct {
	token string
}

func (h hotpepperSearch) Search(keyword string, location *models.Location, count int) (*[]Restaurant, error) {
	r, err := h.SearchNext(keyword, location, 0, count)
	if err != nil {
		return nil, xerrors.Errorf("API request failed: %w", err)
	}
	return r, nil
}

func (h hotpepperSearch) SearchNext(keyword string, location *models.Location, startCount int, count int) (*[]Restaurant, error) {
	u, err := url.Parse(hotpepperBaseUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse base url: %w", err)
	}
	q := u.Query()
	q.Set("key", h.token)
	q.Set("keyword", keyword)
	q.Set("start", fmt.Sprintf("%d", startCount))
	q.Set("count", fmt.Sprintf("%d", count))
	q.Set("format", "json")
	if location != nil {
		q.Set("lat", fmt.Sprintf("%.10f", location.Latitude))
		q.Set("lng", fmt.Sprintf("%.10f", location.Longitude))
		q.Set("range", "5")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to retrieve response: %w", err)
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			println(err)
		}
	}()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to read response body: %w", err)
	}
	var jsonData hotpepperResponseJson
	err = json.Unmarshal(respData, &jsonData)
	if err != nil {
		return nil, xerrors.Errorf("failed to deserialize json: %w", err)
	}

	ret := make([]Restaurant, 0)
	for _, restaurant := range jsonData.Results.Shop {
		ret = append(ret, Restaurant{
			Id:   restaurant.Id,
			Name: restaurant.Name,
			Location: models.Location{
				Latitude:  restaurant.Lat,
				Longitude: restaurant.Lng,
			},
			ImageUrl: restaurant.Photo.Mobile.L,
			PageUrl:  restaurant.Urls.Pc,
			Address:  restaurant.Address,
		})
	}

	return &ret, nil
}
