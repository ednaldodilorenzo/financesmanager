package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.model.CategoryType;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.springframework.web.multipart.MultipartFile;

import java.time.Instant;

public class TransactionDto {
    public record Request(
            @NotBlank(message = "Description must be provided.")
            String description,
            @NotNull(message = "Category must my provided")
            Long categoryId,
            @NotNull(message = "Value must be provided")
            Long value,
            @NotNull(message = "Account must be provided")
            Long accountId,
            @NotNull(message = "Transaction date must be provided")
            Instant transactionDate,
            String detail
    ){}

    public record Response(
            Long id,
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            String detail,
            CategoryDto.Response category,
            AccountDto.Response account) {}

    public record UploadRequest(
            @NotNull(message = "File must be provided")
            MultipartFile file,
            @NotNull(message = "Account must be provided")
            Long accountId,
            @NotNull(message = "Payment month must be provided")
            Integer paymentMonth,
            @NotNull(message = "Payment year must be provided")
            Integer paymentYear,
            @NotBlank(message = "File type must be provided")
            String fileType) {}

    public record UploadResponse(
            Long categoryId,
            Long accountId,
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            Boolean duplicated,
            String detail) {}

    public record TransactionYearlySummary(
            Integer month,
            Integer year,
            Long total,
            CategoryType type) {}
}
