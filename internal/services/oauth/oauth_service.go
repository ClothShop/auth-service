package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/ClothShop/auth-service/internal/models"
	"log"
	"net/http"
	"net/url"
	"os"
)

func GetGoogleAuthURL(state string) (string, error) {
	u, err := url.Parse("https://accounts.google.com/o/oauth2/auth")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URI"))
	q.Set("response_type", "code")
	q.Set("scope", fmt.Sprintf("%s", ScopesToString()))
	q.Set("state", state)
	q.Set("access_type", "offline")
	q.Set("prompt", "consent")

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func ScopesToString() string {
	return "openid email profile"
}

func ExchangeCodeForUser(code string) (*models.User, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	data := getDataForCode(code)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Println("Error sending request POST_FORM", err)
		return nil, err
	}
	defer resp.Body.Close()

	var tokenRes struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		log.Println("Error decoding response body1", err)
		return nil, err
	}

	userInfoReq, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	userInfoReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)
	client := &http.Client{}
	userResp, err := client.Do(userInfoReq)
	if err != nil {
		log.Println("Error sending request GET_FORM", err)
		return nil, err
	}
	defer userResp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		log.Println("Error decoding response body2", err)
		return nil, err
	}

	return &models.User{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

func getDataForCode(code string) url.Values {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URI"))
	data.Set("grant_type", "authorization_code")
	return data
}
