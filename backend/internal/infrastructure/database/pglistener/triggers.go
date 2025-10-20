package pglistener

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// TriggerManager handles creation and management of database triggers
type TriggerManager struct {
	db *sql.DB
}

// NewTriggerManager creates a new trigger manager
func NewTriggerManager(db *sql.DB) *TriggerManager {
	return &TriggerManager{db: db}
}

// CreateNotifyTrigger creates a trigger function and trigger for table changes
func (tm *TriggerManager) CreateNotifyTrigger(ctx context.Context, schema, table, channel string) error {
	// Create the trigger function if it doesn't exist
	err := tm.createTriggerFunction(ctx)
	if err != nil {
		return fmt.Errorf("failed to create trigger function: %w", err)
	}

	// Create the trigger on the specific table
	err = tm.createTableTrigger(ctx, schema, table, channel)
	if err != nil {
		return fmt.Errorf("failed to create table trigger: %w", err)
	}

	return nil
}

// createTriggerFunction creates the generic notify trigger function
func (tm *TriggerManager) createTriggerFunction(ctx context.Context) error {
	query := `
CREATE OR REPLACE FUNCTION notify_table_change()
RETURNS TRIGGER AS $$
DECLARE
    notification json;
    old_data json := NULL;
    new_data json := NULL;
BEGIN
    -- Capture old and new data based on operation
    IF (TG_OP = 'DELETE') THEN
        old_data = row_to_json(OLD);
    ELSIF (TG_OP = 'UPDATE') THEN
        old_data = row_to_json(OLD);
        new_data = row_to_json(NEW);
    ELSIF (TG_OP = 'INSERT') THEN
        new_data = row_to_json(NEW);
    END IF;

    -- Build notification payload
    notification = json_build_object(
        'type', TG_OP,
        'table', TG_TABLE_NAME,
        'schema', TG_TABLE_SCHEMA,
        'timestamp', NOW(),
        'old_data', old_data,
        'new_data', new_data
    );

    -- Send notification
    PERFORM pg_notify(TG_ARGV[0], notification::text);

    -- Return appropriate row
    IF (TG_OP = 'DELETE') THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
`

	_, err := tm.db.ExecContext(ctx, query)
	return err
}

// createTableTrigger creates a trigger on a specific table
func (tm *TriggerManager) createTableTrigger(ctx context.Context, schema, table, channel string) error {
	triggerName := fmt.Sprintf("%s_%s_notify", schema, table)

	query := fmt.Sprintf(`
DROP TRIGGER IF EXISTS %s ON %s.%s;
CREATE TRIGGER %s
AFTER INSERT OR UPDATE OR DELETE ON %s.%s
FOR EACH ROW
EXECUTE FUNCTION notify_table_change('%s');
`,
		triggerName, schema, table,
		triggerName,
		schema, table,
		channel,
	)

	_, err := tm.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create trigger %s: %w", triggerName, err)
	}

	return nil
}

// DropTrigger removes a trigger from a table
func (tm *TriggerManager) DropTrigger(ctx context.Context, schema, table string) error {
	triggerName := fmt.Sprintf("%s_%s_notify", schema, table)

	query := fmt.Sprintf(`
DROP TRIGGER IF EXISTS %s ON %s.%s;
`,
		triggerName, schema, table,
	)

	_, err := tm.db.ExecContext(ctx, query)
	return err
}

// EnableTriggersForTable enables realtime for a specific table
func (tm *TriggerManager) EnableTriggersForTable(ctx context.Context, schema, table string, channels ...string) error {
	channel := fmt.Sprintf("table_changes:%s.%s", schema, table)
	if len(channels) > 0 {
		channel = channels[0]
	}

	return tm.CreateNotifyTrigger(ctx, schema, table, channel)
}

// DisableTriggersForTable disables realtime for a specific table
func (tm *TriggerManager) DisableTriggersForTable(ctx context.Context, schema, table string) error {
	return tm.DropTrigger(ctx, schema, table)
}

// ListTriggers returns all notify triggers in the database
func (tm *TriggerManager) ListTriggers(ctx context.Context) ([]TriggerInfo, error) {
	query := `
SELECT
    t.trigger_schema,
    t.event_object_table,
    t.trigger_name,
    t.action_timing,
    t.event_manipulation
FROM information_schema.triggers t
WHERE t.action_statement LIKE '%notify_table_change%'
ORDER BY t.trigger_schema, t.event_object_table;
`

	rows, err := tm.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var triggers []TriggerInfo
	for rows.Next() {
		var ti TriggerInfo
		err := rows.Scan(
			&ti.Schema,
			&ti.Table,
			&ti.TriggerName,
			&ti.Timing,
			&ti.Event,
		)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, ti)
	}

	return triggers, rows.Err()
}

// TriggerInfo holds information about a database trigger
type TriggerInfo struct {
	Schema      string
	Table       string
	TriggerName string
	Timing      string
	Event       string
}

// String formats trigger info as a string
func (ti TriggerInfo) String() string {
	return fmt.Sprintf("%s.%s: %s %s %s",
		ti.Schema, ti.Table, ti.Timing, ti.Event, ti.TriggerName)
}

// CreateRealtimeChannel creates a standardized channel name for table changes
func CreateRealtimeChannel(schema, table string) string {
	return fmt.Sprintf("realtime:%s:%s", schema, table)
}

// ParseRealtimeChannel parses a realtime channel name into schema and table
func ParseRealtimeChannel(channel string) (schema, table string, ok bool) {
	parts := strings.Split(channel, ":")
	if len(parts) != 3 || parts[0] != "realtime" {
		return "", "", false
	}
	return parts[1], parts[2], true
}
