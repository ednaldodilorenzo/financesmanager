package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import jakarta.validation.constraints.NotBlank;

public class AccountDto {
    public record Response(PublicId id, String name, String type, Integer dueDay) {}

    public record Request(
            @NotBlank(message = "Account name must be provided.")
            String name,
            @NotBlank(message = "Account type must be provided.")
            String type,
            Integer dueDay) {}
}
