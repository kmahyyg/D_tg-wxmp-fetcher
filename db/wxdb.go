package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"

	"bitbucket.org/mutongx/go-utils/log"
	_ "github.com/lib/pq" // Load the PostgreSQL driver

	"bitbucket.org/mutze5/wxfetcher/article"
)

const (
	wechatLinkFormat = "https://mp.weixin.qq.com/s?__biz=%s&mid=%d&idx=%d&sn=%s"
)

const keyChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
	"CreateLinkKey": "" +
		"INSERT INTO key" +
		"(key, uuid)" +
		"VALUES " +
		"($1, $2)",
}
var preparedStmts = map[string]*sql.Stmt{}

// Connect makes a connection to current database
func Connect(ctx context.Context, driver string, source string) error {
	var err error
	db, err = sql.Open(driver, source)
	if err != nil {
		return err
	}
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	for stmtName, stmtContent := range stmts {
		if preparedStmts[stmtName], err = db.PrepareContext(ctx, stmtContent); err != nil {
			return err
		}
	}
	return nil
}

// GetArticleMeta fetches article metadata by the shortened article key
func GetArticleMeta(ctx context.Context, key string) (meta *article.Metadata, err error) {
	var source, signature string
	var accountID, messageID, articleIndex int64
	meta = &article.Metadata{}
	if err = preparedStmts["GetArticleMeta"].QueryRowContext(ctx, key).Scan(
		&source,
		&accountID, &messageID, &articleIndex,
		&signature,
		&meta.Title, &meta.Author, &meta.Brief, &meta.Image,
	); err != nil {
		meta = nil
		return
	}
	switch source {
	case "wechat":
		encodedBiz := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(accountID, 10)))
		meta.Link = fmt.Sprintf(wechatLinkFormat, encodedBiz, messageID, articleIndex, signature)
	}
	return
}

// GetWxArticleKey insert a new WeChat article into database, or return the generated key if already exists
func GetWxArticleKey(ctx context.Context, atc *article.WxArticle) (key string, err error) {
	// Begin a transaction
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Critical("Database", "Error during rollback: %v. Database may be inconsistent.", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Critical("Database", "Error during commit: %v. Database may be inconsistent.", err)
			}
		}
	}()
	if key, err = findWxArticle(ctx, tx, atc); err == nil && key == "" {
		// Insert the article if not found
		key, err = insertWxArticle(ctx, tx, atc)
	}
	if err != nil {
		log.Error("GetWxArticleKey", "Error when generating key for WeChat article %d/%d/%d: %v", atc.AccountID, atc.MessageID, atc.ArticleIdx, err)
	}
	return
}

func findWxArticle(ctx context.Context, tx *sql.Tx, atc *article.WxArticle) (key string, err error) {
	var uuid string
	if err =
		tx.Stmt(preparedStmts["FindArticleByIdentifier"]).QueryRowContext(ctx,
			atc.AccountID, atc.MessageID, atc.ArticleIdx,
		).Scan(&uuid, &key); err == sql.ErrNoRows {
		// Clear the error if no rows found
		err = nil
	}
	return
}

func insertWxArticle(ctx context.Context, tx *sql.Tx, atc *article.WxArticle) (key string, err error) {
	// Generate the author name
	var author string
	if atc.AuthorName != "" {
		author = fmt.Sprintf("%s | %s", atc.AccountName, atc.AuthorName)
	} else {
		author = atc.AccountName
	}
	// Insert the article
	var uuid string
	if err =
		tx.Stmt(preparedStmts["InsertWxArticle"]).QueryRowContext(ctx,
			atc.AccountID, atc.MessageID, atc.ArticleIdx, atc.Signature,
			atc.Title, author, atc.Brief, atc.Timestamp, atc.ArticleImageURL, atc.ContentHTML,
		).Scan(&uuid); err != nil {
		return
	}
	// Generate a new URL Key
	if key, err = randomKey(keyLen); err == nil {
		_, err = tx.Stmt(preparedStmts["CreateLinkKey"]).ExecContext(ctx, key, uuid)
	}
	return
}

func randomKey(n int) (key string, err error) {
	max := big.NewInt(int64(len(keyChars)))
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		var randInt *big.Int
		if randInt, err = rand.Int(rand.Reader, max); err != nil {
			return
		}
		result[i] = keyChars[randInt.Int64()]
	}
	return string(result), nil
}
