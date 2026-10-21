package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.config.parser.csv.ExtractParser;
import br.com.d2s.financesmanager.api.config.parser.csv.TransactionGenericIdentity;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.*;
import br.com.d2s.financesmanager.api.modules.account.AccountDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionMapper;
import br.com.d2s.financesmanager.api.modules.transactionsplit.TransactionSplitDao;
import br.com.d2s.financesmanager.api.modules.transactionsplit.TransactionSplitMapper;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.io.*;
import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;


@Service
@AllArgsConstructor
public class TransactionServiceImpl implements TransactionService {
    private final TransactionDao transactionDao;
    private final TransactionSplitDao transactionSplitDao;
    private final TransactionMapper transactionMapper;
    private final TransactionSplitMapper transactionSplitMapper;
    private final AccountDao accountDao;
    private final ExtractParser extractParser;

    @Override
    public List<TransactionDto.Response> getTransactions(TransactionDto.SearchFilter filter, Long userId) {
        return transactionSplitDao
                .findAll(TransactionSpecs.from(filter, userId))
                .stream()
                .map(transactionSplitMapper::toDto).toList();
    }

    @Override
    public TransactionDto.IndividualResponse getTransaction(Long id, Long userId) {
        return transactionSplitMapper.toIndividualDto(transactionSplitDao.findWithTagsById(new TransactionSplitId(id, userId)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        ));
    }

    @Override
    public List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId) {
        return transactionSplitDao.getSummaryByYear(year, userId);
    }

    @Override
    @Transactional
    public void create(TransactionDto.Request dto, Long userId) {
        TransactionSplit transactionSplit = transactionSplitMapper.toEntity(dto);
        if (dto.paymentDate() == null) {
            transactionSplit.getTransaction().setPaymentDate(dto.transactionDate());
        }
        transactionSplit.setId(new TransactionSplitId(transactionSplitDao.nextId(), userId));
        long nextTransactionId = transactionDao.nextId();
        transactionSplit.setTransactionId(nextTransactionId);
        transactionSplit.getTransaction().setId(new TransactionId(userId, nextTransactionId));

        transactionSplitDao.save(transactionSplit);
    }

    @Override
    @Transactional
    public void update(Long id, TransactionDto.Request dto, Long userId) {
        var transactionSplit = transactionSplitDao.findById(new TransactionSplitId(id, userId)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        );

        transactionSplitMapper.updateEntityFromDto(dto, transactionSplit);
    }

    @Override
    @Transactional(readOnly = true)
    public List<TransactionDto.UploadResponse> parseTransactionImport(TransactionDto.UploadRequest dto, Long userId) {
        var account =  accountDao.findById(new AccountId(userId, dto.accountId().value())).orElseThrow();
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("dd/MM/yyyy");
        Instant fixedPaymentDate = null;
        if (account.getType().equals("C")) {
            fixedPaymentDate = LocalDate
                    .parse((account.getDueDay() < 10 ? "0"+account.getDueDay(): account.getDueDay())
                                    + "/" +dto.paymentMonth() + "/" + dto.paymentYear(),
                            formatter)
                    .atStartOfDay(ZoneOffset.UTC)
                    .toInstant();

        }
        try(InputStream inputStream = dto.file().getInputStream()) {
            return extractParser.parse(inputStream, dto.fileType(), dto.accountId().value(), userId,
                    identities -> findExisting(identities, userId), fixedPaymentDate);
        } catch (IOException ioex) {
            throw new ApiException(ApiExceptionType.BUSINESS_ERROR,  "Failed to process file.", ioex);
        }
    }

    private Set<TransactionGenericIdentity> findExisting(List<TransactionGenericIdentity> identities, Long userId) {
        if (identities.isEmpty()) return Set.of();
        var dates = identities.stream().map(TransactionGenericIdentity::transactionDate).toList();
        var values = identities.stream().map(TransactionGenericIdentity::value).toList();
        return new HashSet<>(transactionDao.findIdentities(userId, Collections.min(dates), Collections.max(dates), values));
    }

    @Override
    @Transactional
    public void createAll(List<TransactionDto.Request> dtos, Long userId) {
        if (dtos.isEmpty()) return;
        var ids = transactionDao.nextIds(dtos.size()).iterator();
        var splitIds = transactionSplitDao.nextIds(dtos.size()).iterator();
        transactionSplitDao.saveAll(dtos.stream().map(dto -> {
            TransactionSplit transactionSplit = transactionSplitMapper.toEntity(dto);
            if (dto.paymentDate() == null) {
                transactionSplit.getTransaction().setPaymentDate(dto.transactionDate());
            }
            transactionSplit.setId(new TransactionSplitId(splitIds.next(), userId));
            long nextTransactionId = ids.next();
            transactionSplit.setTransactionId(nextTransactionId);
            transactionSplit.getTransaction().setId(new TransactionId(userId, nextTransactionId));

            // Ensures the account and category exist and belong to the user
            return transactionSplit;
        }).toList());
    }

    @Override
    @Transactional
    public void delete(Long id, Long userId) {
        var transactionSplit = transactionSplitDao.findById(new TransactionSplitId(id, userId)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        );

        transactionSplitDao.delete(transactionSplit);
    }
}
