package main

type hookBody struct {
	Repository struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		SSHURL   string `json:"ssh_url"`
	} `json:"repository"`
}
