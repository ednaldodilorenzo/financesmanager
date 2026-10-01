package br.com.d2s.financesmanager.api.modules.planning.service;

import br.com.d2s.financesmanager.api.dto.PlanningDto;

import java.util.List;

public interface PlanningService {
    List<PlanningDto> findPlanningByMonthAndYear(Integer month, Integer year, Long userId);
}
