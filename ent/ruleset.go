// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"sifu-box/ent/ruleset"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// RuleSet is the model entity for the RuleSet schema.
type RuleSet struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Tag holds the value of the "tag" field.
	Tag string `json:"tag,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Path holds the value of the "path" field.
	Path string `json:"path,omitempty"`
	// Format holds the value of the "format" field.
	Format string `json:"format,omitempty"`
	// Label holds the value of the "label" field.
	Label string `json:"label,omitempty"`
	// DownloadDetour holds the value of the "download_detour" field.
	DownloadDetour string `json:"download_detour,omitempty"`
	// UpdateInterval holds the value of the "update_interval" field.
	UpdateInterval string `json:"update_interval,omitempty"`
	// NameServer holds the value of the "name_server" field.
	NameServer string `json:"name_server,omitempty"`
	// China holds the value of the "china" field.
	China        bool `json:"china,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*RuleSet) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case ruleset.FieldChina:
			values[i] = new(sql.NullBool)
		case ruleset.FieldID:
			values[i] = new(sql.NullInt64)
		case ruleset.FieldTag, ruleset.FieldType, ruleset.FieldPath, ruleset.FieldFormat, ruleset.FieldLabel, ruleset.FieldDownloadDetour, ruleset.FieldUpdateInterval, ruleset.FieldNameServer:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the RuleSet fields.
func (rs *RuleSet) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case ruleset.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			rs.ID = int(value.Int64)
		case ruleset.FieldTag:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tag", values[i])
			} else if value.Valid {
				rs.Tag = value.String
			}
		case ruleset.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				rs.Type = value.String
			}
		case ruleset.FieldPath:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field path", values[i])
			} else if value.Valid {
				rs.Path = value.String
			}
		case ruleset.FieldFormat:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field format", values[i])
			} else if value.Valid {
				rs.Format = value.String
			}
		case ruleset.FieldLabel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field label", values[i])
			} else if value.Valid {
				rs.Label = value.String
			}
		case ruleset.FieldDownloadDetour:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field download_detour", values[i])
			} else if value.Valid {
				rs.DownloadDetour = value.String
			}
		case ruleset.FieldUpdateInterval:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field update_interval", values[i])
			} else if value.Valid {
				rs.UpdateInterval = value.String
			}
		case ruleset.FieldNameServer:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name_server", values[i])
			} else if value.Valid {
				rs.NameServer = value.String
			}
		case ruleset.FieldChina:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field china", values[i])
			} else if value.Valid {
				rs.China = value.Bool
			}
		default:
			rs.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the RuleSet.
// This includes values selected through modifiers, order, etc.
func (rs *RuleSet) Value(name string) (ent.Value, error) {
	return rs.selectValues.Get(name)
}

// Update returns a builder for updating this RuleSet.
// Note that you need to call RuleSet.Unwrap() before calling this method if this RuleSet
// was returned from a transaction, and the transaction was committed or rolled back.
func (rs *RuleSet) Update() *RuleSetUpdateOne {
	return NewRuleSetClient(rs.config).UpdateOne(rs)
}

// Unwrap unwraps the RuleSet entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rs *RuleSet) Unwrap() *RuleSet {
	_tx, ok := rs.config.driver.(*txDriver)
	if !ok {
		panic("ent: RuleSet is not a transactional entity")
	}
	rs.config.driver = _tx.drv
	return rs
}

// String implements the fmt.Stringer.
func (rs *RuleSet) String() string {
	var builder strings.Builder
	builder.WriteString("RuleSet(")
	builder.WriteString(fmt.Sprintf("id=%v, ", rs.ID))
	builder.WriteString("tag=")
	builder.WriteString(rs.Tag)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(rs.Type)
	builder.WriteString(", ")
	builder.WriteString("path=")
	builder.WriteString(rs.Path)
	builder.WriteString(", ")
	builder.WriteString("format=")
	builder.WriteString(rs.Format)
	builder.WriteString(", ")
	builder.WriteString("label=")
	builder.WriteString(rs.Label)
	builder.WriteString(", ")
	builder.WriteString("download_detour=")
	builder.WriteString(rs.DownloadDetour)
	builder.WriteString(", ")
	builder.WriteString("update_interval=")
	builder.WriteString(rs.UpdateInterval)
	builder.WriteString(", ")
	builder.WriteString("name_server=")
	builder.WriteString(rs.NameServer)
	builder.WriteString(", ")
	builder.WriteString("china=")
	builder.WriteString(fmt.Sprintf("%v", rs.China))
	builder.WriteByte(')')
	return builder.String()
}

// RuleSets is a parsable slice of RuleSet.
type RuleSets []*RuleSet
