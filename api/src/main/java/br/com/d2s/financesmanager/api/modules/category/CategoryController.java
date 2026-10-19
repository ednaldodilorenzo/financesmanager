package br.com.d2s.financesmanager.api.modules.category;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.CategoryDto;
import br.com.d2s.financesmanager.api.modules.category.service.CategoryService;
import jakarta.validation.Valid;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/categories")
@AllArgsConstructor
public class CategoryController {

    private final CategoryService categoryService;

    @GetMapping("/")
    public List<CategoryDto.Response> findAll(@AuthenticationPrincipal Jwt jwt) {
        return categoryService.getAll(Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @GetMapping("/{id}")
    public CategoryDto.Response findById(@PathVariable PublicId id, @AuthenticationPrincipal Jwt jwt) {
        return categoryService.getOne(id.value(), Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public void post(@RequestBody @Valid CategoryDto.Request dto, @AuthenticationPrincipal Jwt jwt) {
        categoryService.create(dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }

    @PatchMapping("/{id}")
    public void patch(@PathVariable PublicId id, @AuthenticationPrincipal Jwt jwt,
                      @RequestBody @Valid CategoryDto.Request dto) {
        categoryService.update(id.value(), dto, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}
