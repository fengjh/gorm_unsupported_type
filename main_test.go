package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initDB()
	migrate()

	retCode := m.Run()

	os.Exit(retCode)
}

func TestCreateSurvey(t *testing.T) {
	var answersData Answers

	err := json.Unmarshal([]byte(`[
		{"question":"What colour is your hair","answers":["increasingly grey"]},
		{"question":"What is up?","answers":["clouds"]},
		{"question":"What are your hobbies?","answers":["skiing","music production","karaoke","karate"]}
	]`), &answersData)

	assertNoErr(t, err)

	survey := Survey{Answers: answersData}

	err = DB.Create(&survey).Error
	assertNoErr(t, err)

	assertNoErr(t, DB.Find(&survey, survey.ID).Error)

	if len(survey.Answers) != 3 {
		t.Fatalf("Expected survey of 3 answers, got %v", survey)
	}

	for i, answerData := range answersData {
		if !isEqualAsJSONString(t, answerData, survey.Answers[i]) {
			t.Fatalf("Unexpected answer %v, wanted %v", survey.Answers[i], answerData)
		}
	}
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func toJSONString(t *testing.T, value interface{}) string {
	jsonStr, err := json.Marshal(value)
	assertNoErr(t, err)

	return string(jsonStr)
}

func isEqualAsJSONString(t *testing.T, expectedValue interface{}, value interface{}) bool {
	expectedJSONStr := toJSONString(t, expectedValue)
	jsonStr := toJSONString(t, value)

	if jsonStr != expectedJSONStr {
		t.Logf("expected JSON string: \n%v\n, but got: \n%v\n", expectedJSONStr, jsonStr)
		return false
	}
	return true
}
