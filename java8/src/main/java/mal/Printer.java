package mal;

import mal.Types.MalType;

public class Printer {

    public static String pr_str(MalType value) {
        return value.print(true);
    }

}
