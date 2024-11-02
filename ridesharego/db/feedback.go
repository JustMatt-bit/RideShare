package db

import (
	"database/sql"
	"main/core"
)

func GetFeedbacks(db *sql.DB) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbackByID(db *sql.DB, id int) (*core.Feedback, error) {
	row := db.QueryRow("SELECT * FROM user_feedback WHERE id = ?", id)
	var f core.Feedback
	if err := row.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func GetFeedbacksByUserID(db *sql.DB, userID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback WHERE owner_user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbacksByRideID(db *sql.DB, rideID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback WHERE ride_id = ?", rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbackByUserIDAndRideID(db *sql.DB, userID, rideID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT uf.id, uf.owner_user_id, uf.ride_id, uf.score, uf.message, uf.created_at FROM user_feedback uf LEFT JOIN ride r ON r.id = uf.ride_id WHERE r.owner_user_id = ? AND uf.ride_id = ?", userID, rideID)
	if err != nil {
		return nil, err
	}

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func CreateFeedback(db *sql.DB, f core.Feedback) (int64, error) {
	result, err := db.Exec("INSERT INTO user_feedback (owner_user_id, ride_id, score, message) VALUES (?, ?, ?, ?)", f.UserID, f.RideID, f.Score, f.Message)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateFeedback(db *sql.DB, id int, f core.Feedback) error {
	_, err := db.Exec("UPDATE user_feedback SET owner_user_id = ?, ride_id = ?, score = ?, message = ? WHERE id = ?",
		f.UserID, f.RideID, f.Score, f.Message, id)
	return err
}

func DeleteFeedback(db *sql.DB, feedbackID int) error {
	_, err := db.Exec("DELETE FROM user_feedback WHERE id = ?", feedbackID)
	return err
}
