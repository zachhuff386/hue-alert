package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zachhuff386/hue-alert/account"
	"github.com/zachhuff386/hue-alert/accounts"
	"github.com/zachhuff386/hue-alert/utils"
)

func callbackGet(c *gin.Context) {
	acctType := c.Params.ByName("type")

	auth, authType, err := account.GetAuth(acctType)
	if err != nil {
		c.JSON(400, &errorData{
			Error:   "unknown_type",
			Message: "Unknown account type",
		})
		return
	}

	params := utils.ParseParams(c.Request)
	var x string
	var y string
	var denied bool

	if authType == accounts.Oauth1 {
		x = params.GetByName("oauth_token")
		y = params.GetByName("oauth_verifier")
		denied = params.GetByName("denied") != ""
	} else {
		x = params.GetByName("state")
		y = params.GetByName("code")

		switch params.GetByName("error") {
		case "":
			denied = false
		case "access_denied":
			denied = true
		default:
			c.AbortWithStatus(400)
			return
		}
	}

	if !denied {
		if x == "" || y == "" {
			c.AbortWithStatus(400)
			return
		}

		err = auth.Authorize(x, y)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}

	account.Authenticated <- true

	c.String(200, "Account successfully added")
}
