package br.com.d2s.financesmanager.api.infra.security;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;

public class CookieBearerTokenResolver implements BearerTokenResolver {

    private static final String COOKIE_NAME = "jwt";

    @Override
    public String resolve(HttpServletRequest request) {
//        if (request.getServletPath().startsWith("/auth/")) {
//            return null;
//        }

        Cookie[] cookies = request.getCookies();
        if (cookies == null) {
            return null;
        }

        for (Cookie cookie : cookies) {
            if (COOKIE_NAME.equals(cookie.getName())) {
                String token = cookie.getValue();

                if (token == null || token.isBlank()) {
                    return null;
                }

                if (token.startsWith("Bearer ")) {
                    return token.substring(7);
                }

                return token;
            }
        }

        return null;
    }
}
