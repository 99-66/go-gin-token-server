package controllers

import (
	"github.com/99-66/go-gin-token-server/config"
	"time"
)

// CreateBlackList RefreshToken을 통해 재발급한 경우 이전 refreshToken을 블랙리스트로 넣는다(재사용 금지)
// TTL을 refresh Token 만료일자와 동일하게 설정한다
func CreateBlackList(domain, token string, expiredTime int64) (string, error) {
	// 현재시간에서 token Expired time 차이를 구해서 expired 시간으로 설정한다
	rt := time.Unix(expiredTime, 0)
	now := time.Now()

	result, err := config.REDIS.Set(config.CTX, token, domain, rt.Sub(now)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// CheckBlackList RefreshToken이 블랙리스트에 있는지 확인한다
func CheckBlackList(domain, token string) bool {
	d, err := config.REDIS.Get(config.CTX, token).Result()
	if err != nil {
		return false
	}

	if d == domain {
		return true
	}

	return false
}