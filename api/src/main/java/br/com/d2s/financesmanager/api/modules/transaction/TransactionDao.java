package br.com.d2s.financesmanager.api.modules.transaction;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface TransactionDao extends JpaRepository<Transaction, Long> {

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
