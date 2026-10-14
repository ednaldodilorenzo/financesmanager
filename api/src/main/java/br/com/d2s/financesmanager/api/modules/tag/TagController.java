package br.com.d2s.financesmanager.api.modules.tag;

import br.com.d2s.financesmanager.api.dto.TransactionTagDto;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Objects;

@RestController
@RequestMapping("/v1/tags")
@AllArgsConstructor
@Slf4j
public class TagController {
    @Autowired
    private TagService tagService;

    @GetMapping("/")
    public List<TransactionTagDto> findByFilter(@RequestParam String filter, @AuthenticationPrincipal Jwt jwt) {
        log.info("Finding by filter: {}", filter);
        return tagService.findByFilter(filter, Long.parseLong(Objects.requireNonNull(jwt.getSubject())));
    }
}

