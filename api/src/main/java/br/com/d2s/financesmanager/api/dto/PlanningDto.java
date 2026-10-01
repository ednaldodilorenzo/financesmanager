package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.model.CategoryType;

public record PlanningDto(
        String name,
        CategoryType type,
        Long total,
        Long planned,
        Long accumulated) {
}
