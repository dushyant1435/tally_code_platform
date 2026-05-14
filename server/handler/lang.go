package handler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// langSpec describes how to prepare and run a particular language. Templates
// support $SRC (the path of the saved source file), $DIR (the temp directory)
// and $OUT (the compiled binary, when applicable).
type langSpec struct {
	Extension string   // file extension, e.g. ".py"
	SrcName   string   // override the source filename (used by Java)
	Compile   []string // optional compile cmd; nil for interpreted langs
	Run       []string // run cmd
}

// Supported languages. The string keys are the same values the frontend
// dropdown sends.
var languages = map[string]langSpec{
	"python": {
		Extension: ".py",
		Run:       []string{"python3", "$SRC"},
	},
	"javascript": {
		Extension: ".js",
		Run:       []string{"node", "$SRC"},
	},
	"cpp": {
		Extension: ".cpp",
		Compile:   []string{"g++", "-O2", "-std=c++17", "$SRC", "-o", "$OUT"},
		Run:       []string{"$OUT"},
	},
	"java": {
		Extension: ".java",
		SrcName:   "Main.java",
		Compile:   []string{"javac", "$SRC"},
		Run:       []string{"java", "-cp", "$DIR", "Main"},
	},
}

// Language returns the spec for the given key, falling back to python if the
// key is empty or unknown.
func Language(key string) (string, langSpec) {
	key = strings.ToLower(strings.TrimSpace(key))
	if l, ok := languages[key]; ok {
		return key, l
	}
	return "python", languages["python"]
}

// preparedExec is the output of prepareExecution: a temp working dir, the run
// args, and a cleanup function the caller must always invoke.
type preparedExec struct {
	Dir         string
	RunArgs     []string
	Cleanup     func()
	CompileLog  string // populated when compile failed
	CompileFail bool
}

// prepareExecution writes user code to a temp dir, compiles it if necessary,
// and returns the args to run plus a cleanup function. When CompileFail is
// true the caller should record a `compilation_error` verdict and skip
// execution.
func prepareExecution(ctx context.Context, langKey, code string) (*preparedExec, error) {
	_, spec := Language(langKey)

	dir, err := os.MkdirTemp("", "tally-")
	if err != nil {
		return nil, fmt.Errorf("mkdir temp: %w", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	srcName := spec.SrcName
	if srcName == "" {
		srcName = "main" + spec.Extension
	}
	srcPath := filepath.Join(dir, srcName)
	outPath := filepath.Join(dir, "program")

	if err := os.WriteFile(srcPath, []byte(code), 0644); err != nil {
		cleanup()
		return nil, fmt.Errorf("write src: %w", err)
	}

	subst := func(tpl string) string {
		s := strings.ReplaceAll(tpl, "$SRC", srcPath)
		s = strings.ReplaceAll(s, "$OUT", outPath)
		s = strings.ReplaceAll(s, "$DIR", dir)
		return s
	}
	render := func(tpls []string) []string {
		out := make([]string, len(tpls))
		for i, t := range tpls {
			out[i] = subst(t)
		}
		return out
	}

	if len(spec.Compile) > 0 {
		cmd := exec.CommandContext(ctx, spec.Compile[0], render(spec.Compile)[1:]...)
		cmd.Dir = dir
		var stderr strings.Builder
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log := strings.TrimSpace(stderr.String())
			if log == "" {
				log = err.Error()
			}
			return &preparedExec{
				Dir:         dir,
				Cleanup:     cleanup,
				CompileFail: true,
				CompileLog:  log,
			}, nil
		}
	}

	return &preparedExec{
		Dir:     dir,
		RunArgs: render(spec.Run),
		Cleanup: cleanup,
	}, nil
}
