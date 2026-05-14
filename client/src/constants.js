// Starter snippets shown in the editor for each language. They demonstrate
// the stdin-reading pattern used by the judge — copy-paste-ready.

export const LANGUAGE_VERSIONS = {
  python: '3.x',
  javascript: 'Node 20',
  cpp: 'g++ -std=c++17',
  java: 'OpenJDK 17',
};

export const CODE_SNIPPETS = {
  python: `# Read whitespace-separated tokens from stdin.
import sys
data = sys.stdin.read().split()

# example: print the first token
if data:
    print(data[0])
`,

  javascript: `// Read all of stdin, then split.
let data = '';
process.stdin.on('data', d => data += d);
process.stdin.on('end', () => {
    const tokens = data.split(/\\s+/).filter(Boolean);
    // example: print the first token
    if (tokens.length) console.log(tokens[0]);
});
`,

  cpp: `#include <bits/stdc++.h>
using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);

    // example: read one int and print it
    int x;
    if (cin >> x) cout << x << "\\n";
    return 0;
}
`,

  // IMPORTANT: the class MUST be named Main. The judge saves your code as
  // Main.java and runs \`java Main\`.
  java: `import java.util.*;
import java.io.*;

public class Main {
    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        StreamTokenizer st = new StreamTokenizer(br);

        // example: read one int and print it
        if (st.nextToken() != StreamTokenizer.TT_EOF) {
            System.out.println((long) st.nval);
        }
    }
}
`,
};
