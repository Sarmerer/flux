package validation

import (
	"fmt"
	"regexp"

	"github.com/flow/internal/domain/entities"
)

func ValidateTrigger(trigger *entities.WorkflowTrigger) error {
	validTriggerTypes := map[string]bool{
		"on_row_created": true,
		"on_row_updated": true,
		"on_row_deleted": true,
		"scheduled":      true,
		"webhook":        true,
	}

	if !validTriggerTypes[trigger.Type] {
		return fmt.Errorf("invalid trigger type: %s", trigger.Type)
	}

	if trigger.Type == "on_row_created" || trigger.Type == "on_row_updated" || trigger.Type == "on_row_deleted" {
		if trigger.TableName == "" {
			return fmt.Errorf("table_name is required for trigger type: %s", trigger.Type)
		}
	}

	if trigger.Type == "scheduled" {
		if trigger.Schedule == "" {
			return fmt.Errorf("schedule is required for scheduled triggers")
		}
		if err := validateCronExpression(trigger.Schedule); err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}
	}

	return nil
}

func validateCronExpression(expr string) error {
	cronRegex := regexp.MustCompile(`^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([0-9]|1[0-9]|2[0-3])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([1-9]|1[0-2])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|[0-6]|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)$`)

	if !cronRegex.MatchString(expr) {
		return fmt.Errorf("invalid cron expression format, expected: 'minute hour day month weekday'")
	}

	return nil
}

func ValidateActions(actions []entities.WorkflowAction) error {
	if len(actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	validActionTypes := map[string]bool{
		"send_webhook": true,
		"send_email":   true,
		"update_row":   true,
		"create_row":   true,
		"delete_row":   true,
	}

	for i, action := range actions {
		if action.ID == "" {
			return fmt.Errorf("action %d: id is required", i)
		}

		if !validActionTypes[action.Type] {
			return fmt.Errorf("action %d: invalid action type: %s", i, action.Type)
		}

		if len(action.Config) == 0 {
			return fmt.Errorf("action %d: config is required", i)
		}

		if err := validateActionConfig(i, &action); err != nil {
			return err
		}
	}

	return nil
}

func validateActionConfig(index int, action *entities.WorkflowAction) error {
	switch action.Type {
	case "send_webhook":
		if _, ok := action.Config["url"].(string); !ok {
			return fmt.Errorf("action %d: webhook url is required", index)
		}
	case "send_email":
		if _, ok := action.Config["to"].(string); !ok {
			return fmt.Errorf("action %d: email recipient (to) is required", index)
		}
	case "update_row", "create_row", "delete_row":
		if _, ok := action.Config["table"].(string); !ok {
			return fmt.Errorf("action %d: table name is required", index)
		}
	}
	return nil
}
