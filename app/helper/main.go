package helper

import (
	mysql_repository "gin-boilerplate/app/repository/mysql"
)

type HelperCollections struct {
	OauthHelper OauthHelper
}

func RegisterHelpers(
	mysqlRepo mysql_repository.MysqlRepositoryCollections,
) HelperCollections {
	oauthHelper := NewOauthHelper()

	return HelperCollections{
		OauthHelper: oauthHelper,
	}
}
