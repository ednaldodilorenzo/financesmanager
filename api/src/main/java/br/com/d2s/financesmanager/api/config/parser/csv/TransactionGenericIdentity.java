package br.com.d2s.financesmanager.api.config.parser.csv;

import java.time.Instant;

public record TransactionGenericIdentity(
        Long value,
        Instant paymentDate,
        Instant transactionDate,
        Long userId) {
}
