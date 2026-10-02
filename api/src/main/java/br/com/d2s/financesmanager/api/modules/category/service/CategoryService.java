package br.com.d2s.financesmanager.api.modules.category.service;

import br.com.d2s.financesmanager.api.dto.CategoryDto;

import java.util.List;

public interface CategoryService {

    void create(CategoryDto.Request dto, Long userId);

    List<CategoryDto.Response> getAll(Long userId);

    void update(Long id, CategoryDto.Request dto, Long userId);

    CategoryDto.Response getOne(Long id, Long userId);
}
