package container

import (
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/adapter/repository/tokenrepo"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/adapter/webapi/user"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/middleware"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/authhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/swaggerhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/userhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/service/authservice"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/service/userservice"
	pb "github.com/bhankey/pharmacy-automatization-user/pkg/api/userservice"
)

func (c *Container) GetV1AuthHandler() *authhandler.AuthHandler {
	const key = "V1AuthHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*authhandler.AuthHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := authhandler.NewAuthHandler(c.getBaseHandler(), c.getAuthSrv())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) GetV1SwaggerHandler() *swaggerhandler.SwaggerHandler {
	const key = "V1SwaggerHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*swaggerhandler.SwaggerHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := swaggerhandler.NewSwaggerHandler(c.getBaseHandler())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) GetV1UserHandler() *userhandler.UserHandler {
	const key = "V1UserHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*userhandler.UserHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := userhandler.NewUserHandler(c.getBaseHandler(), c.getUserSrv(), c.GetAuthMiddleware())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getBaseHandler() *http.BaseHandler {
	const key = "BaseHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*http.BaseHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := http.NewHandler(c.logger)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getUserSrv() *userservice.UserService {
	const key = "UserSrv"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*userservice.UserService)
		if ok {
			return typedDependency
		}
	}

	typedDependency := userservice.NewUserService(
		c.getUserAdapter(),
	)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getAuthSrv() *authservice.AuthService {
	const key = "AuthSrv"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*authservice.AuthService)
		if ok {
			return typedDependency
		}
	}

	typedDependency := authservice.NewAuthService(
		c.getUserAdapter(),
		c.getTokenStorage(),
		c.jwtKey,
	)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getTokenStorage() *tokenrepo.TokenRepo {
	const key = "TokenStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*tokenrepo.TokenRepo)
		if ok {
			return typedDependency
		}
	}

	typedDependency := tokenrepo.NewTokenRepo(c.masterPostgresDB, c.slavePostgresDB)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getUserAdapter() *user.APIClient {
	const key = "UserAdapter"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*user.APIClient)
		if ok {
			return typedDependency
		}
	}

	typedDependency := user.NewUserAPIClient(c.getUserServiceAPIClient())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getUserServiceAPIClient() pb.UserServiceClient {
	const key = "UserServiceAPIClient"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(pb.UserServiceClient)
		if ok {
			return typedDependency
		}
	}

	typedDependency := pb.NewUserServiceClient(c.userServiceConn)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) GetAuthMiddleware() *middleware.AuthMiddleware {
	const key = "AuthMiddleware"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*middleware.AuthMiddleware)
		if ok {
			return typedDependency
		}
	}

	typedDependency := middleware.NewAuthMiddleware(c.logger, c.jwtKey)

	c.dependencies[key] = typedDependency

	return typedDependency
}
