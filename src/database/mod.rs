use crate::models::api::ContextualNews;
use crate::models::keyword::Keyword;
use crate::models::user::User;
use sqlx::types::BigDecimal;
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

pub async fn find_all_users(connection: &Pool<Postgres>) -> Result<Vec<User>, sqlx::Error> {
    sqlx::query_as!(
        User,
        r#"
        SELECT * FROM users
        "#
    )
    .fetch_all(connection)
    .await
}

pub async fn find_keywords_per_user(
    connection: &Pool<Postgres>,
    user_id: &i64,
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

pub async fn create_news(
    connection: &Pool<Postgres>,
    news: ContextualNews,
    user_id: i64,
) -> Result<BigDecimal, sqlx::Error> {
    Ok(sqlx::query!(
        r#"
        INSERT INTO news (id, title, url, description, date_published, user_id )
        VALUES ( $1, $2, $3, $4, $5, $6 )
        ON CONFLICT (id) DO NOTHING 
        RETURNING id
        "#,
        news.id,
        news.title,
        news.url,
        news.description,
        news.date_published,
        user_id
    )
    .fetch_one(connection)
    .await?
    .id)
}
