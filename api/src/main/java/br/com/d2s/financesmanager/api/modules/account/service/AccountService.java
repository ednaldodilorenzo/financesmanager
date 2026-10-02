package br.com.d2s.financesmanager.api.modules.account.service;

import br.com.d2s.financesmanager.api.dto.AccountDto;

import java.util.List;

public interface AccountService {
    void create(AccountDto.Request dto, Long userId);

    List<AccountDto.Response> getAll(Long userId);

    void update(Long id, AccountDto.Request dto, Long userId);

    AccountDto.Response getOne(Long id, Long userId);
}
