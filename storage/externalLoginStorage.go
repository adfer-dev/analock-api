package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/models"
)

const (
	getExternalLoginQuery         = "SELECT * FROM external_login where id = ?;"
	getExternalLoginByClientQuery = "SELECT * FROM external_login where provider_client_id = ?;"
	insertExternalLoginQuery      = "INSERT INTO external_login (provider, provider_client_id, provider_client_token" +
		", user_id) VALUES (?, ?, ?, ?);"
	updateExternalLoginQuery = "UPDATE external_login SET provider = ?, provider_client_id = ?" +
		", user_id = ? WHERE id = ?;"
	deleteExternalLoginQuery = "DELETE FROM external_login WHERE id = ?;"
)

type ExternalLoginStorage struct{}

var externalLoginNotFoundError = &models.DbNotFoundError{DbItem: &models.ExternalLogin{}}
var failedToParseExternalLoginError = &models.DbCouldNotParseItemError{DbItem: &models.ExternalLogin{}}

func (externalLoginStorage *ExternalLoginStorage) Get(id uint) (interface{}, error) {
	result, err := databaseConnection.Query(getExternalLoginQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, externalLoginNotFoundError
	}

	scannedExternalLogin, scanErr := externalLoginStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	externalLogin, ok := scannedExternalLogin.(*models.ExternalLogin)

	if !ok {
		return nil, failedToParseExternalLoginError
	}

	return externalLogin, nil
}

func (externalLoginStorage *ExternalLoginStorage) GetByClientId(clientId string) (interface{}, error) {
	result, err := databaseConnection.Query(getExternalLoginByClientQuery, clientId)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, externalLoginNotFoundError
	}

	scannedExternalLogin, scanErr := externalLoginStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	externalLogin, ok := scannedExternalLogin.(*models.ExternalLogin)

	if !ok {
		return nil, failedToParseExternalLoginError
	}

	return externalLogin, nil
}

func (externalLoginStorage *ExternalLoginStorage) Create(externalLogin interface{}) error {
	dbExternalLogin, ok := externalLogin.(*models.ExternalLogin)

	if !ok {
		return failedToParseUserError
	}

	result, err := databaseConnection.Exec(insertExternalLoginQuery, dbExternalLogin.Provider, dbExternalLogin.ClientId, dbExternalLogin.ClientToken,
		dbExternalLogin.UserRefer)
	if err != nil {
		return err
	}

	externalLoginId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbExternalLogin.Id = uint(externalLoginId)

	return nil
}

func (externalLoginStorage *ExternalLoginStorage) Update(externalLogin interface{}) error {
	dbExternalLogin, ok := externalLogin.(*models.ExternalLogin)

	if !ok {
		return failedToParseUserError
	}

	result, err := databaseConnection.Exec(updateExternalLoginQuery, dbExternalLogin.Provider, dbExternalLogin.ClientId,
		dbExternalLogin.UserRefer)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return externalLoginNotFoundError
	}

	return nil
}

func (externalLoginStorage *ExternalLoginStorage) Delete(id uint) error {
	result, err := databaseConnection.Exec(deleteExternalLoginQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return externalLoginNotFoundError
	}

	return nil
}

func (externalLoginStorage *ExternalLoginStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var externalLogin models.ExternalLogin

	scanErr := rows.Scan(&externalLogin.Id, &externalLogin.Provider, &externalLogin.ClientId, &externalLogin.ClientToken,
		&externalLogin.UserRefer)

	return &externalLogin, scanErr
}
