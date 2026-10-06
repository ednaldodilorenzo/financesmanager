package br.com.d2s.financesmanager.api.modules.account.service;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.AccountId;
import br.com.d2s.financesmanager.api.modules.account.AccountDao;
import br.com.d2s.financesmanager.api.modules.account.AccountMapper;
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class AccountServiceImpl implements AccountService {
    private final AccountDao accountDao;
    private final AccountMapper accountMapper;

    @Override
    @Transactional
    public void create(AccountDto.Request dto, Long userId) {
        var account = accountMapper.toEntity(dto);
        account.setId(new AccountId(userId,accountDao.nextId()));
        accountDao.save(account);
    }

    @Override
    public List<AccountDto.Response> getAll(Long userId) {
        return accountDao
                .findAllById_UserId(userId)
                .stream()
                .map(accountMapper::toDto)
                .toList();
    }

    @Override
    @Transactional
    public void update(Long id, AccountDto.Request dto, Long userId) {
        var account = accountDao
                .findById(new AccountId(userId, id))
                .orElseThrow(() -> new ApiException(ApiExceptionType.NOT_FOUND, "Account not found."));
        account.setType(dto.type());
        account.setName(dto.name());
    }

    @Override
    public AccountDto.Response getOne(Long id, Long userId) {
        return accountMapper
                .toDto(accountDao
                        .findById(new AccountId(id, userId))
                        .orElseThrow(() -> new ApiException(ApiExceptionType.NOT_FOUND, "Account not found.")));
    }
}
