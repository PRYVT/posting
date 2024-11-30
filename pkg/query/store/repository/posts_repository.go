package repository

import (
	"database/sql"
	"time"

	models "github.com/PRYVT/posting/pkg/models/query"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PostRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *PostRepository {
	if db == nil {
		return nil
	}
	return &PostRepository{db: db}
}

func (repo *PostRepository) GetPostById(postId uuid.UUID) (*models.Post, error) {
	stmt, err := repo.db.Prepare("SELECT id, user_id, text, image_base64, change_date FROM posts WHERE id = ? ")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var changeDate string
	var post models.Post
	err = stmt.QueryRow(postId.String()).Scan(&post.Id, &post.UserId, &post.Text, &post.ImageBase64, &changeDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, changeDate)
	if err != nil {
		log.Err(err).Msg("Error while parsing time using empty changedate")
	} else {
		post.ChangeDate = parsedTime
	}

	return &post, nil
}

func (repo *PostRepository) GetAllPosts(limit, offset int) ([]models.Post, error) {
	if limit > 100 {
		limit = 100
	}
	stmt, err := repo.db.Prepare(`SELECT id,  user_id, text, image_base64, change_date FROM posts ORDER BY "order" desc LIMIT ? OFFSET ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var changeDate string
		var userId string
		var post models.Post
		if err := rows.Scan(&post.Id, &userId, &post.Text, &post.ImageBase64, &changeDate); err != nil {
			return nil, err
		}
		post.UserId = uuid.MustParse(userId)
		parsedTime, err := time.Parse(time.RFC3339Nano, changeDate)
		if err != nil {
			log.Err(err).Msg("Error while parsing time using empty changedate")
		} else {
			post.ChangeDate = parsedTime
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *PostRepository) AddOrReplacePost(post *models.Post) error {
	stmt, err := repo.db.Prepare(`
		INSERT INTO posts (id, user_id, text, image_base64, change_date)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			text = excluded.text,
			image_base64 = excluded.image_base64,
			change_date = excluded.change_date
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.UserId.String(), post.Text, post.ImageBase64, post.ChangeDate.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}
