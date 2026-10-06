package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.modules.transaction.service.TransactionService;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
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
    public List<TransactionDto.Response> getTransactions(@RequestParam Integer month, @RequestParam Integer year, @AuthenticationPrincipal Jwt jwt) {
        return transactionService.getTransactions(month, year, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/")
    public void post(@RequestBody @Valid TransactionDto.Request dto, @AuthenticationPrincipal Jwt jwt) {
        transactionService.create(dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/summary/{year}")
    public List<TransactionDto.TransactionYearlySummary> getSummary(@PathVariable("year") Integer year, @AuthenticationPrincipal Jwt jwt) {
        return transactionService.getSummaryByYear(year, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}
