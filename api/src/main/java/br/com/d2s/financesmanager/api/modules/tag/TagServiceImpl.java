package br.com.d2s.financesmanager.api.modules.tag;

import br.com.d2s.financesmanager.api.dto.TransactionTagDto;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class TagServiceImpl implements TagService {
    private TagDao tagDao;
    @Override
    public List<TransactionTagDto> findByFilter(String filter, Long userId) {
        return tagDao.findByUserIdAndTagLike(filter, userId)
                .stream()
                .map(TransactionTagDto::new).toList();
    }
}
