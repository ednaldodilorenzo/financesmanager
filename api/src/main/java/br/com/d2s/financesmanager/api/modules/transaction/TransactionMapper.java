package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface TransactionMapper {
    @Mapping(source = "id.id", target = "id")
    @Mapping(source = "account.id.id", target="account.id")
    @Mapping(source = "category.categoryId.id", target="category.id")
    TransactionDto.Response toDto(Transaction transaction);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "paymentDate", ignore = true)
    @Mapping(target = "account", ignore = true)
    @Mapping(target = "category", ignore = true)
    Transaction toEntity(TransactionDto.Request request);
}
