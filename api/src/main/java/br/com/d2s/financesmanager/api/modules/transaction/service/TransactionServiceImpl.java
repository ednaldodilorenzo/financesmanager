package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionDao;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class TransactionServiceImpl implements TransactionService {
    private final TransactionDao transactionDao;

    @Override
    public List<TransactionDto.Response> getTransactions(Integer month, Integer year, Long userId) {
        return transactionDao.findByMonthAndYear(month, year, userId)
                .stream()
                .map(t ->
                        new TransactionDto.Response(t.getDescription(),
                                t.getValue(),
                                t.getPaymentDate(),
                                t.getTransactionDate(),
                                t.getDetail(),
                                new CategoryDto(t.getCategory().getCategoryId().getId(), t.getCategory().getName(), t.getCategory().getType()),
                                new AccountDto(t.getAccount().getId().getId(), t.getAccount().getName(), t.getAccount().getType())))
                .toList();
    }

    @Override
    public List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId) {
        return transactionDao.getSummaryByYear(year, userId);
    }
}
