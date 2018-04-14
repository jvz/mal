package mal;

import java.util.List;
import java.util.stream.Collectors;

public class Types {
    public interface MalType {
        Object value();

        default MalType meta() {
            return null;
        }

        default String print(boolean pretty) {
            Object value = value();
            return value == null ? "nil" : value.toString();
        }
    }

    public enum MalNil implements MalType {
        NIL;

        @Override
        public Void value() {
            return null;
        }
    }

    public enum MalBool implements MalType {
        FALSE(false),
        TRUE(true);

        private final boolean value;

        MalBool(boolean value) {
            this.value = value;
        }

        @Override
        public Boolean value() {
            return value;
        }
    }

    public static class MalInt implements MalType {
        private final int value;

        public MalInt(int value) {
            this.value = value;
        }

        @Override
        public Integer value() {
            return value;
        }
    }

    private static abstract class AbstractStringType implements MalType {
        final String value;

        private AbstractStringType(String value) {
            this.value = value;
        }

        @Override
        public String value() {
            return value;
        }
    }

    public static class MalSymbol extends AbstractStringType {
        public MalSymbol(String value) {
            super(value);
        }
    }

    public static class MalString extends AbstractStringType {
        public MalString(String value) {
            super(value);
        }
    }

    public static class MalKeyword extends AbstractStringType {
        public MalKeyword(String value) {
            super(value);
        }

        @Override
        public String print(boolean pretty) {
            return ':' + value;
        }
    }

    public static class MalList implements MalType {
        private final List<MalType> value;

        public MalList(List<MalType> value) {
            this.value = value;
        }

        @Override
        public Object value() {
            return value;
        }

        @Override
        public String print(boolean pretty) {
            return value.stream().map(v -> v.print(pretty)).collect(Collectors.joining(" ", "(", ")"));
        }
    }
}
