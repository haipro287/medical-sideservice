package error

import "golang.org/x/xerrors"

var (
	ErrMissingUserId = xerrors.New("ERROR.API.AUTH.MISSING_USER_ID")
	ErrNoPermission  = xerrors.New("ERROR.API.NO_PERMISSION")
)
