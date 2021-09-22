package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DatabaseConditionPending  string = "Pending"
	DatabaseConditionCreating string = "Creating"
	DatabaseConditionCreated  string = "Created"
	DatabaseConditionSynced   string = "Synced"
	DatabaseConditionError    string = "Errored"
	DatabaseConditionUpdating string = "Updating"
	DatabaseConditionUpdated  string = "Updated"
)

const (
	DatabaseConditionReasonPending  string = "PendingDatabase"
	DatabaseConditionReasonCreating string = "CreatingDatabase"
	DatabaseConditionReasonCreated  string = "CreatedDatabase"
	DatabaseConditionReasonSynced   string = "SyncedDatabase"
	DatabaseConditionReasonError    string = "ErroredDatabase"
	DatabaseConditionReasonUpdating string = "UpdatingDatabase"
	DatabaseConditionReasonUpdated  string = "UpdatedDatabase"
)

func (d *Database) PendingCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionPending, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonPending, Message: "Database is pending"}
}

func (d *Database) CreatingCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionCreating, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonCreating, Message: "Database is creating"}
}

func (d *Database) CreatedCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionCreated, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonCreated, Message: "Database successfully created"}
}

func (d *Database) SyncedCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionSynced, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonSynced, Message: "Database successfully synced"}
}

func (d *Database) ErroredCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionError, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonError, Message: "Database is erroring"}
}

func (d *Database) UpdatingCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionUpdating, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonUpdating, Message: "Database is updating"}
}

func (d *Database) UpdatedCondition() *metav1.Condition {
	return &metav1.Condition{Type: DatabaseConditionUpdated, Status: metav1.ConditionTrue,
		Reason: DatabaseConditionReasonUpdated, Message: "Database successfully updated"}
}
