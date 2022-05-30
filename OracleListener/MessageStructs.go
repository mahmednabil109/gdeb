package OracleListener

type UnsubscribeMsg struct {
	VM int
}

type OracleMsg struct {
	Key       string `json:"key"`
	Value     []byte `json:"value"`
	Type      int    `json:"type"`
	Timestamp string `json:"timestamp"`
	Error     bool   `json:"error"`
}

type BroadcastMsg struct {
	Key       string `json:"key"`
	Value     []byte `json:"value"`
	Type      int    `json:"type"`
	Timestamp string `json:"timestamp"`
	Index     int
	Error     bool
}

type SubscribeMsg struct {
	VmId          int
	Url           string
	OracleKey     string
	KeyType       int
	Index         int
	BroadcastChan chan *BroadcastMsg
}
