package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Artist struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	LocationsJson    string   `json:"locations"`
	ConcertDatesJson string   `json:"concertDates"`
	RelationsJson    string   `json:"relations"`
}

type Locations struct {
	ID               int      `json:"id"`
	Locations        []string `json:"locations"`
	ConcertDatesJson string   `json:"dates"`
}

type Dates struct {
	ID               int      `json:"id"`
	Dates        []string 	  `json:"dates"`
}

type DatesLocations map[string][]string

type Relations struct {
	ID             int            `json:"id"`
	DatesLocations DatesLocations `json:"datesLocations"`
}

type Coordinate struct {
	Location string  `json:"location"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

// FetchArtists fetches and processes the list of artists from the API
func FetchArtists() ([]Artist, error) {
	const url = "https://groupietrackers.herokuapp.com/api/artists"

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the artists slice
	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	return artists, nil
}

// FetchRelations fetches the concert dates and locations for a specific artist
func FetchRelations(artistID int) (DatesLocations, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations for artist %d: %w", artistID, err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the relations struct
	var relations Relations
	if err := json.Unmarshal(body, &relations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relations for artist %d: %w", artistID, err)
	}

	return relations.DatesLocations, nil
}

// FetchLocations fetches the concerts locations for a specific artist
func FetchLocations(artistID int) ([]string, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", artistID)

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations for artist %d: %w", artistID, err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the locations struct
	var locations Locations
	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal locations for artist %d: %w", artistID, err)
	}

	return locations.Locations, nil
}

// FetchDates fetches the concerts dates for a specific artist
func FetchDates(artistID int) ([]string, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", artistID)

	// Perform the HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dates for artist %d: %w", artistID, err)
	}
	defer res.Body.Close()

	// Check for a successful response code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the dates struct
	var dates Dates
	if err := json.Unmarshal(body, &dates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dates for artist %d: %w", artistID, err)
	}

	return dates.Dates, nil
}

// FetchGeocode fetches the latitude and longitude of a given location
func FetchGeocode(location string) (float64, float64, error) {
    // URL for OpenStreetMap Nominatim API
    url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&addressdetails=1&limit=1", url.QueryEscape(location))

    // Perform the HTTP GET request
    res, err := http.Get(url)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to fetch geocode for location '%s': %w", location, err)
    }
    defer res.Body.Close()

    // Check for a successful response code
    if res.StatusCode != http.StatusOK {
        return 0, 0, fmt.Errorf("unexpected status code: %d for location '%s'", res.StatusCode, location)
    }

    // Read the response body
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to read response body: %w", err)
    }

    // Unmarshal the JSON response
    var results []struct {
        Latitude  string `json:"lat"`
        Longitude string `json:"lon"`
    }
    if err := json.Unmarshal(body, &results); err != nil {
        return 0, 0, fmt.Errorf("failed to unmarshal geocode response for location '%s': %w", location, err)
    }

    if len(results) == 0 {
        return 0, 0, fmt.Errorf("no results found for location '%s'", location)
    }

    // Convert latitude and longitude to float64
    lat, _ := strconv.ParseFloat(results[0].Latitude, 64)
    lon, _ := strconv.ParseFloat(results[0].Longitude, 64)

    return lat, lon, nil
}


