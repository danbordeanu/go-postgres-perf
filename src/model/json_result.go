package model

type JSONSuccessResult struct {
	Code          int         `json:"code" example:"200"`
	Message       string      `json:"message,omitempty" example:"Success"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationId string      `json:"correlation_id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
	Id            string      `json:"id" example:"123-456-789-abc-def"`
}

type JSONFailureResult struct {
	Code          int         `json:"code" example:"400"`
	Data          interface{} `json:"data,omitempty"`
	Error         string      `json:"error,omitempty" example:"There was an error processing the request"`
	Stack         string      `json:"stacktrace,omitempty"`
	Id string      `json:"id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
}
