package request

type DeleteUserRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
