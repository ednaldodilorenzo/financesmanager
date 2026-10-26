package br.com.d2s.financesmanager.api.config.parser.csv;

import br.com.d2s.financesmanager.api.infra.publicid.PublicId;
import br.com.d2s.financesmanager.api.dto.TransactionDto;
import br.com.d2s.financesmanager.api.model.Transaction;
import com.opencsv.CSVParser;
import com.opencsv.CSVParserBuilder;
import com.opencsv.CSVReader;
import com.opencsv.CSVReaderBuilder;
import com.opencsv.exceptions.CsvException;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Component;

import java.io.*;
import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.function.Function;

@Component
@AllArgsConstructor
public class ConfigurableCsvBankParser implements ExtractParser{
    private final CsvProperties csvProperties;
    @Override
    public boolean supports(String bank) {
        return csvProperties.layouts().containsKey(bank);
    }

    private record ParsedRow(String description, TransactionGenericIdentity identity) {}

    @Override
    public List<TransactionDto.UploadResponse> parse(InputStream inputStream,
                                                     String bank,
                                                     Long accountId,
                                                     Long userId,
                                                     Function<List<TransactionGenericIdentity>, Set<TransactionGenericIdentity>> findDuplicates,
                                                     Instant fixedPaymentDate) throws IOException {
        var layout = csvProperties.layouts().get(bank);
        try(Reader reader = new BufferedReader(new InputStreamReader(inputStream, layout.charset()))) {
            CSVParser parser = new CSVParserBuilder()
                    .withSeparator(layout.delimiter().charAt(0))
                    .build();
            CSVReader csvReader = new CSVReaderBuilder(reader)
                    .withCSVParser(parser)
                    .build();
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern(layout.datePattern());

            List<String[]> rows = csvReader.readAll();
            List<ParsedRow> parsedRows = rows.subList(layout.headerLines(), rows.size()- layout.footerLines())
                    .stream()
                    .filter(s -> s.length > 0)
                    .map(row -> {
                        Instant paymentDate = fixedPaymentDate != null ? fixedPaymentDate :
                                LocalDate.parse(row[layout.paymentDateColumn()], formatter).atStartOfDay(ZoneOffset.UTC).toInstant();
                        Instant transactionDate =
                                LocalDate.parse(row[layout.paymentDateColumn()], formatter).atStartOfDay(ZoneOffset.UTC).toInstant();
                        var value = Long.parseLong(row[layout.ammountColumn()].replace(".", "").replace(",", ""));
                        value = layout.invertSign() ? -value : value;

                        return new ParsedRow(row[layout.descriptionColumn()],
                                new TransactionGenericIdentity(value, paymentDate, transactionDate, userId));
                    }).toList();

            Set<TransactionGenericIdentity> duplicates =
                    findDuplicates.apply(parsedRows.stream().map(ParsedRow::identity).toList());

            return parsedRows.stream().map(p -> new TransactionDto.UploadResponse(
                null,
                PublicId.of(accountId),
                p.description(),
                p.identity().value(),
                p.identity().paymentDate(),
                p.identity().transactionDate(),
                duplicates.contains(p.identity()),
                ""
            )).toList();
        } catch (CsvException e) {
            throw new IOException(e);
        }
    }
}
