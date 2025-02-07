package response

// AuthUserResponse is the structure for the authenticated user's response
type AuthUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Add any other fields you need to send in the response
}
