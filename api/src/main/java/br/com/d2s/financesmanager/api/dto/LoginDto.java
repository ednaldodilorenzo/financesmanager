package br.com.d2s.financesmanager.api.dto;

public class LoginDto {

    public record Request(String email, String password) {
    }

    public record Response(Long id, String name, String token) {
    }
}
