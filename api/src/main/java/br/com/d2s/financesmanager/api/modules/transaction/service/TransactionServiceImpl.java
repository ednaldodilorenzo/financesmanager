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
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

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
    private final TransactionMapper transactionMapper;
    private final AccountDao accountDao;
    private final ExtractParser extractParser;

    @Override
    public List<TransactionDto.Response> getTransactions(TransactionDto.SearchFilter filter, Long userId) {
        return transactionDao
                .findAll(TransactionSpecs.from(filter, userId))
                .stream()
                .map(transactionMapper::toDto).toList();
    }

    @Override
    public TransactionDto.IndividualResponse getTransaction(Long id, Long userId) {
        return transactionMapper.toIndividualDto(transactionDao.findById(new TransactionId(userId, id)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        ));
    }

    @Override
    public List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId) {
        return transactionDao.getSummaryByYear(year, userId);
    }

    @Override
    @Transactional
    public void create(TransactionDto.Request dto, Long userId) {
        Transaction transaction = transactionMapper.toEntity(dto);
        if (dto.paymentDate() == null) {
            transaction.setPaymentDate(dto.transactionDate());
        }
        IO.println(transaction.getPaymentDate());
        transaction.setId(new TransactionId(userId, transactionDao.nextId()));
        // Ensures the account and category exist and belong to the user
        transactionDao.save(transaction);
    }

    @Override
    @Transactional
    public void update(Long id, TransactionDto.Request dto, Long userId) {
        var transaction = transactionDao.findById(new TransactionId(userId, id)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        );

        transactionMapper.updateEntityFromDto(dto, transaction);
    }

    @Override
    public List<TransactionDto.UploadResponse> parseTransactionImport(TransactionDto.UploadRequest dto, Long userId) {
        var account =  accountDao.findById(new AccountId(userId, dto.accountId())).orElseThrow();
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
            return extractParser.parse(inputStream, dto.fileType(), dto.accountId(), userId,
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
        transactionDao.saveAll(dtos.stream().map(dto -> {
            Transaction transaction = transactionMapper.toEntity(dto);
            if (dto.paymentDate() == null) {
                transaction.setPaymentDate(dto.transactionDate());
            }
            transaction.setId(new TransactionId(userId, transactionDao.nextId()));
            // Ensures the account and category exist and belong to the user
            return transaction;
        }).toList());
    }
}
