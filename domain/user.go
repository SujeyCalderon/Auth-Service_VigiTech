package domain

type User struct {
    ID       uint   `json:"id"`
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    Active   bool   `json:"active"`
    Role     string `json:"role"`
    SystemID int    `json:"id_Sistema"`
}

type UserRegister struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    Role     string `json:"role"`
    SystemID *int   `json:"id_Sistema"`
}

