package controllers

import (
	"fmt"
	"github.com/99-66/go-gin-token-server/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

// CreateToken JWT 토큰을 생성한다
func CreateToken(c *gin.Context) {
	// 토큰 생성을 위해 전달받은 Json 파싱
	var tokenReq models.TokenRequest
	err := c.ShouldBindJSON(&tokenReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := CreateJWTToken(tokenReq)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

// CreateJWTToken JWT 토큰을 생성한다
func CreateJWTToken(tr models.TokenRequest) (*models.TokenInfo, error) {
	// JWT 환경변수 설정
	jwtAcSecret := os.Getenv("JWT_ACCESS_SECRET")
	jwtRefSecret := os.Getenv("JWT_REFRESH_SECRET")

	// Token 정보 생성
	t := &models.TokenInfo{}
	t.TokenRequest = tr
	t.AccessExpires = time.Now().Add(time.Minute * 5).Unix()
	t.RefreshExpires = time.Now().Add(time.Hour * 24 * 5).Unix()
	t.Iat = time.Now().Unix()

	// AccessToken 생성
	// 클레임 설정
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = t.ID
	atClaims["domain"] = t.Domain
	atClaims["roles"] = t.Roles
	atClaims["exp"] = t.AccessExpires
	atClaims["iat"] = t.Iat

	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(jwtAcSecret))
	if err != nil {
		return nil, err
	}

	// RefreshToken 생성
	// 클레임 설정
	rtClaims := jwt.MapClaims{}
	rtClaims["id"] = t.ID
	rtClaims["domain"] = t.Domain
	rtClaims["roles"] = t.Roles
	rtClaims["exp"] = t.RefreshExpires
	rtClaims["iat"] = t.Iat

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(jwtRefSecret))
	if err != nil {
		return nil, err
	}

	return t, nil
}

// VerifyToken JWT 토큰을 검증한다
func VerifyToken(c *gin.Context) {
	token, err := VerifyJWTToken(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token parsed error"})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// VerifyJWTToken JWT 토큰의 유효성 검사를 한다
func VerifyJWTToken(r *http.Request) (*jwt.Token, error){
	tokenString := ExtractToken(r)
	jwtAccSecret := os.Getenv("JWT_ACCESS_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtAccSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TokenValid Middleware에서 토큰의 유효성 검사를 한다
//func TokenValid(r *http.Request) error {
//	token, err := VerifyJWTToken(r)
//	if err != nil {
//		return err
//	}
//
//	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
//		return err
//	}
//
//	return nil
//}

// RefreshToken JWT 토큰을 갱신한다
func RefreshToken(c *gin.Context) {
	// 토큰 갱신을 위해 전달받은 Json 파싱
	var tokenReq models.TokenRequest
	err := c.ShouldBindJSON(&tokenReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refreshToken := tokenReq.RefreshRequestToken.Token

	// JWT 환경변수 설정
	jwtRefSecret := os.Getenv("JWT_REFRESH_SECRET")

	// Refresh Token 을 검증한다
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtRefSecret), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 클레임 파싱
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token parsed error"})
		return
	}

	// AccessToken 으로 갱신하는 경우 실패로 반환한다
	authorized := claims["authorized"]
	if authorized != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed"})
		return
	}

	// Refresh Token 갱신을 위한 파라미터 생성
	var NewTokenRequest models.TokenRequest

	userId := claims["id"].(float64)

	// Roles 목록을 String 슬라이스로 변환한다
	roles := claims["roles"].([]interface{})
	strRoles := make([]string, len(roles))
	for i, v := range roles {
		strRoles[i] = fmt.Sprint(v)
	}

	NewTokenRequest.ID = uint64(userId)
	NewTokenRequest.Domain = claims["domain"].(string)
	NewTokenRequest.Roles = strRoles

	// 블랙리스트로 저장되어 있는 토큰인지 확인한다
	ok = CheckBlackList(NewTokenRequest.Domain, refreshToken)
	if ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has already expired"})
		return
	}

	// 새로운 AccessToken, RefreshToken을 생성한다
	n, err := CreateJWTToken(NewTokenRequest)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 이전 RefreshToken은 Redis에 블랙리스트로 저장한다
	expired := claims["exp"].(float64)
	_, err = CreateBlackList(NewTokenRequest.Domain, refreshToken, int64(expired))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, n)
}

// ExtractToken HTTP Request 헤더에서 토큰을 가져온다
func ExtractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")

	if len(strArr) == 2 {
		if strArr[0] == "Bearer" {
			return strArr[1]
		}
	}

	return ""
}


