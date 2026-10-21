package br.com.d2s.financesmanager.api.modules.tag;

import br.com.d2s.financesmanager.api.model.Transaction;
import br.com.d2s.financesmanager.api.model.TransactionId;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.Repository;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface TagDao extends Repository<Transaction, TransactionId> {
    @Query("select distinct tg from TransactionSplit ts join ts.tags tg where ts.id.userId = :userId and lower(tg) " +
            "like concat('%', lower(:filter), '%')")
    List<String> findByUserIdAndTagLike(@Param("filter") String filter, @Param("userId") Long userId);
}
