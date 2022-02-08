package response

//CreateNewContentResponse Create content response payload
type CreateNewContentResponse struct {
	ID string `json:"id"`
}

//NewCreateNewContentResponse construct CreateNewContentResponse
func NewCreateNewContentResponse(id string) *CreateNewContentResponse {
	return &CreateNewContentResponse{
		id,
	}
}
