package br.com.d2s.financesmanager.api.config.parser.csv;

import br.com.d2s.financesmanager.api.dto.TransactionDto;
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
import java.util.List;
import java.util.function.Predicate;

@Component
@AllArgsConstructor
public class ConfigurableCsvBankParser implements ExtractParser{
    private final CsvProperties csvProperties;
    @Override
    public boolean supports(String bank) {
        return csvProperties.layouts().containsKey(bank);
    }

    @Override
    public List<TransactionDto.UploadResponse> parse(InputStream inputStream, String bank, Long accountId,
                                                     Long userId, Predicate<TransactionGenericIdentity> isDuplicate,
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
            return rows.subList(layout.headerLines(), rows.size()- layout.footerLines())
                    .stream()
                    .filter(s -> s.length > 0)
                    .map(row -> {
                        Instant paymentDate = fixedPaymentDate != null ? fixedPaymentDate :
                                LocalDate.parse(row[layout.paymentDateColumn()], formatter).atStartOfDay(ZoneOffset.UTC).toInstant();
                        Instant transactionDate =
                                LocalDate.parse(row[layout.paymentDateColumn()], formatter).atStartOfDay(ZoneOffset.UTC).toInstant();
                        Long value = Long.parseLong(row[layout.ammountColumn()].replace(".", "").replace(",", ""));

                        return new TransactionDto.UploadResponse(
                            0L,
                            accountId,
                            row[layout.descriptionColumn()],
                            value,
                            paymentDate,
                            transactionDate,
                                isDuplicate.test(new TransactionGenericIdentity(value, paymentDate, transactionDate, 1L)),
                                "");
                    }).toList();
        } catch (CsvException e) {
            throw new IOException(e);
        }
    }
}
