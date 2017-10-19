// Copyright 2015 Brian J. Downs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openweathermap

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ForecastWeatherList holds specific query data
type Forecast5WeatherList struct {
	Dt      int       `json:"dt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
	Wind    Wind      `json:"wind"`
	Speed   float64   `json:"speed"`
	Deg     int       `json:"deg"`
}

// ForecastWeatherData will hold returned data from queries
type Forecast5WeatherData struct {
	COD     string                 `json:"cod"`
	Message float64                `json:"message"`
	City    City                   `json:"city"`
	Cnt     int                    `json:"cnt"`
	List    []Forecast5WeatherList `json:"list"`
	Unit    string
	Lang    string
	Key     string
	*Settings
}

// NewForecast returns a new HistoricalWeatherData pointer with
// the supplied arguments.
func NewForecast5(unit, lang, apiKey string, options ...Option) (*Forecast5WeatherData, error) {
	if !ValidAPIKey(apiKey) {
		return nil, errInvalidKey
	}

	unitChoice := strings.ToUpper(unit)
	langChoice := strings.ToUpper(lang)

	f := &Forecast5WeatherData{
		Settings: NewSettings(),
		Key:      apiKey,
	}

	if ValidDataUnit(unitChoice) {
		f.Unit = DataUnits[unitChoice]
	} else {
		return nil, errUnitUnavailable
	}

	if ValidLangCode(langChoice) {
		f.Lang = langChoice
	} else {
		return nil, errLangUnavailable
	}

	if err := setOptions(f.Settings, options); err != nil {
		return nil, err
	}
	return f, nil
}

// DailyByName will provide a forecast for the location given for the
// number of days given.
func (f *Forecast5WeatherData) DailyByName(location string, days int) error {
	response, err := f.client.Get(fmt.Sprintf(forecast5Base, f.Key, fmt.Sprintf("%s=%s", "q", url.QueryEscape(location)), f.Unit, f.Lang, days))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&f); err != nil {
		return err
	}

	return nil
}

// DailyByCoordinates will provide a forecast for the coordinates ID give
// for the number of days given.
func (f *Forecast5WeatherData) DailyByCoordinates(location *Coordinates, days int) error {
	response, err := f.client.Get(fmt.Sprintf(forecast5Base, f.Key, fmt.Sprintf("lat=%f&lon=%f", location.Latitude, location.Longitude), f.Unit, f.Lang, days))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&f); err != nil {
		return err
	}

	return nil
}

// DailyByID will provide a forecast for the location ID give for the
// number of days given.
func (f *Forecast5WeatherData) DailyByID(id, days int) error {
	response, err := f.client.Get(fmt.Sprintf(forecast5Base, f.Key, fmt.Sprintf("%s=%s", "id", strconv.Itoa(id)), f.Unit, f.Lang, days))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&f); err != nil {
		return err
	}

	return nil
}
