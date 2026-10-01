package br.com.d2s.financesmanager.api.modules.planning.service;

import br.com.d2s.financesmanager.api.dto.PlanningDto;
import br.com.d2s.financesmanager.api.modules.budget.BudgetDao;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class PlanningServiceImpl implements PlanningService {
    private final BudgetDao budgetDao;

    @Override
    public List<PlanningDto> findPlanningByMonthAndYear(Integer month, Integer year, Long userId) {
        return budgetDao.findPlanningByMonthAndYearBudget(month, year, userId);
    }
}
