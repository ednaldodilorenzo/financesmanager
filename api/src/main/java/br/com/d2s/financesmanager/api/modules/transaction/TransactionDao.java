package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import br.com.d2s.financesmanager.api.model.TransactionId;
import br.com.d2s.financesmanager.api.model.TransactionSpecs;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;

import java.time.Instant;
import java.util.List;

public interface TransactionDao extends JpaRepository<Transaction, TransactionId>, JpaSpecificationExecutor<Transaction> {

    @Query(value = "select nextval('transaction_id_seq')", nativeQuery = true)
    Long nextId();

    Boolean existsByValueAndPaymentDateAndTransactionDateAndId_UserId(Long value, Instant paymentDate, Instant transactionDate, Long userId);

    @Query("""
        select t from Transaction t
        where EXTRACT(YEAR FROM t.paymentDate) = :year
          and EXTRACT(MONTH FROM t.paymentDate) = :month
          and t.id.userId = :userId
    """)
    List<Transaction> findByMonthAndYear(Integer month, Integer year, Long userId);

    @Query("""
        select
            new br.com.d2s.financesmanager.api.dto.TransactionDto$TransactionYearlySummary(
            EXTRACT(MONTH FROM t.paymentDate),
            EXTRACT(YEAR FROM t.paymentDate),
            SUM(t.value),
            t.category.type)
        from Transaction t
        where EXTRACT(YEAR FROM t.paymentDate) = :year
          and t.id.userId = :userId
          group by EXTRACT(MONTH FROM t.paymentDate),
                       EXTRACT(YEAR FROM t.paymentDate),
                       t.category.type
          order by EXTRACT(MONTH FROM t.paymentDate)
    """)
    List<TransactionDto.TransactionYearlySummary> getSummaryByYear(Integer year, Long userId);
}
