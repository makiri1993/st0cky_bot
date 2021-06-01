use crate::database::create_user;
use sqlx::{Pool, Postgres};

pub struct User {
    pub id: i64,
    pub name: Option<String>,
    pub automatic_sending: Option<bool>,
    // pub created_at: Option<Timestamp>,
    // pub updated_at: Option<Timestamp>,
    // pub deleted_at: Option<Timestamp>,
}

pub async fn handle_user_creation(
    connection: &Pool<Postgres>,
    telegram_details: (i64, Option<String>),
) -> User {
    log::info!("Creating user");
    let new_user = User {
        id: telegram_details.0,
        name: telegram_details.1,
        automatic_sending: Option::from(true),
    };
    let result = create_user(
        connection,
        new_user.id,
        &new_user.name,
        new_user.automatic_sending,
    )
    .await;

    let user_name = new_user.name.clone().unwrap();
    let user_id = new_user.id;
    match result {
        Ok(pg_response) => {
            if pg_response.rows_affected() == 1 {
                log::info!(
                    "Inserting user {} with id {} was succesfull",
                    user_name,
                    user_id
                );
            } else {
                log::info!("User {} with id {} already existed", user_name, user_id);
            }
            new_user
        }
        Err(err) => {
            log::error!("User {} was not created!\n{}", user_name, err);
            new_user
        }
    }
}
