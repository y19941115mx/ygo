package user

import "github.com/y19941115mx/ygo/app/provider/user"

// ConvertUserToDTO 将user转换为UserDTO
func ConvertUserToDTO(user *user.User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		ID:        int64(user.ID),
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt,
	}
}
