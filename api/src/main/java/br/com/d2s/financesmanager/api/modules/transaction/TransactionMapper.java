package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.infra.publicid.PublicIdMapper;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.dto.TransactionTagDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingTarget;

@Mapper(componentModel = "spring", uses = PublicIdMapper.class)
public interface TransactionMapper {
    @Mapping(source = "id.id", target = "id")
    @Mapping(source = "account.id.id", target="account.id")
    TransactionDto.Response toDto(Transaction transaction);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "account", ignore = true)
    Transaction toEntity(TransactionDto.Request request);

    // Splits are synchronized by the service, since they need ids and may be added or removed.
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "account", ignore = true)
    @Mapping(target = "splits", ignore = true)
    @Mapping(target = "paymentDate", defaultExpression = "java(request.transactionDate())")
    void updateEntityFromDto(TransactionDto.SplitedRequest request, @MappingTarget Transaction transaction);
}
