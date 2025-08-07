package response

type UserStatusResponse struct {
	Sub      string `json:"sub"`
	Provider string `json:"provider"`
	Status   int    `json:"status"`
}
