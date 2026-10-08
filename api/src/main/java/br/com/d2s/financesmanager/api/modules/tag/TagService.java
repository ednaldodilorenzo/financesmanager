package br.com.d2s.financesmanager.api.modules.tag;

import br.com.d2s.financesmanager.api.dto.TransactionTagDto;

import java.util.List;

public interface TagService {
    List<TransactionTagDto> findByFilter(String filter, Long userId);
}
