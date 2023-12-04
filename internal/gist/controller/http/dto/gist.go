package dto

// swagger:model CreateGistRequest
type CreateGistRequest struct {
	GistRequest   GistRequest   `json:"gist"`
	CommitRequest CommitRequest `json:"commit"`
	FilesRequest  []FileRequest `json:"files"`
}

type GistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visible     *bool  `json:"visible"`
}

type CommitRequest struct {
	Comment string `json:"comment"`
}

type FileRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
