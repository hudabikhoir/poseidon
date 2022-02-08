package request

import "boilerplate-golang-v2/business/content/spec"

//CreateContentRequest create content request payload
type CreateContentRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

//ToUpsertContentSpec convert into content.UpsertContentSpec object
func (req *CreateContentRequest) ToUpsertContentSpec() *spec.UpsertContentSpec {
	var upsertContentSpec spec.UpsertContentSpec
	upsertContentSpec.Name = req.Name
	upsertContentSpec.Description = req.Description
	upsertContentSpec.Tags = req.Tags

	return &upsertContentSpec
}
