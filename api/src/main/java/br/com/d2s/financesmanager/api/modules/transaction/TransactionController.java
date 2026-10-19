package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.modules.transaction.service.TransactionService;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/transactions")
@AllArgsConstructor
public class TransactionController {

    private TransactionService transactionService;

    @GetMapping("/")
    public List<TransactionDto.Response> getTransactions(@ModelAttribute TransactionDto.SearchFilter filter, @AuthenticationPrincipal Jwt jwt) {
        return transactionService.getTransactions(filter, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public void post(@RequestBody @Valid TransactionDto.Request dto, @AuthenticationPrincipal Jwt jwt) {
        transactionService.create(dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/{id}")
    public TransactionDto.IndividualResponse get(@PathVariable("id") PublicId id, @AuthenticationPrincipal Jwt jwt) {
        return transactionService.getTransaction(id.value(), Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/summary/{year}")
    public List<TransactionDto.TransactionYearlySummary> getSummary(@PathVariable("year") Integer year, @AuthenticationPrincipal Jwt jwt) {
        return transactionService.getSummaryByYear(year, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping(value = "/upload", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public List<TransactionDto.UploadResponse> uploadImportFile(@ModelAttribute TransactionDto.UploadRequest dto,
                                    @AuthenticationPrincipal Jwt jwt) {
        return transactionService.parseTransactionImport(dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/list")
    public void postList(@RequestBody List<TransactionDto.Request> dtos, @AuthenticationPrincipal Jwt jwt) {
        transactionService.createAll(dtos, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PatchMapping("/{id}")
    public void patch(@PathVariable PublicId id, @RequestBody @Valid TransactionDto.Request dto, @AuthenticationPrincipal Jwt jwt) {
        transactionService.update(id.value(), dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}
