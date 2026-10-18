package br.com.d2s.financesmanager.api.modules.category.service;

import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.infra.exception.ApiException;
import br.com.d2s.financesmanager.api.infra.exception.ApiExceptionType;
import br.com.d2s.financesmanager.api.model.CategoryId;
import br.com.d2s.financesmanager.api.modules.category.CategoryDao;
import br.com.d2s.financesmanager.api.modules.category.CategoryMapper;
import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class CategoryServiceImpl implements CategoryService {

    private final CategoryDao categoryDao;

    private final CategoryMapper categoryMapper;

    @Override
    @Transactional
    public void create(CategoryDto.Request dto, Long userId) {
        var category = categoryMapper.toEntity(dto);
        category.setId(new CategoryId(categoryDao.nextId(), userId));
        category.setFilter(dto.name());
        categoryDao.save(category);
    }

    @Override
    public List<CategoryDto.Response> getAll(Long userId) {
        return categoryDao.findAllById_UserId(userId).stream().map(categoryMapper::toDto).toList();
    }

    @Override
    @Transactional
    public void update(Long id, CategoryDto.Request dto, Long userId) {
        var category = categoryDao.findById(new CategoryId(id, userId)).orElseThrow(() ->
                new ApiException(ApiExceptionType.NOT_FOUND, "Category not found!"));
        category.setName(dto.name());
        category.setFilter(dto.name());
        category.setType(dto.type());
    }

    @Override
    public CategoryDto.Response getOne(Long id, Long userId) {
        return categoryMapper
                .toDto(categoryDao
                        .findById(
                                new CategoryId(id, userId)).orElseThrow(() ->
                                new ApiException(ApiExceptionType.NOT_FOUND, "Category not found!")));
    }
}
