package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.model.CategoryType;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.springframework.web.multipart.MultipartFile;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

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
            Instant paymentDate,
            String detail,
            List<TransactionTagDto> tags
    ){
        public Request {
            if (tags == null) {
                tags = new ArrayList<>();
            }
        }
    }

    public record Response(
            Long id,
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            String detail,
            CategoryDto.Response category,
            AccountDto.Response account) {}

    public record IndividualResponse(
            Long id,
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            String detail,
            CategoryDto.Response category,
            AccountDto.Response account,
            List<TransactionTagDto> tags) {

        public IndividualResponse {
            if (tags == null) {
                tags = new ArrayList<>();
            }
        }
    }

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

    public record SearchFilter(
            Instant startDate,
            Instant endDate,
            Integer month,
            Integer year,
            String description,
            List<String> tags
    ){}
}
