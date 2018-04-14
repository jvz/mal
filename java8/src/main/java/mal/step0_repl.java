package mal;

import java.io.IOException;

public class step0_repl implements Step<String, String> {

    public static void main(String[] args) throws IOException {
        step0_repl instance = new step0_repl();
        String line;
        while ((line = IO.prompt("user> ")) != null) {
            System.out.println(instance.rep(line, ""));
        }
    }

    @Override
    public String read(String str) {
        return str;
    }

    @Override
    public String eval(String ast, String env) {
        return ast;
    }

    @Override
    public String print(String ast) {
        return ast;
    }
}
