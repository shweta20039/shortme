package api


type Version struct {
	Version string `json:"version"`
}

type Err struct {
	Msg string `json:"msg"`
}

type ShortReq struct {
	LongURL	string `json:"longURL"`
}

type ShortResp struct {
	ShortURL	string	`json:"shortURL"`
}

type ExpandReq struct {
	ShortURL	string	`json:"shortURL"`
}

type ExpandResp struct {
	LongURL 	string	`json:"longURL"`
}

