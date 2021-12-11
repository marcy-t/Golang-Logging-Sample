package domain

type SampleRequest struct {
	Host string
}

type SamplePingListResponse struct {
	SamplePingList []*SamplePing `json:"sample_ping_list"`
}

type SamplePing struct {
	ID      int    `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SamplePingResponse struct {
	ID      int    `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
