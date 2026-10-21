package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.dto.TransactionDto;


import java.util.List;

public interface TransactionService {
    List<TransactionDto.Response> getTransactions(TransactionDto.SearchFilter filter, Long userId);

    TransactionDto.IndividualResponse getTransaction(Long id, Long userId);

    List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId);

    void create(TransactionDto.Request dto, Long userId);

    void update(Long id, TransactionDto.Request dto, Long userId);

    List<TransactionDto.UploadResponse> parseTransactionImport(TransactionDto.UploadRequest dto, Long userId);

    void createAll(List<TransactionDto.Request> dtos, Long userId);

    void delete(Long id, Long userId);
}
