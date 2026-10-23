package br.com.d2s.financesmanager.api.infra.token;

import org.springframework.security.core.Authentication;

public interface JwtManager {
    String generateAccessToken(Authentication authentication);
    String generateRegistrationToken(String email);
    String extractRegistrationInfo(String token);
}
