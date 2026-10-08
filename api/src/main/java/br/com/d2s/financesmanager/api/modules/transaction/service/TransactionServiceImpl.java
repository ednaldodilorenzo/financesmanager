package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.config.parser.csv.ExtractParser;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.*;
import br.com.d2s.financesmanager.api.modules.account.AccountDao;
import br.com.d2s.financesmanager.api.modules.category.CategoryDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionDao;
import br.com.d2s.financesmanager.api.modules.transaction.TransactionMapper;
import com.opencsv.CSVReader;
import com.opencsv.CSVReaderBuilder;
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.io.*;
import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.List;



@Service
@AllArgsConstructor
public class TransactionServiceImpl implements TransactionService {
    private final TransactionDao transactionDao;
    private final TransactionMapper transactionMapper;
    private final AccountDao accountDao;
    private final ExtractParser extractParser;

    @Override
    public List<TransactionDto.Response> getTransactions(Integer month, Integer year, Long userId) {
        return transactionDao.findByMonthAndYear(month, year, userId)
                .stream()
                .map(transactionMapper::toDto)
                .toList();
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
                    t -> transactionDao.existsByValueAndPaymentDateAndTransactionDateAndId_UserId(t.value(),
                            t.paymentDate(), t.transactionDate(), t.userId()), fixedPaymentDate);
        } catch (IOException ioex) {
            throw new ApiException(ApiExceptionType.UNPROCESSABLE_ENTITY,  "Failed to process file.", ioex);
        }
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
