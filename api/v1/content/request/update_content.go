package request

import "boilerplate-golang-v2/business/content/spec"

//UpdateContentRequest update content request payload
type UpdateContentRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Version     int      `json:"version" validate:"required"`
}

//ToUpsertContentSpec convert into content.UpsertContentSpec object
func (req *UpdateContentRequest) ToUpsertContentSpec() *spec.UpsertContentSpec {
	var upsertContentSpec spec.UpsertContentSpec
	upsertContentSpec.Name = req.Name
	upsertContentSpec.Description = req.Description
	upsertContentSpec.Tags = req.Tags

	return &upsertContentSpec
}
