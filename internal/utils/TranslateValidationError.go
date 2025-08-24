package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorsToMap mengubah ValidationErrors menjadi map[field]pesan.
func ValidationErrorsToMap(errs validator.ValidationErrors) map[string]string {
    errorsMap := make(map[string]string)
    for _, e := range errs {
        switch e.Field() {
        case "Name":
            switch e.Tag() {
            case "required":
                errorsMap["name"] = "Nama wajib diisi"
            case "min":
                errorsMap["name"] = fmt.Sprintf("Nama minimal %s karakter", e.Param())
            }
        case "Email":
            switch e.Tag() {
            case "required":
                errorsMap["email"] = "Email wajib diisi"
            case "email":
                errorsMap["email"] = "Format email tidak valid"
            }
        case "Password":
            switch e.Tag() {
            case "required":
                errorsMap["password"] = "Kata sandi wajib diisi"
            case "min":
                errorsMap["password"] = fmt.Sprintf("Kata sandi minimal %s karakter", e.Param())
            }
        // Tambahkan field lain sesuai kebutuhan
        default:
            // Nama field diubah menjadi huruf kecil sebagai key
            errorsMap[e.Field()] = e.Error()
        }
    }
    return errorsMap
}
