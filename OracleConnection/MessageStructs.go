package OracleConnection

type UnsubscribeMsg struct {
	VM int
}

type ResponseMsg struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type BroadcastMsg struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Index     int
	Error     error
}

type SubscribeMsg struct {
	VM            int
	Url           string
	OracleKey     string
	Index         int
	BroadcastChan chan *BroadcastMsg
}
