package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.model.CategoryType;

public record CategoryDto(Long id, String name, CategoryType type) {
}
