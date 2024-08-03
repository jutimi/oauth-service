package middleware

import (
	"net/http"
	"oauth-server/external/client"
	"oauth-server/package/errors"
	"oauth-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/jutimi/grpc-service/workspace"
)

type ownerMiddleware struct{}

func NewOwnerMiddleware() Middleware {
	return &ownerMiddleware{}
}

func (owner *ownerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resErr := errors.New(errors.ErrCodeForbidden)

		payload, err := utils.GetScopeContext[*utils.WorkspacePayload](c, string(utils.WORKSPACE_CONTEXT_KEY))
		if err != nil {
			c.JSON(http.StatusForbidden, utils.FormatErrorResponse(resErr))
			return
		}

		clientGRPC := client.NewWsClient()
		defer clientGRPC.CloseConn()

		userWSId := payload.UserWorkspaceID.String()
		userWS, err := clientGRPC.GetUserWSByFilter(c, &workspace.GetUserWorkspaceByFilterParams{
			Id: &userWSId,
		})
		if err != nil {
			c.JSON(http.StatusForbidden, utils.FormatErrorResponse(resErr))
			return
		}
		if userWS.Data.Role != workspace.UserWorkspaceRole_OWNER {
			c.JSON(http.StatusForbidden, utils.FormatErrorResponse(resErr))
			return
		}

		c.Next()
	}
}
