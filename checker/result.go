package checker

// ScanResult 是一次完整检测的“最终结论”
type ScanResult struct {
	TargetName    string
	StartURL      string
	FinalURL      string
	FinalTitle    string
	Accessible    bool
	AllowPublic   bool
	Steps         []RedirectStep
	Error         error
	Risk          bool
	StatusCode    int
	ContentLength int64
}
