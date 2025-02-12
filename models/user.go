package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Email  string
	Name   string
	UserID int64
}

type UserData struct {
	OAuthID string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	OAuth   string `json:"oauth,omitempty"`
	IpAddr  string `json:"_,omitempty"`
}

type UserSession struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	LoginToken string    `db:"login_token"`
	Superseded bool      `db:"superseded"`
	CreatedAt  time.Time `db:"created_at"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

// to use main db that initialised in main.go
type UserModel struct {
	db    *sql.DB
	redis *redis.Client
}

func (u *UserModel) UserInit(db *sql.DB, redis *redis.Client) UserModel {
	return UserModel{
		db:    db,
		redis: redis,
	}
}
func (m *UserModel) AuthUser(userdata UserData) (string, error) {
	var user_id int64
	err := m.db.QueryRow("SELECT id  FROM users WHERE email = ?", userdata.Email).Scan(&user_id)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		// User is not registered, save the user information to the database
		stmt := "INSERT INTO users (email, name, oauth, oauth_id, picture) VALUES (?, ?, ?,?,?)"
		result, err := m.db.Exec(stmt, userdata.Email, userdata.Name, userdata.OAuth, userdata.OAuthID, userdata.Picture)
		if err != nil {
			return "", err
		}

		user_id, err = result.LastInsertId()
		if err != nil {
			return "", err
		}
	}

	login_token, err := m.generateToken(userdata.Name, user_id)
	if err != nil {
		return "", err
	}

	err = m.SetLoginToken(login_token, userdata.IpAddr, user_id)
	return login_token, err
}

func (m *UserModel) ValidToken(login_token string) (int, error) {
	var user_id int
	err := m.db.QueryRow("SELECT user_id FROM users_session WHERE login_token = ? AND superseded = 0", login_token).Scan(&user_id)
	return user_id, err
}

func (m *UserModel) SetLoginToken(token, ip_addr string, user_id int64) error {
	err := m.Logout(user_id)
	stmt := fmt.Sprintf("INSERT INTO users_session (login_token, user_id, ip_addr) VALUES ('%s', '%d', '%s' )", token, user_id, ip_addr)
	if err == nil {
		_, err = m.db.Exec(stmt)
	}

	return err
}

// logout
func (m *UserModel) Logout(user_id int64) error {
	_, err := m.db.Exec("UPDATE `users_session` SET `superseded` = 1 WHERE  `id` = ?", user_id)
	return err
}

func (m *UserModel) ActivityLog(activity string, uid int64) {
	_, _ = m.db.Exec("UPDATE `user_log` SET superseded = 1 WHERE activity = ? AND uid = ?", activity, uid)
	_, _ = m.db.Exec("INSERT INTO `user_log` SET  activity = ? , uid = ?, superseded = 0", activity, uid)
}

func (m *UserModel) generateToken(username string, user_id int64) (string, error) {
	// Create the claims

	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	mySigningKeyStr := os.Getenv("SIGNING_KEY")
	mySigningKey := []byte(mySigningKeyStr)
	claims := MyCustomClaims{
		username,
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(), // Token expires in 4 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(mySigningKey)
}
