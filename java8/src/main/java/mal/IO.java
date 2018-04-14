package mal;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

public class IO {
    private static final BufferedReader stdin = new BufferedReader(new InputStreamReader(System.in));

    public static String prompt(String prompt) throws IOException {
        System.out.print(prompt);
        return stdin.readLine();
    }

    public static String readLine() throws IOException {
        return stdin.readLine();
    }
}
