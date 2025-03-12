package httpdto

type AddAliasRequest struct {
	Url string `form:"url" validate:"required,url"`
}

type AddAliasResponse struct {
	Error string `json:"error,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type GetUrlResponse struct {
	Error string `json:"error,omitempty"`
	Url   string `json:"url,omitempty"`
}
