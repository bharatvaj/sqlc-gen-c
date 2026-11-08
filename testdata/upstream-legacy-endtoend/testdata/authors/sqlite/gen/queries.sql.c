
#include "sqlite3.h"
#include "models.h"



typedef struct {
    sql_text Name;
    sql_nulltext Bio;
} CreateAuthorParams;

int sql_CreateAuthor(sql_context *ctx, CreateAuthorParams arg, Author *result) {

    const char *sql = "-- name: CreateAuthor :one\n"
                      "INSERT INTO authors (\n"
                      "  name, bio\n"
                      ") VALUES (\n"
                      "  ?, ?\n"
                      ")\n"
                      "RETURNING id, name, bio\n";

    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(ctx->db, sql, -1, &stmt, 0);
    if (rc != SQLITE_OK) return rc;

    sqlite3_bind_text(stmt, 1, arg.Name.data, -1, SQLITE_STATIC);
    
    sqlite3_bind_text(stmt, 1, arg.Bio.data, -1, SQLITE_STATIC); // nullable
    
    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) return rc;

    rc = sqlite3_finalize(stmt);
    if (rc != SQLITE_OK) return rc;

	return SQLITE_OK;
}






int sql_GetAuthor(sql_context *ctx, sql_int64 id, Author *result) {

    const char *sql = "-- name: GetAuthor :one\n"
                      "SELECT id, name, bio FROM authors\n"
                      "WHERE id = ? LIMIT 1\n";

    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(ctx->db, sql, -1, &stmt, 0);
    if (rc != SQLITE_OK) return rc;

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) return rc;

    rc = sqlite3_finalize(stmt);
    if (rc != SQLITE_OK) return rc;

	return SQLITE_OK;
}




