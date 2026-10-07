package br.com.d2s.financesmanager.api.modules.budget.service;

import br.com.d2s.financesmanager.api.dto.BudgetDto;
import br.com.d2s.financesmanager.api.dto.CategoryDto;


import java.util.List;

public interface BudgetService {
    void create(BudgetDto.Request dto, Long userId);

    void update(Long id, BudgetDto.Request dto, Long userId);

    void delete(Long id, Long userId);

    List<BudgetDto.Response> findAll(Integer year, Long userId);

    BudgetDto.Response findById(Long id, Long userId);
}
