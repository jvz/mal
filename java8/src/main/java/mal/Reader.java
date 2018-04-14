package mal;

import mal.Types.MalInt;
import mal.Types.MalList;
import mal.Types.MalSymbol;
import mal.Types.MalType;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Reader {

    private final List<String> tokens;
    private int pos;

    private Reader(List<String> tokens) {
        this.tokens = tokens;
    }

    private String peek() {
        return pos < tokens.size() ? tokens.get(pos) : null;
    }

    private String next() {
        return tokens.get(pos++);
    }

    private MalType readForm() {
        String token = peek();
        if (token == null) {
            throw new MissingTokenException("end of expression");
        }
        switch (token.charAt(0)) {
            case '(':
                return readList();
            case ')':
                throw new UnexpectedTokenException(token);
            default:
                return readAtom();
        }
    }

    private MalType readList() {
        next(); // (
        List<MalType> acc = new ArrayList<>();
        while (peek() != null) {
            if (")".equals(peek())) break;
            acc.add(readForm());
        }
        next(); // )
        return new MalList(acc);
    }

    private static final Pattern INTEGRAL = Pattern.compile("-?\\d+");

    private MalType readAtom() {
        String token = next();
        if (INTEGRAL.matcher(token).matches()) {
            return new MalInt(Integer.parseInt(token));
        }
        return new MalSymbol(token);
    }

    public static MalType read_str(String str) {
        Reader reader = tokenize(str);
        return reader.readForm();
    }

    private static final Pattern TOKENIZER =
            Pattern.compile("[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\\\.|[^\\\\\"])*\"|;.*|[^\\s\\[\\]{}('\"`,;)]*)");

    private static Reader tokenize(String str) {
        List<String> tokens = new ArrayList<>();
        Matcher matcher = TOKENIZER.matcher(str);
        while (matcher.find()) {
            tokens.add(matcher.group(1));
        }
        return new Reader(tokens);
    }

    public static class UnexpectedTokenException extends IllegalArgumentException {
        public UnexpectedTokenException(String token) {
            super("Unexpected token: " + token);
        }
    }

    public static class MissingTokenException extends IllegalArgumentException {
        public MissingTokenException(String token) {
            super("Missing token: " + token);
        }
    }

}
