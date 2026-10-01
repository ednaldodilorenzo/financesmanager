package br.com.d2s.financesmanager.api.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

public record RegisterDto(@NotBlank(message = "")
                          @Email(message = "")
                          String email) {
}
