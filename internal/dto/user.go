package dto

// RegisterUserRequest merepresentasikan payload JSON untuk registrasi pengguna baru.
type RegisterUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest merepresentasikan payload JSON untuk login pengguna.
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// UserResponse merepresentasikan data pengguna yang dikirim dalam response.
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// LoginResponse mengembalikan data pengguna dan token JWT setelah login berhasil.
type LoginResponse struct {
    User  UserResponse `json:"user"`
    Token string       `json:"token"`
}

// menggunakan tag validate agar library validator/v10 dapat memeriksa input secara otomatis.