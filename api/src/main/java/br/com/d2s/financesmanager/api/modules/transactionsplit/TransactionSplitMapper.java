package br.com.d2s.financesmanager.api.modules.transactionsplit;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.dto.TransactionTagDto;
import br.com.d2s.financesmanager.api.infra.publicid.PublicIdMapper;
import br.com.d2s.financesmanager.api.model.Account;
import br.com.d2s.financesmanager.api.model.Category;
import br.com.d2s.financesmanager.api.model.Transaction;
import br.com.d2s.financesmanager.api.model.TransactionSplit;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingTarget;

@Mapper(componentModel = "spring", uses = PublicIdMapper.class)
public interface TransactionSplitMapper {
    @Mapping(source = "id.id", target = "id")
    @Mapping(source = "transaction.id.id", target = "transactionId")
    @Mapping(source = "transaction.account", target = "account")
    @Mapping(source = "transaction.paymentDate", target = "paymentDate")
    @Mapping(source = "transaction.transactionDate", target = "transactionDate")
    TransactionDto.Response toDto(TransactionSplit transactionSplit);

    @Mapping(source = "id.id", target = "id")
    @Mapping(source = "transaction.account", target = "account")
    @Mapping(source = "transaction.paymentDate", target = "paymentDate")
    @Mapping(source = "transaction.transactionDate", target = "transactionDate")
    TransactionDto.IndividualResponse toIndividualDto(TransactionSplit transactionSplit);

    @Mapping(source = "id.id", target = "id")
    TransactionDto.SplitIndividualResponse toSplitIndividualDto(Transaction transaction);

    @Mapping(source = "id.id", target = "id")
    TransactionDto.TransactionSplitResponse toSplitResponse(TransactionSplit transactionSplit);

    @Mapping(target = "id", ignore = true)
    @Mapping(source = "accountId", target = "transaction.accountId")
    @Mapping(source = "transactionDate", target = "transaction.transactionDate")
    @Mapping(source = "paymentDate", target = "transaction.paymentDate")
    @Mapping(source = "value", target = "transaction.value")
    TransactionSplit toEntity(TransactionDto.Request request);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "account", ignore = true)
    @Mapping(target = "splits", ignore = true)
    @Mapping(target = "paymentDate", defaultExpression = "java(request.transactionDate())")
    Transaction toTransaction(TransactionDto.SplitedRequest request);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "category", ignore = true)
    @Mapping(source = "transaction", target = "transaction")
    @Mapping(source = "transaction.id.id", target = "transactionId")
    @Mapping(source = "split.value", target = "value")
    TransactionSplit toEntity(TransactionDto.TransactionSplit split, Transaction transaction);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "transactionId", ignore = true)
    @Mapping(target = "transaction", ignore = true)
    @Mapping(target = "category", ignore = true)
    void updateSplitFromDto(TransactionDto.TransactionSplit split, @MappingTarget TransactionSplit transactionSplit);

    @Mapping(source = "id.id", target = "id")
    AccountDto.Response toAccountDto(Account account);

    @Mapping(source = "id.id", target = "id")
    CategoryDto.Response toCategoryDto(Category category);

    default TransactionTagDto toTagDto(String tag) {
        return new TransactionTagDto(tag);
    }

    default String toTag(TransactionTagDto tagDto) {
        return tagDto.tag();
    }

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "transactionId", ignore = true)
    @Mapping(source = "accountId", target = "transaction.accountId")
    @Mapping(source = "transactionDate", target = "transaction.transactionDate")
    @Mapping(source = "paymentDate", target = "transaction.paymentDate")
    @Mapping(source = "value", target = "transaction.value")
    void updateEntityFromDto(TransactionDto.Request request, @MappingTarget TransactionSplit transactionSplit);
}
