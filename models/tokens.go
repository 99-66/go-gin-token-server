package models

import "encoding/json"

// RefreshRequestToken TokenRequest 임베딩 전용
// TokenInfo에 임베딩되었을 때는 사용하지 않는다
type RefreshRequestToken struct {
	Token string `json:"refresh_token,omitempty"`
}

// TokenRequest 토큰 생성을 위해 서버로 전달한 정보의 Structure
type TokenRequest struct {
	ID	uint64	`json:"id"`
	Domain 	string `json:"domain"`
	Roles []string `json:"roles"`
	RefreshRequestToken
}

// TokenInfo 발급한 토큰의 정보를 위한 Structure
type TokenInfo struct {
	Authorized bool `json:"authorized"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessExpires int64 `json:"access_expire"`
	RefreshExpires int64 `json:"refresh_expire"`
	Iat int64 `json:"iat"`
	TokenRequest
}

func (t TokenInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}