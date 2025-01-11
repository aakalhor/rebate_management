package http

import (
	"awesomeProject2/rebate/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CustomError struct {
	Status       int
	CodeResponse interface{} // Detailed error response (could be a string, JSON object, etc.)
}

type CodeResponse struct {
	Code    string
	Message string
}

// ErrHandler processes errors based on the provided ErrorMap
func ErrHandler(context *gin.Context, err error, errMap ErrorMap) {
	// Look up the error in the map
	if customErr, exists := errMap[err]; exists {
		context.JSON(customErr.Status, gin.H{
			"error": customErr.CodeResponse,
		})
		return
	}

	// Default response for unhandled errors
	context.JSON(http.StatusInternalServerError, gin.H{
		"error":  "An unexpected error occurred",
		"detail": err.Error(),
	})
}

type ErrorMap map[error]CustomError

var errMap = ErrorMap{
	domain.ErrInternalServerError: CustomError{
		Status:       http.StatusInternalServerError,
		CodeResponse: ErrInternalServerError,
	},
	domain.ErrInvalidInterval: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrInvalidInterval,
	},
	domain.ErrNotEligible: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrNotEligible,
	},
	domain.ErrClaimNotFound: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrClaimNotFound,
	},
	domain.ErrRebateAlreadyClaimed: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrRebateAlreadyClaimed,
	},
	domain.ErrMultipleProgramName: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrMultipleProgramName,
	},
	domain.ErrTransactionCanNotCreate: CustomError{
		Status:       http.StatusInternalServerError,
		CodeResponse: ErrTransactionCanNotCreate,
	},
	domain.ErrTransactionNotFound: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrTransactionNotFound,
	},
	domain.ErrRebateNotFound: CustomError{
		Status:       http.StatusBadRequest,
		CodeResponse: ErrRebateNotFound,
	},
	domain.ErrFailedToListClaims: CustomError{
		Status:       http.StatusInternalServerError,
		CodeResponse: ErrFailedToListClaims,
	},
	domain.ErrFailedToGetCache: CustomError{
		Status:       http.StatusInternalServerError,
		CodeResponse: ErrFailedToGetCache,
	},
	domain.ErrFailedToStoreCache: CustomError{
		Status:       http.StatusInternalServerError,
		CodeResponse: ErrFailedToStoreCache,
	},
}

var (
	ErrInternalServerError = CodeResponse{
		Code:    "err-internal-server-error",
		Message: "An unexpected error occurred. Please try again later.",
	}
	ErrInvalidInterval = CodeResponse{
		Code:    "err-invalid-interval",
		Message: "The provided date interval is invalid.",
	}
	ErrNotEligible = CodeResponse{
		Code:    "err-not-eligible",
		Message: "The user is not eligible for this action.",
	}
	ErrClaimNotFound = CodeResponse{
		Code:    "err-claim-not-found",
		Message: "No claim found with the specified ID.",
	}
	ErrRebateAlreadyClaimed = CodeResponse{
		Code:    "err-rebate-already-claimed",
		Message: "The rebate has already been claimed.",
	}
	ErrMultipleProgramName = CodeResponse{
		Code:    "err-multiple-program-name",
		Message: "A rebate program with the same name already exists.",
	}
	ErrTransactionCanNotCreate = CodeResponse{
		Code:    "err-transaction-cannot-create",
		Message: "The transaction could not be created. Please check the input data.",
	}
	ErrTransactionNotFound = CodeResponse{
		Code:    "err-transaction-not-found",
		Message: "No transaction found with the specified ID.",
	}
	ErrRebateNotFound = CodeResponse{
		Code:    "err-rebate-not-found",
		Message: "No rebate program found with the specified ID.",
	}
	ErrFailedToListClaims = CodeResponse{
		Code:    "err-failed-to-list-claims",
		Message: "Failed to retrieve the list of claims. Please try again later.",
	}
	ErrFailedToGetCache = CodeResponse{
		Code:    "err-failed-to-get-cache",
		Message: "Failed to fetch cache date.",
	}
	ErrFailedToStoreCache = CodeResponse{
		Code:    "err-failed-to-store-cache",
		Message: "Failed to store cache date.",
	}
)
