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
		workspacePayload, _ := c.Get(string(utils.WORKSPACE_CONTEXT_KEY))
		payload, _ := workspacePayload.(*utils.WorkspacePayload)

		clientGRPC := client.NewWorkspaceClient()
		defer clientGRPC.CloseConn()

		userWorkspaceId := payload.UserWorkspaceId.String()
		userWorkspace, err := clientGRPC.GetUserWorkspaceByFilter(c, &workspace.GetUserWorkspaceByFilterParams{
			Id: &userWorkspaceId,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.FormatErrorResponse(resErr))
			return
		}
		if userWorkspace.Data.Role != workspace.UserWorkspaceRole_OWNER {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.FormatErrorResponse(resErr))
			return
		}

		c.Next()
	}
}
