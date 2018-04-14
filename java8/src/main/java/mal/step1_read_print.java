package mal;

import java.io.IOException;
import mal.Types.*;

public class step1_read_print implements Step<MalType, String> {

    public static void main(String[] args) throws IOException {
        step1_read_print instance = new step1_read_print();
        String line;
        while ((line = IO.prompt("user> ")) != null) {
            try {
                System.out.println(instance.rep(line, ""));
            } catch (Exception e) {
                System.err.print("Error: ");
                e.printStackTrace();
            }
        }
    }

    @Override
    public MalType read(String str) {
        return Reader.read_str(str);
    }

    @Override
    public MalType eval(MalType ast, String env) {
        return ast;
    }

    @Override
    public String print(MalType ast) {
        return Printer.pr_str(ast);
    }
}
