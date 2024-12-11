package traceUtils

// StackTraceConfig allows configuring the detail level of the printed stack trace.
type StackTraceConfig struct {
	SkipFrames        int
	IncludeSourceCode bool
	IncludePC         bool
	ShortFuncNames    bool
	ShowFullPath      bool
	ShowLineNumbers   bool
	FrameSeparator    string
	ChunkSeparator    string
	ChunkIndentation  string
}

// NewStackTrace - returns a nicely formatted stack trace according to cfg, default is full verbose stack trace.
// modified from https://github.com/gin-gonic/gin/blob/master/recovery.go#L111-L169
func NewStackTrace(opts ...StackTraceOption) []byte {
	cfg := StackTraceConfig{
		SkipFrames:        0,
		IncludeSourceCode: true,
		IncludePC:         true,
		ShortFuncNames:    true,
		ShowFullPath:      true,
		ShowLineNumbers:   true,
		FrameSeparator:    "\n",
		ChunkSeparator:    "\n",
		ChunkIndentation:  "\t",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	var frames []string
	var lines [][]byte
	var lastFile string

	for i := cfg.SkipFrames; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Determine what file/line info to show
		var displayFile string
		if cfg.ShowFullPath {
			displayFile = file
		} else {
			displayFile = filepath.Base(file)
		}

		var frameHeader string
		if cfg.ShowLineNumbers {
			if cfg.IncludePC {
				frameHeader = fmt.Sprintf("%s:%d (0x%x)", displayFile, line, pc)
			} else {
				frameHeader = fmt.Sprintf("%s:%d", displayFile, line)
			}
		} else {
			if cfg.IncludePC {
				frameHeader = fmt.Sprintf("%s (0x%x)", displayFile, pc)
			} else {
				frameHeader = displayFile
			}
		}

		funcName := resolveFuncName(pc, cfg.ShortFuncNames)

		var frameChunks []string
		frameChunks = append(frameChunks, frameHeader)

		if cfg.IncludeSourceCode {
			if file != lastFile {
				data, err := os.ReadFile(file)
				if err == nil {
					lines = bytes.Split(data, []byte{'\n'})
					lastFile = file
				} else {
					lines = nil
				}
			}
			code := source(lines, line)
			frameChunks = append(frameChunks, fmt.Sprintf("%s%s: %s", cfg.ChunkIndentation, funcName, code))
		} else {
			frameChunks = append(frameChunks, fmt.Sprintf("%s%s", cfg.ChunkIndentation, funcName))
		}

		frames = append(frames, strings.Join(frameChunks, cfg.ChunkSeparator))
	}

	// Join all frames with the configured frameSeparator
	output := strings.Join(frames, cfg.FrameSeparator)
	return []byte(output)
}

// resolveFuncName returns the function name based on the config.
func resolveFuncName(pc uintptr, shortNames bool) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return unknown
	}

	if shortNames {
		name := []byte(fn.Name())
		if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
			name = name[lastSlash+1:]
		}
		name = bytes.ReplaceAll(name, centerDot, dot)
		if period := bytes.Index(name, dot); period >= 0 {
			name = name[period+1:]
		}
		return name
	}

	return []byte(fn.Name())
}

// source returns a space-trimmed slice of the nth line.
func source(lines [][]byte, n int) []byte {
	n-- // stack traces are 1-indexed
	if n < 0 || n >= len(lines) {
		return unknown
	}
	return bytes.TrimSpace(lines[n])
}

var (
	slash     = []byte("/")
	dot       = []byte(".")
	centerDot = []byte("·")
	unknown   = []byte("???")
)

type StackTraceOption func(*StackTraceConfig)

func WithSkipFrames(skip int) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.SkipFrames = skip
	}
}

func WithIncludeSourceCode(include bool) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.IncludeSourceCode = include
	}
}

func WithIncludePC(include bool) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.IncludePC = include
	}
}

func WithShortFuncNames(short bool) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.ShortFuncNames = short
	}
}

func WithShowFullPath(full bool) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.ShowFullPath = full
	}
}

func WithShowLineNumbers(show bool) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.ShowLineNumbers = show
	}
}

func WithFrameSeparator(frameSeparator string) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.FrameSeparator = frameSeparator
	}
}

func WithChunkSeparator(chunkSeparator string) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.ChunkSeparator = chunkSeparator
	}
}

func WithChunkIndentation(chunkIndentation string) StackTraceOption {
	return func(cfg *StackTraceConfig) {
		cfg.ChunkIndentation = chunkIndentation
	}
}
