package br.com.d2s.financesmanager.api.model;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import jakarta.persistence.criteria.Expression;
import jakarta.persistence.criteria.JoinType;
import org.springframework.data.jpa.domain.Specification;

import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneId;
import java.util.ArrayList;
import jakarta.persistence.criteria.Predicate;

public final class TransactionSpecs {
    public static Specification<TransactionSplit> from(TransactionDto.SearchFilter filter, Long userId) {
        return (root, query, cb) -> {
            if (query.getResultType() != long.class && query.getResultType() != Long.class) {
                root.fetch("category", JoinType.INNER);
                root.fetch("transaction", JoinType.INNER).fetch("account", JoinType.INNER);
            }
            var predicates = new ArrayList<Predicate>();
            predicates.add(cb.equal(root.get("id").get("userId"), userId));
            Expression<Instant> paymentDate = root.get("transaction").get("paymentDate");
            // Date ranges instead of date_part so the payment_date index can be used.
            // payment_date is stored in the JVM local time zone, so boundaries are computed in it too.
            if (filter.year() != null) {
                ZoneId zone = ZoneId.systemDefault();
                LocalDate start = LocalDate.of(filter.year(), filter.month() != null ? filter.month() : 1, 1);
                LocalDate end = filter.month() != null ? start.plusMonths(1) : start.plusYears(1);
                predicates.add(cb.greaterThanOrEqualTo(paymentDate, start.atStartOfDay(zone).toInstant()));
                predicates.add(cb.lessThan(paymentDate, end.atStartOfDay(zone).toInstant()));
            } else if (filter.month() != null)
                predicates.add(cb.equal(cb.function("date_part", Integer.class, cb.literal("month"), paymentDate), filter.month()));
            if (filter.startDate() != null)
                predicates.add(cb.greaterThanOrEqualTo(paymentDate, filter.startDate()));
            if (filter.endDate() != null)
                predicates.add(cb.lessThanOrEqualTo(paymentDate, filter.endDate()));
            if (filter.description() != null && !filter.description().isBlank())
                predicates.add(cb.like(cb.lower(root.get("description")), "%" + filter.description().trim().toLowerCase() + "%"));
            if (filter.tags() != null && !filter.tags().isEmpty()) {
                predicates.add(root.join("tags").in(filter.tags()));
                query.distinct(true);
            }
            return cb.and(predicates.toArray(Predicate[]::new));
        };
    }
}
