package store

// CommitMetadata stores the prompt and resulting diff for a commit
type CommitMetadata struct {
	Commit    string   `json:"commit"`
	Prompt    string   `json:"prompt"`
	Details   string   `json:"details,omitempty"`
	Timestamp string   `json:"timestamp"`
	Diff      string   `json:"diff,omitempty"`
	Files     []string `json:"files,omitempty"`
	Model     string   `json:"model,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// Store manages vibe diff metadata
type Store struct {
	dataPath string
}

// NewStore creates a new store with the given data path
func NewStore(dataPath string) *Store {
	return &Store{dataPath: dataPath}
}
