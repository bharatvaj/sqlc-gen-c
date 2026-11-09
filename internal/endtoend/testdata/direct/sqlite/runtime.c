#include "gen/queries.h"
#include <sqlite3.h>
#include <string.h>

#define CHECK(x) do { if (!(x)) return __LINE__; } while (0)

int main(void)
{
    sqlite3 *db = NULL;
    int changed = 0;
    sqlc_Get_row row = {0};
    CHECK(sqlite3_open(":memory:", &db) == SQLITE_OK);
    CHECK(sqlite3_exec(db, "CREATE TABLE records(id TEXT PRIMARY KEY NOT NULL, value TEXT NOT NULL, note TEXT, state INTEGER NOT NULL)", NULL, NULL, NULL) == SQLITE_OK);
    CHECK(sqlc_Insert(db, "key", "value", &changed) == SQLITE_OK && changed == 1);
    CHECK(sqlc_Get(db, "missing", &row) == SQLITE_NOTFOUND);
    CHECK(sqlc_Get(db, "key", &row) == SQLITE_OK);
    CHECK(!strcmp(row.value, "value") && row.note == NULL && row.state == 7);
    CHECK(sqlite3_next_stmt(db, NULL) == NULL);
    sqlc_Get_free(&row);
    CHECK(sqlc_Insert(db, "key", "other", &changed) == SQLITE_CONSTRAINT);
    CHECK(sqlite3_next_stmt(db, NULL) == NULL);
    CHECK(sqlite3_close(db) == SQLITE_OK);
    return 0;
}
