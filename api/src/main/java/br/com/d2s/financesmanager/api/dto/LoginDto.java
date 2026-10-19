package br.com.d2s.financesmanager.api.dto;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
public class LoginDto {

    public record Request(String email, String password) {
    }

    public record Response(PublicId id, String name, String token) {
    }
}
