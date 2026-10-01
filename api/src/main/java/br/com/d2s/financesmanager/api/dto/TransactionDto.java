package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.model.CategoryType;

import java.time.Instant;

public class TransactionDto {
    public record Response(
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            String detail,
            CategoryDto category,
            AccountDto account) {}

    public record TransactionYearlySummary(
            Integer month,
            Integer year,
            Long total,
            CategoryType type) {}
}
