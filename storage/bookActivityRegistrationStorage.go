package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/models"
)

const (
	getBookActivityRegistrationByIdentifierQuery = "SELECT arb.id, arb.internet_archive_id, ar.registration_date, ar.user_id FROM activity_registration_book arb INNER JOIN activity_registration ar ON (arb.registration_id = ar.id) WHERE arb.id = ?;"
	getUserBookActivityRegistrationsQuery        = "SELECT arb.id, arb.internet_archive_id, ar.registration_date, ar.user_id FROM activity_registration_book arb INNER JOIN activity_registration ar ON (arb.registration_id = ar.id) WHERE ar.user_id = ?;"
	insertBookActivityRegistrationQuery          = "INSERT INTO activity_registration_book (internet_archive_id, registration_id) VALUES (?, ?);"
	updateBookActivityRegistrationQuery          = "UPDATE activity_registration_book SET internet_archive_id = ? WHERE id = ?;"
	deleteBookActivityRegistrationQuery          = "DELETE FROM activity_registration_book WHERE id = ?;"
)

type BookActivityRegistrationStorage struct{}

var bookActivityRegistrationNotFoundError = &models.DbNotFoundError{DbItem: &models.BookActivityRegistration{}}
var failedToParseBookActivityRegistrationError = &models.DbCouldNotParseItemError{DbItem: &models.BookActivityRegistration{}}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) Get(id uint) (interface{}, error) {
	result, err := databaseConnection.Query(getBookActivityRegistrationByIdentifierQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, bookActivityRegistrationNotFoundError
	}

	scannedBookActivityRegistration, scanErr := bookActivityRegistrationStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	bookActivityRegistration, ok := scannedBookActivityRegistration.(models.BookActivityRegistration)

	if !ok {
		return nil, failedToParseBookActivityRegistrationError
	}

	return &bookActivityRegistration, nil
}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) GetByUserId(userId uint) (interface{}, error) {
	var userBookActivityRegistrations []*models.BookActivityRegistration
	result, err := databaseConnection.Query(getUserBookActivityRegistrationsQuery, userId)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		scannedBookActivityRegistration, scanErr := bookActivityRegistrationStorage.Scan(result)

		if scanErr != nil {
			return nil, scanErr
		}
		bookActivityRegistration, ok := scannedBookActivityRegistration.(models.BookActivityRegistration)

		if !ok {
			return nil, failedToParseBookActivityRegistrationError
		}

		userBookActivityRegistrations = append(userBookActivityRegistrations, &bookActivityRegistration)
	}

	return userBookActivityRegistrations, nil
}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) Create(bookRegistration interface{}) error {
	dbBookRegistration, ok := bookRegistration.(*models.BookActivityRegistration)

	if !ok {
		return failedToParseDiaryEntryError
	}

	result, err := databaseConnection.Exec(insertBookActivityRegistrationQuery,
		dbBookRegistration.InternetArchiveIdentifier,
		dbBookRegistration.Registration.Id)

	if err != nil {
		return err
	}

	bookRegistrationId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbBookRegistration.Id = uint(bookRegistrationId)

	return nil
}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) Update(bookRegistration interface{}) error {
	dbBookRegistration, ok := bookRegistration.(*models.BookActivityRegistration)

	if !ok {
		return failedToParseBookActivityRegistrationError
	}

	result, err := databaseConnection.Exec(updateDiaryEntryQuery,
		dbBookRegistration.InternetArchiveIdentifier,
		dbBookRegistration.Id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return bookActivityRegistrationNotFoundError
	}

	return nil
}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) Delete(id uint) error {
	result, err := databaseConnection.Exec(deleteBookActivityRegistrationQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return bookActivityRegistrationNotFoundError
	}

	return nil
}

func (bookActivityRegistrationStorage *BookActivityRegistrationStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var bookActivityRegistration models.BookActivityRegistration

	scanErr := rows.Scan(&bookActivityRegistration.Id, &bookActivityRegistration.InternetArchiveIdentifier,
		&bookActivityRegistration.Registration.RegistrationDate, &bookActivityRegistration.Registration.UserRefer)

	return bookActivityRegistration, scanErr
}
