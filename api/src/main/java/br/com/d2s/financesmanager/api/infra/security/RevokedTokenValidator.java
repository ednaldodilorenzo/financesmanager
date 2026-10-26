package br.com.d2s.financesmanager.api.infra.security;

import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.core.OAuth2TokenValidator;
import org.springframework.security.oauth2.core.OAuth2TokenValidatorResult;
import org.springframework.security.oauth2.jwt.Jwt;

@RequiredArgsConstructor
public class RevokedTokenValidator implements OAuth2TokenValidator<Jwt> {
    private final TokenRevocationStore tokenRevocationStore;

    @Override
    public OAuth2TokenValidatorResult validate(Jwt token) {
        String jti = token.getId();
        if (jti == null ||  tokenRevocationStore.isRevoked(jti)) {
            return OAuth2TokenValidatorResult.failure();
        }
        return OAuth2TokenValidatorResult.success();
    }
}
