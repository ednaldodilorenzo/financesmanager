package br.com.d2s.financesmanager.api.modules.budget;

import br.com.d2s.financesmanager.api.dto.PlanningDto;
import br.com.d2s.financesmanager.api.model.Budget;
import br.com.d2s.financesmanager.api.model.BudgetId;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface BudgetDao extends JpaRepository<Budget, BudgetId> {

    @Query(value = "select nextval('budget_id_seq')", nativeQuery = true)
    Long nextId();

    @Query("""
    select new br.com.d2s.financesmanager.api.dto.PlanningDto(
        c.name,
        c.type,
        coalesce((select sum(ts.value)
                  from TransactionSplit ts
                      join ts.transaction t
                  where ts.category = c
                    and extract(year from t.paymentDate) = :year
                    and extract(month from t.paymentDate) = :month), 0L),
        b.value,
        coalesce((select sum(ts.value)
                  from TransactionSplit ts
                          join ts.transaction t
                  where ts.category = c
                    and extract(year from t.paymentDate) = :year
                    and extract(month from t.paymentDate) <= :month), 0L)
    )
    from Budget b
    join b.category c
    where b.year = :year
      and b.id.userId = :userId
    """)
    List<PlanningDto> findPlanningByMonthAndYearBudget(Integer month, Integer year, Long userId);

    List<Budget> findAllByYearAndId_UserId(Integer year, Long userId);
}
