package cmd

type Access_info struct {
	Client_id  string `json:"client_id"`
	Expires_in int    `json:"expires_in"`
	Scope      string `json:"scope"`
	User_id    string `json:"user_id"`
	Grade      string `json:"grade"`
}
