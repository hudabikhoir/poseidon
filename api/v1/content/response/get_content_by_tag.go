package response

import "boilerplate-golang-v2/business/content"

//GetContentByTagResponse Get content by tag response payload
type GetContentByTagResponse struct {
	Contents []*GetContentByIDResponse `json:"contents"`
}

//NewGetContentByTagResponse construct GetContentByTagResponse
func NewGetContentByTagResponse(contents []content.Content) *GetContentByTagResponse {
	var contentResponses []*GetContentByIDResponse
	contentResponses = make([]*GetContentByIDResponse, 0)

	for _, content := range contents {
		contentResponses = append(contentResponses, NewGetContentByIDResponse(content))
	}

	return &GetContentByTagResponse{
		contentResponses,
	}
}
