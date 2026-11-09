
#include "sqlite3.h"
#include "models.h"




int sql_GetAuthor(sql_context *ctx, int64_t id, Author *result) {

    const char *sql = "-- name: GetAuthor :one\n"
                      "SELECT id, name FROM authors WHERE id = $1\n";

    sqlite3_stmt *stmt;
    int rc = sqlite3_prepare_v2(ctx->db, sql, -1, &stmt, 0);
    if (rc != SQLITE_OK) return rc;

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) return rc;

    rc = sqlite3_finalize(stmt);
    if (rc != SQLITE_OK) return rc;

	return SQLITE_OK;
}
