package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShortForecastResponse struct {
	Short            string `json:"short_forecast"`
	Characterization string `json:"characterization"`
}

type nwsForecastResponse struct {
	Temperature   float64 `json:"temperature"`
	ShortForecast string  `json:"shortForecast"`
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func GetShortForecast(lat string, long string) (*ShortForecastResponse, *RequestError) {
	link, err := getForecastLink(lat, long)
	if err != nil {
		return nil, err
	}

	forecast, err := getForecast(link)
	if err != nil {
		return nil, err
	}

	return &ShortForecastResponse{
		Short:            forecast.ShortForecast,
		Characterization: forecast.getCharacterization(),
	}, nil
}

func getForecastLink(lat string, long string) (string, *RequestError) {
	url := fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)

		return "", &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return "", &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var responseMap gin.H
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		log.Println(err)

		return "", &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if statusCode, ok := responseMap["status"].(float64); ok {
		if statusCode == 404 && strings.Contains(responseMap["detail"].(string), "Unable to provide data for requested point") {
			return "", &RequestError{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New(responseMap["detail"].(string)),
			}
		}
	}

	return responseMap["properties"].(map[string]any)["forecast"].(string), nil
}

func getForecast(link string) (*nwsForecastResponse, *RequestError) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var responseMap gin.H
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	periodsData, err := json.Marshal(responseMap["properties"].(map[string]any)["periods"].([]any)[0].(map[string]any))
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var forecastResp nwsForecastResponse
	err = json.Unmarshal(periodsData, &forecastResp)
	if err != nil {
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return &forecastResp, nil
}

func (f *nwsForecastResponse) getCharacterization() string {
	switch {
	case f.Temperature >= 78:
		return "Hot"
	case f.Temperature < 78 && f.Temperature > 57:
		return "Moderate"
	default:
		return "Cold"
	}
}
