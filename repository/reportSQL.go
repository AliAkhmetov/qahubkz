package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/heroku/go-getting-started/models"

	_ "github.com/mattn/go-sqlite3"
)

type reportSQL struct {
	db *sql.DB
}

// NewReportsSQL create new database struct
func NewReportsSQL(db *sql.DB) *reportSQL {
	return &reportSQL{db: db}
}

// CreateReport
func (r *reportSQL) CreateReport(report models.Report) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (created_by, post_id, created_at, moderator_msg, status) values (?,?,?,?,?) RETURNING id`, reportsTable)
	row := r.db.QueryRow(query, report.CreatedBy, report.PostId, report.CreatedAt, report.ModeratorMsg, report.Status)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateReport
func (r *reportSQL) UpdateReport(report models.Report) error {
	query := fmt.Sprintf(`UPDATE %s SET updated_at=?, admin_msg=?, status=? where id=?`, reportsTable)
	if _, err := r.db.Exec(query, report.UpdatedAt, report.AdminMsg, report.Status, report.Id); err != nil {
		return fmt.Errorf("can't update Report: %w", err)
	}
	return nil
}

// GetReportsByUserId
func (r *reportSQL) GetReportsByUserId(userId int) ([]models.Report, error) {
	var reports []models.Report
	query := fmt.Sprintf(`SELECT * FROM %s WHERE created_by = ?`, reportsTable)

	rows, err := r.db.Query(query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no reports found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get reports: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var UpdatedAt sql.NullTime
		var AdminMsg sql.NullString

		var report models.Report
		if err = rows.Scan(
			&report.Id,
			&report.CreatedBy,
			&report.PostId,
			&report.CreatedAt,
			&UpdatedAt,
			&report.ModeratorMsg,
			&AdminMsg,
			&report.Status,
		); err != nil {
			return nil, fmt.Errorf("can't scan report: %w", err)
		}
		report.UpdatedAt = UpdatedAt.Time
		report.AdminMsg = AdminMsg.String

		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get reports: %w", err)
	}
	return reports, nil
}

// GetReportsById
func (r *reportSQL) GetReportsById(reportId int) (models.Report, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, reportsTable)
	var report models.Report //добавил обратно
	row := r.db.QueryRow(query, reportId)

	var UpdatedAt sql.NullTime
	var AdminMsg sql.NullString

	err := row.Scan(
		&report.Id,
		&report.CreatedBy,
		&report.PostId,
		&report.CreatedAt,
		&UpdatedAt,
		&report.ModeratorMsg,
		&AdminMsg,
		&report.Status,
	)
	if err != nil {
		return report, fmt.Errorf("can't scan report: %w", err)
	}
	report.UpdatedAt = UpdatedAt.Time
	report.AdminMsg = AdminMsg.String

	return report, nil
}

// GetAllReports
func (r *reportSQL) GetAllReports() ([]models.Report, error) {
	var reports []models.Report
	query := fmt.Sprintf(`SELECT * FROM %s;`, reportsTable)

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no reports found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get reports: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var UpdatedAt sql.NullTime
		var AdminMsg sql.NullString

		var report models.Report

		if err = rows.Scan(
			&report.Id,
			&report.CreatedBy,
			&report.PostId,
			&report.CreatedAt,
			&UpdatedAt,
			&report.ModeratorMsg,
			&AdminMsg,
			&report.Status,
		); err != nil {
			return nil, fmt.Errorf("can't scan report: %w", err)
		}
		report.UpdatedAt = UpdatedAt.Time
		report.AdminMsg = AdminMsg.String
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get reports: %w", err)
	}
	return reports, nil
}
