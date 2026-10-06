package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.*;
import br.com.d2s.financesmanager.api.modules.account.AccountDao;
import br.com.d2s.financesmanager.api.modules.category.CategoryDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionMapper;
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class TransactionServiceImpl implements TransactionService {
    private final TransactionDao transactionDao;
    private final TransactionMapper transactionMapper;
    private final AccountDao accountDao;
    private final CategoryDao categoryDao;

    @Override
    public List<TransactionDto.Response> getTransactions(Integer month, Integer year, Long userId) {
        return transactionDao.findByMonthAndYear(month, year, userId)
                .stream()
                .map(transactionMapper::toDto)
                .toList();
    }

    @Override
    public List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId) {
        return transactionDao.getSummaryByYear(year, userId);
    }

    @Override
    @Transactional
    public void create(TransactionDto.Request dto, Long userId) {
        Transaction transaction = transactionMapper.toEntity(dto);
        transaction.setPaymentDate(dto.transactionDate());
        transaction.setId(new TransactionId(userId, transactionDao.nextId()));
        // Ensures the account and category exist and belong to the user
        Account account = accountDao.findById(new AccountId(userId, dto.accountId())).orElseThrow();
        transaction.setAccount(account);
        Category category = categoryDao.findById(new CategoryId(dto.categoryId(), userId)).orElseThrow();
        transaction.setCategory(category);
        transactionDao.save(transaction);
    }

    @Override
    public List<TransactionDto.Response> getAll(Long userId) {
        return List.of();
    }

    @Override
    public void update(Long id, TransactionDto.Request dto, Long userId) {

    }

    @Override
    public TransactionDto.Response getOne(Long id, Long userId) {
        return null;
    }
}
