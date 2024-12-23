package edu_school

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/s21platform/school-service/internal/model"
)

func GetPeers(ctx context.Context, token, campusUuid string, offset, limit int64) ([]string, error) {
	url := fmt.Sprintf("https://edu-api.21-school.ru/services/21-school/api/v1/campuses/%s/participantslimit=%d&offset=%d", campusUuid, limit, offset)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot get req, resp status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}

	var peers model.PeersResponse
	if err := json.Unmarshal(body, &peers); err != nil {
		return nil, fmt.Errorf("cannot parse response body: %v", err)
	}

	return peers.Peers, nil
}
