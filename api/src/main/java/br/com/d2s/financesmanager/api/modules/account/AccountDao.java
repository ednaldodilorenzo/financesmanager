package br.com.d2s.financesmanager.api.modules.account;

import br.com.d2s.financesmanager.api.model.Account;
import br.com.d2s.financesmanager.api.model.AccountId;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Set;

@Repository
public interface AccountDao extends JpaRepository<Account, AccountId> {
    @Query(value = "select nextval('account_id_seq')", nativeQuery = true)
    Long nextId();

    List<Account> findAllById_UserId(Long userId);

    @Query("SELECT a.id.id FROM Account a WHERE a.id.id IN :accountIds AND a.id.userId = :userId")
    Set<Long> findExistingIdsByIdInAndUserId(
            @Param("accountIds") Set<Long> accountIds,
            @Param("userId") Long userId
    );
}
