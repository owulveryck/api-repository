package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

type returnType struct {
	ID string
}

func postRequest(jsonStr []byte) (time.Duration, int, returnType, error) {
	var ret returnType
	var retCode int
	var elapsed time.Duration
	req, err := http.NewRequest("POST", testURL+"/products/models", bytes.NewBuffer(jsonStr))
	if err != nil {
		return elapsed, retCode, ret, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	now := time.Now()
	resp, err := client.Do(req)
	elapsed = time.Since(now)
	if err != nil {
		return elapsed, retCode, ret, err
	}
	retCode = resp.StatusCode
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ret)
	return elapsed, retCode, ret, err
}

func TestModelCreate(t *testing.T) {
	simpleModel := Model{
		Code:           uuid.New().String(),
		Name:           "driller",
		AttributesCode: []string{"a1234", "a2345"},
	}
	jsonStr, err := json.Marshal([]Model{simpleModel})
	if err != nil {
		t.Fatal(err)
	}
	d, code, _, err := postRequest(jsonStr)
	t.Log(d)
	if err != nil {
		t.Fatal(err)
	}
	if code != http.StatusOK {
		t.Fatal("Non 200 status: ", code)
	}
}

func waitForSaved(id string) error {
	req, err := http.NewRequest("GET", testURL+"/jobs/"+id, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	for {
		time.Sleep(200 * time.Millisecond)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusMultiStatus {
			return fmt.Errorf("Expecting 207, got %v", resp.StatusCode)
		}
		var t transaction
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&t)
		if err != nil {
			return err
		}
		exit := true
		for _, e := range t.Elements {
			if e.Status != http.StatusOK && e.Status != http.StatusInternalServerError {
				exit = false
				continue
			}
		}
		if exit {
			return nil
		}
	}
}

func TestModelCreate_100(t *testing.T) {
	numberElements := 100
	models := make([]Model, numberElements)
	for i := 0; i < numberElements; i++ {
		models[i] = Model{
			Code:           uuid.New().String(),
			Name:           "driller",
			AttributesCode: []string{"a1234", "a2345"},
		}
	}
	jsonStr, err := json.Marshal(models)
	if err != nil {
		t.Fatal(err)
	}
	d, code, ret, err := postRequest(jsonStr)
	t.Log(d)
	if err != nil {
		t.Fatal(err)
	}
	if code != http.StatusOK {
		t.Fatal("Non 200 status: ", code)
	}
	t.Log("Waiting for elements to be saved")
	err = waitForSaved(ret.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkModelCreate(b *testing.B) {
	random := func() (string, error) {
		numberElements := 100
		models := make([]Model, numberElements)
		for i := 0; i < numberElements; i++ {
			models[i] = Model{
				Code:           uuid.New().String(),
				Name:           "driller",
				AttributesCode: []string{"a1234", "a2345"},
			}
		}
		jsonStr, err := json.Marshal(models)
		if err != nil {
			return "", err
		}
		_, _, ret, err := postRequest(jsonStr)
		if err != nil {
			return "", err
		}
		return ret.ID, nil
	}

	var ids []string
	for n := 0; n < b.N; n++ {
		id, err := random()
		if err != nil {
			b.Fatal(err)
		}
		ids = append(ids, id)
	}
	/*
		for _, id := range ids {
			err := waitForSaved(id)
			if err != nil {
				b.Fatal(err)
			}
		}
	*/
}
