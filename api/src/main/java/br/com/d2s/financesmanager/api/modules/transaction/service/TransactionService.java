package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.dto.TransactionDto;


import java.util.List;

public interface TransactionService {
    List<TransactionDto.Response> getTransactions(Integer month, Integer year, Long userId);

    List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId);

    void create(TransactionDto.Request dto, Long userId);

    List<TransactionDto.Response> getAll(Long userId);

    void update(Long id, TransactionDto.Request dto, Long userId);

    TransactionDto.Response getOne(Long id, Long userId);
}
