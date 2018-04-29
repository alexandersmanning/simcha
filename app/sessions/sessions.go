package sessions

//type SessionStore interface {
//	CurrentUser(env *config.Env, r *http.Request) (*models.User, error)
//	Login(u *models.User, w http.ResponseWriter, r *http.Request) error
//	IsLoggedIn(env *config.Env, r *http.Request) (bool, error)
//	Logout(env *config.Env, w http.ResponseWriter, r *http.Request) error
//}
//
//type Session struct {
//	*sessions.CookieStore
//}

//var currentUser *models.User

//func InitStore(secret string) *Session {
//	cookieStore := sessions.NewCookieStore([]byte(secret))
//	return &Session{ cookieStore }
//}

