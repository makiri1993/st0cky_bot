use crate::database::{create_keyword, find_keywords_per_user};
use sqlx::{Pool, Postgres};

#[derive(Clone)]
pub struct Keyword {
    pub id: i32,
    pub searchterm: Option<String>,
    pub user_id: i64,
}

impl Keyword {
    pub fn new(id: i32, searchterm: Option<String>, user_id: i64) -> Self {
        Self {
            id,
            searchterm,
            user_id,
        }
    }
}

#[derive(Debug, Clone)]
pub struct NewKeyword {
    pub searchterm: Option<String>,
    pub user_id: i64,
}

pub async fn handle_keyword_creation(
    connection: &Pool<Postgres>,
    searchterm: &str,
    user_id: i64,
) -> Result<Keyword, sqlx::Error> {
    let searchterm_already_saved =
        check_if_keyword_already_existed(connection, searchterm, user_id).await;

    if let Some(found_keyword) = searchterm_already_saved {
        Ok(Keyword::new(
            found_keyword.id,
            found_keyword.searchterm.clone(),
            found_keyword.user_id,
        ))
    } else {
        let searchterm = searchterm;
        let id = create_keyword(connection, searchterm, user_id).await?;

        log::info!("Inserting keyword {:?} was successful", searchterm);
        Ok(Keyword::new(id, Some(searchterm.to_string()), user_id))
    }
}

async fn check_if_keyword_already_existed(
    connection: &Pool<Postgres>,
    searchterm: &str,
    user_id: i64,
) -> Option<Keyword> {
    let keywords = find_keywords_per_user(connection, user_id).await;

    keywords
        .unwrap()
        .into_iter()
        .find(|keyword| keyword.searchterm.clone().unwrap() == searchterm)
}
