package br.com.d2s.financesmanager.api.dto;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.NotBlank;

public record SignupDto(
        @NotBlank(message = "Name must be provided")
        String name,
        @NotBlank(message = "Password must be provided")
        String password,
        @NotBlank(message = "Password confirmation must be provided")
        String confirmPassword,
        @NotBlank(message = "Token must be provided")
        String token
) {

    @JsonIgnore
    @AssertTrue(message = "Password and confirmation must match")
    public boolean isPasswordConfirmed() {
        if (password == null || confirmPassword == null) {
            return true;
        }
        return password.equals(confirmPassword);
    }
}
