package chicustom

type (
	responseDesc struct {
		ID string `json:"id"`
		EN string `json:"en"`
	}

	meta struct {
		Version        string `json:"version"`
		Status         string `json:"status"`
		ApiEnvironment string `json:"api_environment"`
	}

	success struct {
		ResponseCode string       `json:"response_code"`
		ResponseDesc responseDesc `json:"response_desc"`
		Data         interface{}  `json:"data"`
		HttpStatus   int          `json:"http_status"`
		Meta         meta         `json:"meta"`
	}

	failed struct {
		ResponseCode string       `json:"response_code"`
		ResponseDesc responseDesc `json:"response_desc"`
		HttpStatus   int          `json:"http_status"`
		Meta         meta         `json:"meta"`
	}

	optionSuccess func(*success)
	optionFailed  func(*failed)
)

func NewResponseDesc(id string, en string) responseDesc {
	return responseDesc{ID: id, EN: en}
}

func NewMeta(version string, status string, apiEnvironment string) meta {
	return meta{Version: version, Status: status, ApiEnvironment: apiEnvironment}
}

func NewSucces(responseCode string, responseDesc responseDesc, data interface{}, httpStatus int, meta meta) *success {
	return &success{
		ResponseCode: responseCode,
		ResponseDesc: responseDesc,
		Data:         data,
		HttpStatus:   httpStatus,
		Meta:         meta,
	}
}

func NewFailed(responseCode string, responseDesc responseDesc, httpStatus int, meta meta) *failed {
	return &failed{
		ResponseCode: responseCode,
		ResponseDesc: responseDesc,
		HttpStatus:   httpStatus,
		Meta:         meta,
	}
}
