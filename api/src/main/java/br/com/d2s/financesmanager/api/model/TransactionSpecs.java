package br.com.d2s.financesmanager.api.model;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
import jakarta.persistence.criteria.Expression;
import jakarta.persistence.criteria.JoinType;
import org.springframework.data.jpa.domain.Specification;

import java.time.Instant;
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
            if (filter.year() != null)
                predicates.add(cb.equal(cb.function("date_part", Integer.class, cb.literal("year"), paymentDate), filter.year()));
            if (filter.month() != null)
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
