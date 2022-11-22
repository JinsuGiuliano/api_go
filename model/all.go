package model

type HTTPErrorCode int

type HTTPError struct {
	Code  HTTPErrorCode
	Error string
}

type Msg struct {
	Text  string
	Token string `json:",omitempty"`
}

type LoginStatus int
