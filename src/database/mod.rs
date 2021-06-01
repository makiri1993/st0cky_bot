use crate::models::keyword::Keyword;
use sqlx::{postgres::PgQueryResult, Pool, Postgres};

pub async fn create_user(
    connection: &Pool<Postgres>,
    id: i64,
    name: &Option<String>,
    automatic_sending: Option<bool>,
) -> Result<PgQueryResult, sqlx::Error> {
    sqlx::query!(
        r#"
INSERT INTO users ( id, name, automatic_sending )
VALUES ( $1, $2, $3 )
ON CONFLICT (id) DO NOTHING
        "#,
        id,
        name.as_ref().unwrap(),
        automatic_sending
    )
    .execute(connection)
    .await
}

pub async fn find_keywords_per_user(
    connection: &Pool<Postgres>,
    user_id: i64,
) -> Result<Vec<Keyword>, sqlx::Error> {
    sqlx::query_as!(
        Keyword,
        r#"
        SELECT id,searchterm, user_id
        FROM keywords 
        WHERE user_id = $1
        "#,
        user_id
    )
    .fetch_all(connection)
    .await
}

pub async fn create_keyword(
    connection: &Pool<Postgres>,
    searchterm: &str,
    user_id: i64,
) -> Result<i32, sqlx::Error> {
    Ok(sqlx::query!(
        r#"
            INSERT INTO keywords ( searchterm, user_id )
            VALUES ( $1, $2 ) 
            RETURNING id
            "#,
        searchterm,
        user_id
    )
    .fetch_one(connection)
    .await?
    .id)
}
