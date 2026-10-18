package br.com.d2s.financesmanager.api.modules.category;

import br.com.d2s.financesmanager.api.model.Category;
import br.com.d2s.financesmanager.api.model.CategoryId;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface CategoryDao extends JpaRepository<Category, CategoryId> {
    @Query(value = "select nextval('category_id_seq')", nativeQuery = true)
    Long nextId();

    List<Category> findAllById_UserId(Long userId);
}
