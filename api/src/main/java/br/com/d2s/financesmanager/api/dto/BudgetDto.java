package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import jakarta.validation.constraints.NotNull;

public class BudgetDto {
    public record Request(
            @NotNull(message = "Year must be provided")
            Integer year,
            @NotNull(message = "categoryId must be provided")
            PublicId categoryId,
            @NotNull(message = "Value must be provided")
            Long value){}

    public record Response(PublicId id, Integer year, PublicId categoryId, Long value){}
}
