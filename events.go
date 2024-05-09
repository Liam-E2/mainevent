package main

type Question struct{
	Id int `json:"question_id"`
	Text string `json:"question_text"`
	Answer string `json:"question_answer"`
}
