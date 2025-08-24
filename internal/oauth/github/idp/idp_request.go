package idp

type signInPayload struct {
	RequestURI        string `json:"requestUri"`
	PostBody          string `json:"postBody"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}
