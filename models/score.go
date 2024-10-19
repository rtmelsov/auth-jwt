package models

import "encoding/xml"

// Структура для тела запроса
type ScoreCardsBody struct {
	GetScoreCards string `json:"GetScoreCards"`
}

type ScoreCardsRequest struct {
	GetScoreCards string `json:"get_score_cards" binding:"required"`
}

// Структуры для XML ответа
type ScoreCardsEnvelopeXml struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    ScoreCardsBodyXml `xml:"Body"`
}

type ScoreCardsBodyXml struct {
	GetScoreCardsResponse GetScoreCardsResponseXml `xml:"GetScoreCardsResponse"`
}

type GetScoreCardsResponseXml struct {
	ReturnData ScoreCardsReturnDataXml `xml:"return"`
}

type ScoreCardsReturnDataXml struct {
	Attributes []ScoreCardsAttributeXml `xml:"attributes"`
	Name       string                   `xml:"name"`
}

type ScoreCardsAttributeXml struct {
	Name string `xml:"name"`
}

type ScoreRequest struct {
	Score struct {
		ScoreCard  string `json:"ScoreCard"`
		Attributes struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Values struct {
				Id    string `json:"id"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"attributes"`
	} `json:"Score"`
}

type ScoreResponseXml struct {
	Envelope EnvelopeScoreXml `json:"Envelope"`
}

type EnvelopeScoreXml struct {
	Body    ScoreBodyXml `json:"Body"`
	XmlnsS  string       `json:"_xmlns:S"`
	PrefixS string       `json:"__prefix"`
}

type ScoreBodyXml struct {
	ScoreResponse ScoreResponseDetailsXml `json:"ScoreResponse"`
	PrefixS       string                  `json:"__prefix"`
}

type ScoreResponseDetailsXml struct {
	Return ReturnDetailsXml `json:"return"`
	Xmlns  string           `json:"_xmlns"`
}

type ReturnDetailsXml struct {
	IdQuery                         string `json:"IdQuery"`
	ErrorCode                       string `json:"ErrorCode"`
	ErrorString                     string `json:"ErrorString"`
	Score                           string `json:"Score"`
	OneYearProbabilityOfDefault     string `json:"OneYearProbabilityOfDefault"`
	RiskGrade                       string `json:"RiskGrade"`
	ScoreByML                       string `json:"ScoreByML"`
	OneYearProbabilityOfDefaultByML string `json:"OneYearProbabilityOfDefaultByML"`
	RiskGradeByML                   string `json:"RiskGradeByML"`
	Causes                          Causes `json:"Causes"`
}

type Causes struct {
	Name      string `json:"name"`
	CauseText string `json:"causeText"`
}
