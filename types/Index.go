package types

type Index struct {
	Type string `json:"type"`
	/*
		可以是Topic/Discussion/Floor
	*/
	/*
		Topic      *Topic      `json:"topic,omitempty"`
		Discussion *Discussion `json:"discussion,omitempty"`
		Floor      *Floor      `json:"floor,omitempty"`
	*/
	Data interface{} `json:"data"`
}
