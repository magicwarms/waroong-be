package config

type DefaultResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  string      `json:"errors"`
}

type ValidationDefaultResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

type ResponseField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}
