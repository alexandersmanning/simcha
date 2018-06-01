package database

import "github.com/alexandersmanning/simcha/app/models"

//UserStore is the interface for all User functions that interact with the database
type UserStore interface {
	GetUserByEmailAndPassword(email, password string) (models.User, error)
	UpdatePassword(u models.UserAction, previousPassword, password, confirmationPassword string) error
	UserExists(email string) (bool, error)
	CreateUser(u models.UserAction) error
}

//GetUserByEmailAndPassword checks if the user is in the database, and if it is verifies if the password matches
func (db *DB) GetUserByEmailAndPassword(email, password string) (models.User, error) {
	u := models.User{}
	rows, err := db.Query(
		`SELECT id, email, password_digest FROM users WHERE email = $1`,
		email,
	)

	if err != nil {
		return models.User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&u.Id, &u.Email, &u.PasswordDigest); err != nil {
			return models.User{}, err
		}
	}

	err = u.ComparePassword(password)

	if u.Email == "" || err != nil {
		return models.User{}, &models.ModelError{"Email or Password", "was not found, or does not match our records"}
	}

	return u, nil
}

func (db *DB) UpdatePassword(ua models.UserAction, previousPassword, password, confirmationPassword string) error {
	//Verify password for the new user
	if err := ua.ComparePassword(previousPassword); err != nil {
		return &models.ModelError{"Previous Password", "Does not match current password"}
	}

	ua.SetPassword(password, confirmationPassword)

	digest, err := ua.CreateDigest()
	if err != nil {
		return err
	}
	// save user
	rows, err := db.Query(`
		UPDATE users SET password_digest = $1 WHERE id = $2
	`, digest, ua.User().Id)

	defer rows.Close()

	if err != nil {
		return err
	}

	//remove all existing sessions if successful
	if err := db.RemoveAllUserSessions(ua.User().Id); err != nil {
		return err
	}

	//logout will be done in the controller
	return nil
}

//UserExists checks the existence of an email
func (db *DB) UserExists(email string) (bool, error) {
	var count int
	rows, err := db.Query("SELECT COUNT(*) FROM users WHERE email = $1", email)

	if err != nil {
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}

		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

//CreateUser adds user to system if they do not already exist, and have an appropriate email/password
func (db *DB) CreateUser(ua models.UserAction) error {
	if exists, err := db.UserExists(ua.User().Email); err != nil {
		return err
	} else if exists {
		return &models.ModelError{"Email", "already exists in the system"}
	}

	// set password
	digest, err := ua.CreateDigest()
	if  err != nil {
		return err
	}

	ua.SetDigest(digest)
	ua.SetTimestamps()

	createdAt, modifiedAt := ua.Timestamps()

	rows, err := db.Query(`
		INSERT INTO users (email, password_digest, created_at, modified_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, ua.User().Email, digest, createdAt, modifiedAt)

	if err != nil {
		return err
	}

	defer rows.Close()

	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return err
		}
	}

	ua.SetID(id)

	return nil
}
