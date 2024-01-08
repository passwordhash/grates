package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"grates/pkg/utils"
	"net/http/httptest"
)

func main() {
	type UserSignUp struct {
		Email    string `json:"email" binding:"email"`
		Password string `json:"password" binding:"password"`
	}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
			is := utils.IsPassword(fl.Field().String())
			logrus.Info(is, "   is")
			return is
		})
		logrus.Info(err)
	}

	r.POST("/sign-up", func(context *gin.Context) {
		var input UserSignUp

		if err := context.ShouldBindJSON(&input); err != nil {
			logrus.Error(err)
			context.JSON(400, "error")
			return
		}

		logrus.Info(input)

		context.JSON(200, input)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(
		`{"email":"asdf@mail.ru", "password":"asdfasdf"}`,
	))
	r.ServeHTTP(w, req)

}
