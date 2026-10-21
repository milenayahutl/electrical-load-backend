
package currentuser

import "github.com/gin-gonic/gin"

const UserIDKey = "authenticated_user_id"

func GetCurrentUserID(ctx *gin.Context) uint {
	value, exists := ctx.Get(UserIDKey)

	if !exists {
		return 0
	}

	userID, ok := value.(uint)

	if !ok {
		return 0
	}

	return userID
}
