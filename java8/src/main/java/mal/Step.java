package mal;

public interface Step<T, E> {
    T read(String str);

    T eval(T ast, E env);

    String print(T ast);

    default String rep(String str, E env) {
        return print(eval(read(str), env));
    }
}
