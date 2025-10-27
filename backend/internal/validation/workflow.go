package validation

import (
	"fmt"
	"regexp"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
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
		return errors.NewValidationError("Invalid trigger type").
			WithField("trigger.type").
			WithDetails(fmt.Sprintf("Type '%s' is not supported", trigger.Type))
	}

	if trigger.Type == "on_row_created" || trigger.Type == "on_row_updated" || trigger.Type == "on_row_deleted" {
		if trigger.TableName == "" {
			return errors.NewValidationError("Table name is required").
				WithField("trigger.table_name").
				WithDetails(fmt.Sprintf("Required for trigger type '%s'", trigger.Type))
		}
	}

	if trigger.Type == "scheduled" {
		if trigger.Schedule == "" {
			return errors.NewValidationError("Schedule is required").
				WithField("trigger.schedule").
				WithDetails("Required for scheduled triggers")
		}
		if err := validateCronExpression(trigger.Schedule); err != nil {
			return errors.NewValidationError("Invalid cron expression").
				WithField("trigger.schedule").
				WithDetails(err.Error())
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
		return errors.NewValidationError("At least one action is required").WithField("actions")
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
			return errors.NewValidationError("Action ID is required").
				WithField(fmt.Sprintf("actions[%d].id", i))
		}

		if !validActionTypes[action.Type] {
			return errors.NewValidationError("Invalid action type").
				WithField(fmt.Sprintf("actions[%d].type", i)).
				WithDetails(fmt.Sprintf("Type '%s' is not supported", action.Type))
		}

		if len(action.Config) == 0 {
			return errors.NewValidationError("Action config is required").
				WithField(fmt.Sprintf("actions[%d].config", i))
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
			return errors.NewValidationError("Webhook URL is required").
				WithField(fmt.Sprintf("actions[%d].config.url", index))
		}
	case "send_email":
		if _, ok := action.Config["to"].(string); !ok {
			return errors.NewValidationError("Email recipient is required").
				WithField(fmt.Sprintf("actions[%d].config.to", index))
		}
	case "update_row", "create_row", "delete_row":
		if _, ok := action.Config["table"].(string); !ok {
			return errors.NewValidationError("Table name is required").
				WithField(fmt.Sprintf("actions[%d].config.table", index))
		}
	}
	return nil
}
