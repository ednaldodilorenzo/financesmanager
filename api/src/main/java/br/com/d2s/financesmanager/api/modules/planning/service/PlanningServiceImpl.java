package br.com.d2s.financesmanager.api.modules.planning.service;

import br.com.d2s.financesmanager.api.dto.PlanningDto;
import br.com.d2s.financesmanager.api.modules.budget.BudgetDao;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.ZoneId;
import java.util.List;

@Service
@AllArgsConstructor
public class PlanningServiceImpl implements PlanningService {
    private final BudgetDao budgetDao;

    @Override
    public List<PlanningDto> findPlanningByMonthAndYear(Integer month, Integer year, Long userId) {
        // Date ranges instead of extract(year/month) so the payment_date index can be used.
        // payment_date is stored in the JVM local time zone, so boundaries are computed in it too.
        ZoneId zone = ZoneId.systemDefault();
        LocalDate monthStart = LocalDate.of(year, month, 1);
        return budgetDao.findPlanningByMonthAndYearBudget(
                year,
                monthStart.withDayOfYear(1).atStartOfDay(zone).toInstant(),
                monthStart.atStartOfDay(zone).toInstant(),
                monthStart.plusMonths(1).atStartOfDay(zone).toInstant(),
                userId);
    }
}
