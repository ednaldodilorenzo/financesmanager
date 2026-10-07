package br.com.d2s.financesmanager.api.dto;

import jakarta.validation.constraints.NotNull;

public class BudgetDto {
    public record Request(
            @NotNull(message = "Year must be provided")
            Integer year,
            @NotNull(message = "categoryId must be provided")
            Long categoryId,
            @NotNull(message = "Value must be provided")
            Long value){}

    public record Response(Long id, Integer year, Long categoryId, Long value){}
}
