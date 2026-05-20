package repository

import (
	"errors"
	"strings"

	"portfolio/backend/internal/apperr"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

func parseMySQLError(
	err error,
) error {

	if err == nil {
		return nil
	}

	var mysqlErr *mysqlDriver.MySQLError

	if !errors.As(
		err,
		&mysqlErr,
	) {

		return err
	}

	switch mysqlErr.Number {

	// =========================
	// duplicate entry
	// =========================

	case 1062:

		switch {

		case containsConstraint(
			mysqlErr.Message,
			"uq_users_email",
		):

			return apperr.ErrDuplicateEmail

		case containsConstraint(
			mysqlErr.Message,
			"uq_users_auth_user_id",
		):

			return apperr.ErrDuplicateAuthUserID
		}

	// =========================
	// foreign key violation
	// =========================

	case 1452:

		return apperr.ErrForeignKeyViolation
	}

	return err
}

func containsConstraint(
	message string,
	constraints ...string,
) bool {

	for _, c := range constraints {

		if strings.Contains(
			message,
			c,
		) {

			return true
		}
	}

	return false
}
