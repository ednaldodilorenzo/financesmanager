package br.com.d2s.financesmanager.api.modules.account;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.model.Account;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface AccountMapper {
    Account toEntity(AccountDto.Request dto);

    @Mapping(source = "id.id", target = "id")
    AccountDto.Response toDto(Account entity);
}
