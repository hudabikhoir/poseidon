package spec

//UpsertContentSpec create and update content spec
type UpsertContentSpec struct {
	Name        string   `validate:"required"`
	Description string   `validate:"required,min=3"`
	Tags        []string `validate:"required"`
}
