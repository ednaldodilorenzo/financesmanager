package br.com.d2s.financesmanager.api.dto;

import jakarta.validation.constraints.NotBlank;

public class AccountDto {
    public record Response(Long id, String name, String type) {}

    public record Request(
            @NotBlank(message = "Account name must be provided.")
            String name,
            @NotBlank(message = "Account type must be provided.")
            String type) {}
}
