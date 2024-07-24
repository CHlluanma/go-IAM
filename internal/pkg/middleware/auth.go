package middleware

import "github.com/gin-gonic/gin"

type AuthStrategy interface {
	AuthExecute() gin.HandlerFunc
}

type AuthOperator struct {
	strategy AuthStrategy
}

func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

func (operator *AuthOperator) AuthExecute() gin.HandlerFunc {
	return operator.strategy.AuthExecute()
}
