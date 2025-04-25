package mapper

import (
	"github.com/ClothShop/auth-service/internal/dtos"
	"github.com/ClothShop/auth-service/internal/models"
	"github.com/jinzhu/copier"
)

func ToMeResponse(user *models.User) *dtos.MeResponse {
	var res dtos.MeResponse
	_ = copier.Copy(&res, &user)
	res.Id = user.ID.String()
	res.Role = string(user.Role)
	return &res
}
