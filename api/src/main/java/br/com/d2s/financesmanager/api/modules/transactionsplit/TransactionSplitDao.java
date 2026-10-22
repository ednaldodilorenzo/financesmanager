package br.com.d2s.financesmanager.api.modules.transactionsplit;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import br.com.d2s.financesmanager.api.model.TransactionId;
import br.com.d2s.financesmanager.api.model.TransactionSplit;
import br.com.d2s.financesmanager.api.model.TransactionSplitId;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface TransactionSplitDao extends JpaRepository<TransactionSplit, TransactionSplitId>, JpaSpecificationExecutor<TransactionSplit> {
    @Query(value = "select nextval('transaction_split_id_seq')", nativeQuery = true)
    Long nextId();

    // TransactionDao
    @Query(value = "select nextval('transaction_split_id_seq') from generate_series(1, :count)", nativeQuery = true)
    List<Long> nextIds(int count);

    boolean existsByTransactionIdAndId_UserId(Long transactionId, Long userId);

    @EntityGraph(attributePaths = {"tags", "category", "transaction", "transaction.account"})
    Optional<TransactionSplit> findWithTagsById(TransactionSplitId id);

    @Query("""
        select ts from TransactionSplit ts
            join ts.transaction t
        where EXTRACT(YEAR FROM t.paymentDate) = :year
          and EXTRACT(MONTH FROM t.paymentDate) = :month
          and ts.id.userId = :userId
    """)
    List<TransactionSplit> findByMonthAndYear(Integer month, Integer year, Long userId);

    @Query("""
        select
            new br.com.d2s.financesmanager.api.dto.TransactionDto$TransactionYearlySummary(
            EXTRACT(MONTH FROM t.paymentDate),
            EXTRACT(YEAR FROM t.paymentDate),
            SUM(ts.value),
            ts.category.type)
        from TransactionSplit ts
            join ts.transaction t
        where EXTRACT(YEAR FROM t.paymentDate) = :year
          and ts.id.userId = :userId
          group by EXTRACT(MONTH FROM t.paymentDate),
                       EXTRACT(YEAR FROM t.paymentDate),
                       ts.category.type
          order by EXTRACT(MONTH FROM t.paymentDate)
    """)
    List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId);
}
