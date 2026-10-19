package br.com.d2s.financesmanager.api.infra.publicid;

import br.com.d2s.financesmanager.api.dto.AccountDto;
import br.com.d2s.financesmanager.api.infra.exception.GlobalExceptionHandler;
import br.com.d2s.financesmanager.api.modules.account.AccountController;
import br.com.d2s.financesmanager.api.modules.account.service.AccountService;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;
import static org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors.csrf;
import static org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors.jwt;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.patch;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(AccountController.class)
@Import({PublicIdConfig.class, PublicIdCodec.class, GlobalExceptionHandler.class})
@TestPropertySource(properties = "APP_PUBLIC_ID_ALPHABET=KIHZiM4R6o1p2hDOf8YXLeAJjw0UnPkgzdQaTGsNWmFyBcVqEr7uv5t3S9Clbx")
class PublicIdWebTest {
    @Autowired
    private MockMvc mvc;

    @Autowired
    private PublicIdCodec codec;

    @MockitoBean
    private AccountService accountService;

    @Test
    void decodesPathVariableAndEncodesResponseId() throws Exception {
        when(accountService.getOne(15L, 1L)).thenReturn(new AccountDto.Response(new PublicId(15L), "Bank", "C", 10));

        mvc.perform(get("/accounts/{id}", codec.encode(15L)).with(jwt().jwt(j -> j.subject("1"))))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(codec.encode(15L)));
    }

    @Test
    void undecodablePathIdIsNotFound() throws Exception {
        mvc.perform(get("/accounts/{id}", "15").with(jwt().jwt(j -> j.subject("1"))))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.detail").value("Resource not found"));

        verifyNoInteractions(accountService);
    }

    @Test
    void patchDecodesPathVariable() throws Exception {
        mvc.perform(patch("/accounts/{id}", codec.encode(15L))
                        .with(jwt().jwt(j -> j.subject("1"))).with(csrf())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"name\":\"Bank\",\"type\":\"C\",\"dueDay\":10}"))
                .andExpect(status().isOk());

        verify(accountService).update(eq(15L), any(), eq(1L));
    }
}
