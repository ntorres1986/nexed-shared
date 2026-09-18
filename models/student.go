package models

import "time"

// Student represents a school-age child linked to a PADRE user ("hijo").
// Unlike User, a Student never authenticates — no password, role, or
// permissions; it exists only to be enrolled (see Enrollment).
// InstitutionID is denormalized rather than derived through Parent, so
// every institution-scoped query stays a single WHERE clause — the same
// pattern already used by Role.
type Student struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	InstitutionID uint         `gorm:"column:institution_id;not null;index;uniqueIndex:idx_students_institution_document" json:"institution_id"`
	Institution   *Institution `gorm:"foreignKey:InstitutionID" json:"-"`
	ParentID      uint         `gorm:"column:parent_id;not null;index" json:"parent_id"`
	Parent        *User        `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	FirstName string `gorm:"column:first_name;size:100;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;size:100;not null" json:"last_name"`
	// Unique together with InstitutionID: two students in the same school
	// can't share a document id, but the same document id may reappear in
	// a different institution.
	DocumentID string `gorm:"column:document_id;size:50;not null;uniqueIndex:idx_students_institution_document" json:"document_id"`
	// GradeID is optional (a student may be registered before their grade
	// is assigned) but required before an Enrollment can be created — see
	// EnrollmentService.Register, which reads the tuition amount from it.
	GradeID *uint        `gorm:"column:grade_id;index" json:"grade_id,omitempty"`
	Grade   *SchoolGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`

	// Code is an internal school code, distinct from DocumentID (the
	// identity document number) — e.g. a student roll/enrollment number.
	Code           string        `gorm:"column:code;size:50" json:"code,omitempty"`
	DocumentTypeID *uint         `gorm:"column:document_type_id;index" json:"document_type_id,omitempty"`
	DocumentType   *DocumentType `gorm:"foreignKey:DocumentTypeID" json:"document_type,omitempty"`
	BirthDate      *time.Time    `gorm:"column:birth_date" json:"birth_date,omitempty"`
	BirthPlace     string        `gorm:"column:birth_place;size:150" json:"birth_place,omitempty"`
	Gender         string        `gorm:"column:gender;size:20" json:"gender,omitempty"`
	BloodType      string        `gorm:"column:blood_type;size:5" json:"blood_type,omitempty"`
	Phone          string        `gorm:"column:phone;size:30" json:"phone,omitempty"`
	Email          string        `gorm:"column:email;size:150" json:"email,omitempty"`
	Neighborhood   string        `gorm:"column:neighborhood;size:150" json:"neighborhood,omitempty"`
	Address        string        `gorm:"column:address;size:255" json:"address,omitempty"`
	// Stratum is the Colombian socioeconomic stratum, 1-6 — nil when unset.
	Stratum        *int   `gorm:"column:stratum" json:"stratum,omitempty"`
	AdditionalInfo string `gorm:"column:additional_info;type:text" json:"additional_info,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Student) TableName() string {
	return "students"
}

func (s Student) FullName() string {
	return s.FirstName + " " + s.LastName
}
