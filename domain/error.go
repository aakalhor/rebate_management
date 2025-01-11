package domain

import "errors"

var (
	ErrInternalServerError     = errors.New("internal server error")
	ErrInvalidInterval         = errors.New("invalid interval")
	ErrNotEligible             = errors.New("criteria is not eligible")
	ErrClaimNotFound           = errors.New("claim does not exist")
	ErrRebateAlreadyClaimed    = errors.New("rebate already claimed")
	ErrMultipleProgramName     = errors.New("a program with the same name already exists")
	ErrTransactionCanNotCreate = errors.New("transaction can not be created")
	ErrTransactionNotFound     = errors.New("transaction has not founded")
	ErrRebateNotFound          = errors.New("rebate has not founded")
	ErrFailedToListClaims      = errors.New("failed to list claims within interval")
	ErrFailedToGetCache        = errors.New("failed to get cache date")
	ErrFailedToStoreCache      = errors.New("failed to store cache date")
)
