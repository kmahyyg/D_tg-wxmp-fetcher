package db

import (
	"context"
	"database/sql"
	"math/rand"
	"net/http"

	"bitbucket.org/mutongx/go-utils/log"
	_ "github.com/lib/pq" // Load the PostgreSQL driver

	"bitbucket.org/mutze5/wxfetcher/article"
	"bitbucket.org/mutze5/wxfetcher/proto"
)

const keyChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenKeyChars = int64(len(keyChars))

var keyLen = 8

var db *sql.DB

var stmts = map[string]string{
	"GetArticleMeta": "" +
		"SELECT source, account_id, message_id, article_index, signature, title, author, brief, image_url FROM article, key WHERE key.key = $1 AND article.uuid = key.uuid",
	"FindArticleByIdentifier": "" +
		"SELECT article.uuid, key.key FROM article, key WHERE account_id = $1 AND message_id = $2 AND article_index = $3 AND article.uuid = key.uuid",
	"InsertWxArticle": "" +
		"INSERT INTO article " +
		"(uuid, source, account_id, message_id, article_index, signature, title, author, brief, timestamp, image_url, body) " +
		"VALUES " +
		"(gen_random_uuid(), 'wechat', $1, $2, $3, $4, $5, $6, $7, to_timestamp($8), $9, $10) " +
		"RETURNING uuid",
	"UpdateWxArticle": "" +
		"UPDATE article " +
		"SET source = 'wechat', signature = $4, title = $5, author = $6, brief = $7, timestamp = to_timestamp($8), image_url = $9, body = $10 " +
		"WHERE account_id = $1 AND message_id = $2 AND article_index = $3",
	"CreateLinkKey": "" +
		"INSERT INTO key" +
		"(key, uuid)" +
		"VALUES " +
		"($1, $2)",
}
var preparedStmts = map[string]*sql.Stmt{}

// Connect makes a connection to current database
func Connect(ctx context.Context, driver string, source string) (err error) {
	db, err = sql.Open(driver, source)
	if err != nil {
		return
	}
	if err = db.PingContext(ctx); err != nil {
		return
	}
	for stmtName, stmtContent := range stmts {
		if preparedStmts[stmtName], err = db.PrepareContext(ctx, stmtContent); err != nil {
			return
		}
	}
	return
}

// GetArticleMeta fetches article metadata by the shortened article key
func GetArticleMeta(ctx context.Context, key string) (meta *proto.ArticleMeta, err error) {
	var source string
	var accountID, messageID, articleIdx sql.NullInt64
	var signature sql.NullString
	var title, author, brief, image sql.NullString
	if err = preparedStmts["GetArticleMeta"].QueryRowContext(ctx, key).Scan(
		&source,
		&accountID, &messageID, &articleIdx, &signature,
		&title, &author, &brief, &image,
	); err != nil {
		return
	}
	switch source {
	case "legacy":
		meta, err = updateLegacyWxArticle(ctx,
			accountID.Int64, messageID.Int64, articleIdx.Int64, signature.String,
		)
	case "wechat":
		meta = &proto.ArticleMeta{
			Link:   article.WxArticleLink(accountID.Int64, messageID.Int64, articleIdx.Int64, signature.String),
			Title:  title.String,
			Author: author.String,
			Image:  image.String,
			Brief:  brief.String,
		}
	}
	return
}

// ProcessWxArticle insert a new WeChat article into database, or return the generated key if already exists
func ProcessWxArticle(ctx context.Context, atc *article.WxArticle) (meta *proto.ArticleMeta, key string, err error) {
	// Begin a transaction
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return
	}
	defer func() {
		if err != nil {
			log.Error("ProcessWxArticle", "Error when generating key for WeChat article %d/%d/%d: %v. Rolling back.", atc.AccountID, atc.MessageID, atc.ArticleIdx, err)
			if err := tx.Rollback(); err != nil {
				log.Critical("Database", "Error during rollback: %v. Database may be inconsistent.", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Critical("Database", "Error during commit: %v. Database may be inconsistent.", err)
			}
		}
	}()
	// Generate ArticleMeta for later use
	var uuid string
	meta = atc.Meta()
	// Find the article first
	err = tx.Stmt(preparedStmts["FindArticleByIdentifier"]).QueryRowContext(ctx, atc.AccountID, atc.MessageID, atc.ArticleIdx).Scan(&uuid, &key)
	// Not found
	if err == sql.ErrNoRows {
		// Clear the error
		err = nil
		// Insert a new entry
		if err = tx.Stmt(preparedStmts["InsertWxArticle"]).QueryRowContext(ctx,
			atc.AccountID, atc.MessageID, atc.ArticleIdx, atc.Signature,
			meta.Title, meta.Author, meta.Brief, meta.Timestamp, meta.Image,
			atc.ContentHTML,
		).Scan(&uuid); err != nil {
			return
		}
		// Generate and insert key
		key = randomKey(keyLen)
		_, err = tx.Stmt(preparedStmts["CreateLinkKey"]).ExecContext(ctx, key, uuid) // Ignore error handling since this is the last statement
	}
	return
}

func updateLegacyWxArticle(ctx context.Context, accountID int64, messageID int64, articleIdx int64, signature string) (meta *proto.ArticleMeta, err error) {
	// Fetch article content and parse it
	url := article.WxArticleLink(accountID, messageID, articleIdx, signature)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	atc, err := article.NewFromWxStream(resp.Body)
	if err != nil {
		return
	}
	// Update the article
	meta = atc.Meta()
	_, err = preparedStmts["UpdateWxArticle"].ExecContext(ctx,
		atc.AccountID, atc.MessageID, atc.ArticleIdx, atc.Signature,
		meta.Title, meta.Author, meta.Brief, meta.Timestamp, meta.Image,
		atc.ContentHTML)
	return
}

func randomKey(n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = keyChars[rand.Int63n(lenKeyChars)]
	}
	return string(result)
}
