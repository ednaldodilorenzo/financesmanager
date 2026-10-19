package br.com.d2s.financesmanager.api.infra.exception;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import lombok.extern.slf4j.Slf4j;
import org.jspecify.annotations.Nullable;
import org.springframework.beans.TypeMismatchException;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ProblemDetail;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

@Slf4j
@ControllerAdvice
public class GlobalExceptionHandler extends ResponseEntityExceptionHandler {

    @ExceptionHandler(ApiException.class)
    public ProblemDetail handleApiException(ApiException exception) {
        var statusCode = switch (exception.getType()) {
            case NOT_FOUND -> HttpStatus.NOT_FOUND;
            case BAD_REQUEST -> HttpStatus.BAD_REQUEST;
            case BUSINESS_ERROR -> HttpStatus.UNPROCESSABLE_CONTENT;
            case CONFLICT -> HttpStatus.CONFLICT;
        };

        log.warn("ApiException [{}]: {}", exception.getType(), exception.getMessage());
        return ProblemDetail.forStatusAndDetail(statusCode, exception.getMessage());
    }

    @ExceptionHandler(Exception.class)
    public ProblemDetail handleUnexpectedException(Exception exception) {
        log.error("Unexpected error", exception);
        return ProblemDetail.forStatusAndDetail(HttpStatus.INTERNAL_SERVER_ERROR, "Unexpected internal error");
    }

    @Override
    protected @Nullable ResponseEntity<Object> handleTypeMismatch(
            TypeMismatchException ex, HttpHeaders headers, HttpStatusCode status, WebRequest request) {

        // An undecodable id in the URL can't match any resource; answer the same as an id owned by another user.
        if (ex.getRequiredType() == PublicId.class) {
            var body = ProblemDetail.forStatusAndDetail(HttpStatus.NOT_FOUND, "Resource not found");
            return handleExceptionInternal(ex, body, headers, HttpStatus.NOT_FOUND, request);
        }
        return super.handleTypeMismatch(ex, headers, status, request);
    }

    @Override
    protected @Nullable ResponseEntity<Object> handleExceptionInternal(
            Exception ex, @Nullable Object body, HttpHeaders headers,
            HttpStatusCode statusCode, WebRequest request) {

        if (statusCode.is5xxServerError()) {
            log.error("Error processing {}: {}", request.getDescription(false), ex.getMessage(), ex);
        } else {
            log.warn("Invalid request {}: {}", request.getDescription(false), ex.getMessage());
        }

        return super.handleExceptionInternal(ex, body, headers, statusCode, request);
    }
}
