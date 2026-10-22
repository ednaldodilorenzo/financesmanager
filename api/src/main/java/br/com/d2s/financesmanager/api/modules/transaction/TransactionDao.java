package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.config.parser.csv.TransactionGenericIdentity;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import br.com.d2s.financesmanager.api.model.TransactionId;
import br.com.d2s.financesmanager.api.model.TransactionSplit;
import br.com.d2s.financesmanager.api.model.TransactionSplitId;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;

import java.time.Instant;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

public interface TransactionDao extends JpaRepository<Transaction, TransactionId>, JpaSpecificationExecutor<Transaction> {

    @Query(value = "select nextval('transaction_id_seq')", nativeQuery = true)
    Long nextId();

    // TransactionDao
    @Query(value = "select nextval('transaction_id_seq') from generate_series(1, :count)", nativeQuery = true)
    List<Long> nextIds(int count);

    // Loads tags with the transaction so they can be mapped after the transaction closes (open-in-view is off).
    @EntityGraph(attributePaths = {"account", "splits", "splits.category", "splits.tags"})
    Optional<Transaction> findWithSplitsById(TransactionId id);

//    Boolean existsByValueAndPaymentDateAndTransactionDateAndId_UserId(Long value, Instant paymentDate, Instant transactionDate, Long userId);

//    @Query("""
//        select t from Transaction t
//        where EXTRACT(YEAR FROM t.paymentDate) = :year
//          and EXTRACT(MONTH FROM t.paymentDate) = :month
//          and t.id.userId = :userId
//    """)
//    List<Transaction> findByMonthAndYear(Integer month, Integer year, Long userId);

    @Query("""
        select new br.com.d2s.financesmanager.api.config.parser.csv.TransactionGenericIdentity(
            t.value, t.paymentDate, t.transactionDate, t.id.userId
        )
        from Transaction t
        where t.value in :values
          and t.id.userId = :userId
          and t.transactionDate between :from and :to
    """)
    List<TransactionGenericIdentity> findIdentities(Long userId, Instant from, Instant to, Collection<Long> values);

//    @Query("""
//        select
//            new br.com.d2s.financesmanager.api.dto.TransactionDto$TransactionYearlySummary(
//            EXTRACT(MONTH FROM t.paymentDate),
//            EXTRACT(YEAR FROM t.paymentDate),
//            SUM(t.value),
//            t.category.type)
//        from Transaction t
//        where EXTRACT(YEAR FROM t.paymentDate) = :year
//          and t.id.userId = :userId
//          group by EXTRACT(MONTH FROM t.paymentDate),
//                       EXTRACT(YEAR FROM t.paymentDate),
//                       t.category.type
//          order by EXTRACT(MONTH FROM t.paymentDate)
//    """)
//    List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId);
}
