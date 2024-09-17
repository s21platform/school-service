package edu_school

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type ErrorOfGettingTribes struct {
	Status        int    `json:"status"`
	ExceptionUuid string `json:"exceptionUUID"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

type Tribe struct {
	Id   int32  `json:"coalitionId"`
	Name string `json:"name"`
}

type TribesResponse struct {
	Tribes []Tribe `json:"coalitions"`
}

func GetTribesOfCampus(campusUuid, token string) (*TribesResponse, error) {
	limit, offset := 50, 0
	requestUrl := fmt.Sprintf("https://edu-api.21-school.ru/services/21-school/api/v1/campuses/%s/coalitions?limit=%d&offset=%d", campusUuid, limit, offset)

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var tribesError ErrorOfGettingTribes
		err = json.Unmarshal(body, &tribesError)
		if err != nil {
			return nil, err
		}
		return nil, status.Errorf(codes.Unknown, "Error getting list of tribes. Status code: %d: %s: %s", tribesError.Status, tribesError.Code, tribesError.Message)
	}

	var result TribesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
