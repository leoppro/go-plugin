package sink

import "context"

// Sink is an abstraction for anything that a changefeed may emit into.
type Sink interface {

	// EmitRowChangedEvents sends Row Changed Event to Sink
	// EmitRowChangedEvents may write rows to downstream directly;
	EmitRowChangedEvents(ctx context.Context, rows ...*RowChangedEvent) error

	// EmitDDLEvent sends DDL Event to Sink
	// EmitDDLEvent should execute DDL to downstream synchronously
	EmitDDLEvent(ctx context.Context, ddl *DDLEvent) error

	// FlushRowChangedEvents flushes each row which of commitTs less than or equal to `resolvedTs` into downstream.
	// TiCDC guarantees that all of Event which of commitTs less than or equal to `resolvedTs` are sent to Sink through `EmitRowChangedEvents`
	FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) error

	// EmitCheckpointTs sends CheckpointTs to Sink
	// TiCDC guarantees that all Events **in the cluster** which of commitTs less than or equal `checkpointTs` are sent to downstream successfully.
	EmitCheckpointTs(ctx context.Context, ts uint64) error

	// Close closes the Sink
	Close() error
}

// RowChangedEvent represents a row changed event
type RowChangedEvent struct {
	StartTs  uint64 `json:"start-ts"`
	CommitTs uint64 `json:"commit-ts"`

	RowID int64 `json:"row-id"`

	Table *TableName `json:"table"`

	Delete bool `json:"delete"`

	// if the table of this row only has one unique index(includes primary key),
	// IndieMarkCol will be set to the name of the unique index
	IndieMarkCol string             `json:"indie-mark-col"`
	Columns      map[string]*Column `json:"columns"`
	Keys         []string           `json:"keys"`
}

// Column represents a column value in row changed event
type Column struct {
	Type        byte        `json:"t"`
	WhereHandle *bool       `json:"h,omitempty"`
	Value       interface{} `json:"v"`
}

// DDLEvent represents a DDL event
type DDLEvent struct {
	StartTs  uint64
	CommitTs uint64
	Schema   string
	Table    string
	Query    string
	Type     int
}

// TableName represents name of a table, includes table name and schema name.
type TableName struct {
	Schema    string `toml:"db-name" json:"db-name"`
	Table     string `toml:"tbl-name" json:"tbl-name"`
	Partition int64  `json:"partition"`
}
