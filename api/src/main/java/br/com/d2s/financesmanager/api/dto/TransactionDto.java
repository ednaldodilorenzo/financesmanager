package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.model.CategoryType;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import org.springframework.web.multipart.MultipartFile;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

public class TransactionDto {
    public record Request(
            @NotBlank(message = "Description must be provided.")
            String description,
            @NotNull(message = "Category must my provided")
            PublicId categoryId,
            @NotNull(message = "Value must be provided")
            Long value,
            @NotNull(message = "Account must be provided")
            PublicId accountId,
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
            PublicId id,
            PublicId transactionId,
            String description,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            String detail,
            CategoryDto.Response category,
            AccountDto.Response account) {}

    public record IndividualResponse(
            PublicId id,
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
            PublicId accountId,
            @NotNull(message = "Payment month must be provided")
            Integer paymentMonth,
            @NotNull(message = "Payment year must be provided")
            Integer paymentYear,
            @NotBlank(message = "File type must be provided")
            String fileType) {}

    public record UploadResponse(
            PublicId categoryId,
            PublicId accountId,
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

    public record SplitedRequest(
            PublicId id,
            @NotNull(message = "Value must be provided")
            Long value,
            @NotNull(message = "Account must be provided")
            PublicId accountId,
            @NotNull(message = "Transaction date must be provided")
            Instant transactionDate,
            Instant paymentDate,
            @Size(min = 1, message = "Splits are empty")
            List<TransactionSplit> splits
    ){
        public SplitedRequest {
            if (splits == null) {
                splits = new ArrayList<>();
            }
        }
    }

    public record TransactionSplit(
            // null for a split that does not exist yet
            PublicId id,
            @NotBlank(message = "Description must be provided.")
            String description,
            @NotNull(message = "Category must my provided")
            PublicId categoryId,
            @NotNull(message = "Value must be provided")
            Long value,
            String detail,
            List<TransactionTagDto> tags
    ) {
        public TransactionSplit {
            if (tags == null) {
                tags = new ArrayList<>();
            }
        }
    }

    public record SplitIndividualResponse(
            PublicId id,
            Long value,
            Instant paymentDate,
            Instant transactionDate,
            AccountDto.Response account,
            List<TransactionSplitResponse> splits) {

        public SplitIndividualResponse {
            if (splits == null) {
                splits = new ArrayList<>();
            }
        }
    }

    public record TransactionSplitResponse(
            PublicId id,
            String description,
            CategoryDto.Response category,
            Long value,
            String detail,
            List<TransactionTagDto> tags
    ) {
        public TransactionSplitResponse {
            if (tags == null) {
                tags = new ArrayList<>();
            }
        }
    }
}
