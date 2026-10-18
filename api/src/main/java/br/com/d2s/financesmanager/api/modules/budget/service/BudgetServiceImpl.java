package br.com.d2s.financesmanager.api.modules.budget.service;

import br.com.d2s.financesmanager.api.dto.BudgetDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.Budget;
import br.com.d2s.financesmanager.api.model.BudgetId;
import br.com.d2s.financesmanager.api.modules.budget.BudgetDao;
import br.com.d2s.financesmanager.api.modules.budget.BudgetMapper;
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class BudgetServiceImpl implements BudgetService {
    private BudgetDao budgetDao;
    private BudgetMapper budgetMapper;

    @Override
    @Transactional
    public void create(BudgetDto.Request dto, Long userId) {
        Budget budget = budgetMapper.toEntity(dto);
        budget.setId(new BudgetId(budgetDao.nextId(), userId));
        budgetDao.save(budget);
    }

    @Override
    @Transactional
    public void update(Long id, BudgetDto.Request dto, Long userId) {
        var budget = budgetDao.findById(new BudgetId(id, userId)).orElseThrow(() ->
                new ApiException(ApiExceptionType.NOT_FOUND, "Budget not found!"));
        budget.setYear(dto.year());
        budget.setCategoryId(dto.categoryId());
        budget.setValue(dto.value());
    }

    @Override
    public void delete(Long id, Long userId) {
        budgetDao.deleteById(new BudgetId(id, userId));
    }

    public List<BudgetDto.Response> findAll(Integer year, Long userId) {
        return budgetDao.findAllByYearAndId_UserId(year, userId).stream().map(budgetMapper::toDto).toList();
    }

    @Override
    public BudgetDto.Response findById(Long id, Long userId) {
        return budgetMapper
                .toDto(budgetDao
                        .findById(new BudgetId(id, userId))
                        .orElseThrow( () -> new ApiException(ApiExceptionType.NOT_FOUND, "BUDGET_NOT_FOUND")));
    }
}
