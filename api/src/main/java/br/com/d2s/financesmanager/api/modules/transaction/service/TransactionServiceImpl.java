package br.com.d2s.financesmanager.api.modules.transaction.service;

import br.com.d2s.financesmanager.api.config.parser.csv.ExtractParser;
import br.com.d2s.financesmanager.api.config.parser.csv.TransactionGenericIdentity;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.*;
import br.com.d2s.financesmanager.api.modules.account.AccountDao;
import br.com.d2s.financesmanager.api.modules.category.CategoryDao;
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
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;


@Service
@AllArgsConstructor
public class TransactionServiceImpl implements TransactionService {
    private final TransactionDao transactionDao;
    private final TransactionSplitDao transactionSplitDao;
    private final TransactionSplitMapper transactionSplitMapper;
    private final TransactionMapper transactionMapper;
    private final AccountDao accountDao;
    private final CategoryDao categoryDao;
    private final ExtractParser extractParser;

    @Override
    public List<TransactionDto.Response> getTransactions(TransactionDto.SearchFilter filter, Long userId) {
        return transactionSplitDao
                .findAll(TransactionSpecs.from(filter, userId))
                .stream()
                .map(transactionSplitMapper::toDto).toList();
    }

    @Override
    public TransactionDto.IndividualResponse getTransactionSplit(Long id, Long userId) {
        return transactionSplitMapper.toIndividualDto(transactionSplitDao.findWithTagsById(new TransactionSplitId(id, userId)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        ));
    }

    @Override
    public TransactionDto.SplitIndividualResponse getTransaction(Long id, Long userId) {
        return transactionSplitMapper.toSplitIndividualDto(transactionDao
                .findWithSplitsById(new TransactionId(userId, id))
                .orElseThrow(() -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")));
    }

    @Override
    public List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId) {
        return transactionSplitDao.getSummaryByYear(year, userId);
    }

    @Override
    @Transactional
    public void create(TransactionDto.Request dto, Long userId) {
        if (!categoryDao.existsById(new CategoryId(dto.categoryId().value(), userId))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND,
                    "Category %d not found".formatted(dto.categoryId().value()));
        }

        if (!accountDao.existsById(new AccountId(userId, dto.accountId().value()))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND, "Account %d not found".formatted(dto.accountId().value()));
        }

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
    public void createSplitedTransactions(TransactionDto.SplitedRequest dto, Long userId) {
        validateSplitedRequest(dto, userId);
        var transaction = transactionSplitMapper.toTransaction(dto);
        transaction.setId(new TransactionId(userId, transactionDao.nextId()));
        transactionDao.save(transaction);

        var splitIds = transactionSplitDao.nextIds(dto.splits().size()).iterator();
        transactionSplitDao.saveAll(dto.splits().stream().map(split -> {
            var transactionSplit = transactionSplitMapper.toEntity(split, transaction);
            transactionSplit.setId(new TransactionSplitId(splitIds.next(), userId));
            return transactionSplit;
        }).toList());
    }

    @Override
    @Transactional
    public void updateSplitedTransaction(Long id, TransactionDto.SplitedRequest dto, Long userId) {
        validateSplitedRequest(dto, userId);

        var transaction = transactionDao.findWithSplitsById(new TransactionId(userId, id)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        );
        transactionMapper.updateEntityFromDto(dto, transaction);

        // Splits left in this map after the loop were removed by the client
        var remainingSplits = transaction.getSplits().stream()
                .collect(Collectors.toMap(s -> s.getId().getId(), Function.identity()));

        long newSplitsCount = dto.splits().stream().filter(s -> s.id() == null).count();
        Iterator<Long> splitIds = newSplitsCount == 0
                ? Collections.emptyIterator()
                : transactionSplitDao.nextIds((int) newSplitsCount).iterator();

        for (var split : dto.splits()) {
            if (split.id() == null) {
                var transactionSplit = transactionSplitMapper.toEntity(split, transaction);
                transactionSplit.setId(new TransactionSplitId(splitIds.next(), userId));
                transactionSplitDao.save(transactionSplit);
            } else {
                var transactionSplit = remainingSplits.remove(split.id().value());
                if (transactionSplit == null) {
                    throw new ApiException(ApiExceptionType.NOT_FOUND,
                            "Transaction split %d not found".formatted(split.id().value()));
                }
                transactionSplitMapper.updateSplitFromDto(split, transactionSplit);
            }
        }

        transaction.getSplits().removeAll(remainingSplits.values());
        transactionSplitDao.deleteAll(remainingSplits.values());
    }

    private void validateSplitedRequest(TransactionDto.SplitedRequest dto, Long userId) {
        if (dto.value() != dto.splits().stream().mapToLong(TransactionDto.TransactionSplit::value).sum()) {
            throw new ApiException(ApiExceptionType.BUSINESS_ERROR, "Split value sum differs from transaction value");
        }

        if (!accountDao.existsById(new AccountId(userId, dto.accountId().value()))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND, "Account %d not found".formatted(dto.accountId().value()));
        }

        var requestedCategoryIds =
                dto.splits().stream().map(s -> s.categoryId().value()).collect(Collectors.toSet());
        validateMissingIds(requestedCategoryIds, "Category",
                (idSet -> categoryDao.findExistingIdsByIdInAndUserId(idSet, userId)));
    }

    @Override
    @Transactional
    public void update(Long id, TransactionDto.Request dto, Long userId) {
        if (!categoryDao.existsById(new CategoryId(dto.categoryId().value(), userId))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND,
                    "Category %d not found".formatted(dto.categoryId().value()));
        }

        if (!accountDao.existsById(new AccountId(dto.accountId().value(), userId))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND, "Account %d not found".formatted(dto.accountId().value()));
        }

        var transactionSplit = transactionSplitDao.findById(new TransactionSplitId(id, userId)).orElseThrow(
                () -> new ApiException(ApiExceptionType.NOT_FOUND, "Transaction not found")
        );

        transactionSplitMapper.updateEntityFromDto(dto, transactionSplit);
    }

    @Override
    @Transactional(readOnly = true)
    public List<TransactionDto.UploadResponse> parseTransactionImport(TransactionDto.UploadRequest dto, Long userId) {
        if (!accountDao.existsById(new AccountId(userId, dto.accountId().value()))) {
            throw new ApiException(ApiExceptionType.NOT_FOUND, "Account %d not found".formatted(dto.accountId().value()));
        }
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

        var requestedCategoryIds =
                dtos.stream().map(t -> t.categoryId().value()).collect(Collectors.toSet());
        validateMissingIds(requestedCategoryIds, "Category",
                (idSet -> categoryDao.findExistingIdsByIdInAndUserId(idSet, userId)));

        var requestedAccountIds =
                dtos.stream().map(t -> t.accountId().value()).collect(Collectors.toSet());
        validateMissingIds(requestedAccountIds, "Account",
                (idSet -> accountDao.findExistingIdsByIdInAndUserId(idSet, userId)));

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

        // A transaction without splits is meaningless, so it goes away with its last split
        if (!transactionSplitDao.existsByTransactionIdAndId_UserId(transactionSplit.getTransactionId(), userId)) {
            transactionDao.deleteById(new TransactionId(userId, transactionSplit.getTransactionId()));
        }
    }

    public void validateMissingIds(Set<Long> requestedIds, String entityName,
                                   Function<Set<Long>, Set<Long>> findExisting) {
        if (requestedIds == null || requestedIds.isEmpty()) {
            return;
        }

        // 1. Busca no banco de dados apenas os IDs que realmente existem
        Set<Long> existingIds = findExisting.apply(requestedIds);

        // 2. Cria uma cópia dos IDs requisitados para não alterar o parâmetro original
        Set<Long> missingIds = new HashSet<>(requestedIds);

        // 3. Remove todos os que existem. Sobram apenas os inexistentes!
        missingIds.removeAll(existingIds);

        if (!missingIds.isEmpty()) {
            throw new ApiException(ApiExceptionType.NOT_FOUND,
                    "%s id(s) not found: ".formatted(entityName) +  missingIds);
        }
    }
}
