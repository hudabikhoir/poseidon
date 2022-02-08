package response

import (
	"boilerplate-golang-v2/business/content"
	"time"
)

//GetContentByIDResponse Get content by ID response payload
type GetContentByIDResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	Version     int       `json:"version"`
}

//NewGetContentByIDResponse construct GetContentByIDResponse
func NewGetContentByIDResponse(content content.Content) *GetContentByIDResponse {
	var contentResponse GetContentByIDResponse
	contentResponse.ID = content.ID
	contentResponse.Name = content.Name
	contentResponse.Description = content.Description
	contentResponse.Tags = content.Tags
	contentResponse.ModifiedAt = content.ModifiedAt
	contentResponse.Version = content.Version

	return &contentResponse
}
