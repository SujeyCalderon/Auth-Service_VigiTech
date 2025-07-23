package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"auth-service/application"
	"auth-service/domain"
)


func RegisterUserHandler(svc *application.Service) gin.HandlerFunc {
  return func(c *gin.Context) {
    var req domain.UserRegister
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos"})
      return
    }

    user, err := svc.RegisterUser(req.Name, req.Email, req.Password, req.Role, req.SystemID)
    if err != nil {
      switch err {
      case domain.ErrEmailAlreadyExists:
        c.JSON(http.StatusConflict, gin.H{"success": false, "message": "El email ya está registrado"})
      case domain.ErrUnauthorized:
        c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "ID de sistema inválido para este rol"})
      default:
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "No se pudo registrar el usuario"})
      }
      return
    }

    token, err := svc.TokenService.GenerateToken(strconv.Itoa(int(user.ID)), 24*time.Hour, user.Role)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error generando token"})
      return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
  }
}


func LoginUserHandler(svc *application.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos"})
            return
        }

        token, user, err := svc.Login(req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "success":  true,
            "token":    token,
            "userName": user.Name,
            "role":     user.Role,
            "isActive": user.Active,
        })
    }
}


func UpdateUserStatusHandler(svc *application.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDCtx, _ := c.Get("userID")
		roleCtx, _ := c.Get("role")
		role, _ := roleCtx.(string)
		idParam := c.Param("id")


		if role != "admin" && idParam != userIDCtx {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "No tienes permisos"})
			return
		}

		idUint, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID inválido"})
			return
		}
		var body struct{ Active bool `json:"active"` }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos inválidos"})
			return
		}
		if err := svc.UpdateUserStatus(uint(idUint), body.Active); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error actualizando estado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Estado actualizado"})
	}
}

func GetAllUsersHandler(svc *application.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIfc, _ := c.Get("role")
		role, _ := roleIfc.(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "No tienes permiso para esto"})
			return
		}

		users, err := svc.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error obteniendo usuarios"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "users": users})
	}
}