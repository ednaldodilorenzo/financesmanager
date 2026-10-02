package br.com.d2s.financesmanager.api.modules.account;

import br.com.d2s.financesmanager.api.model.Account;
import br.com.d2s.financesmanager.api.model.AccountId;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface AccountDao extends JpaRepository<Account, AccountId> {
    List<Account> findAllById_UserId(Long userId);
}
