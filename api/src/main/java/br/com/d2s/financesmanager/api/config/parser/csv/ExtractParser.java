package br.com.d2s.financesmanager.api.config.parser.csv;

import br.com.d2s.financesmanager.api.dto.TransactionDto;

import java.io.IOException;
import java.io.InputStream;
import java.time.Instant;
import java.util.List;
import java.util.Set;
import java.util.function.Function;
import java.util.function.Predicate;

public interface ExtractParser {
    boolean supports(String bank);

    List<TransactionDto.UploadResponse> parse(
            InputStream inputStream,
            String bank,
            Long accountId,
            Long userId,
            Function<List<TransactionGenericIdentity>, Set<TransactionGenericIdentity>> findDuplicates,
            Instant fixedPaymentDate) throws IOException;
}
