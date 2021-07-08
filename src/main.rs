use std::env;

use sqlx::{PgPool, Pool, Postgres};

use crate::bot::{fetch_news_per_keyword, pretty_print_news, run};
use crate::database::{create_news, find_all_users, find_keywords_per_user};

mod bot;
mod console;
mod database;
mod models;

#[tokio::main]
async fn main() {
    log::set_logger(&console::CONSOLE_LOGGER).unwrap();
    log::set_max_level(log::LevelFilter::Debug);

    log::info!("Initialized logging...");
    let client = reqwest::Client::new();
    let client_clone = client.clone();
    let connection = init_db().await;
    let connection_clone = connection.clone();

    tokio::join!(
        async move {
            loop {
                fetch_news_periodically_and_save(&connection_clone, &client_clone).await;
                tokio::time::sleep(tokio::time::Duration::from_secs(5)).await;
            }
        },
        run(connection, client)
    );
}

async fn init_db() -> Pool<Postgres> {
    log::info!("Initialize db...");

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    PgPool::connect(&database_url)
        .await
        .unwrap_or_else(|_| panic!("Error connecting to {}", database_url))
}

async fn fetch_news_periodically_and_save(connection: &Pool<Postgres>, client: &reqwest::Client) {
    let users = find_all_users(&connection).await;

    for user in users.unwrap() {
        let user_id = user.id;
        let keywords = find_keywords_per_user(&connection, &user_id).await;
        for keyword in keywords.unwrap() {
            let result_news = fetch_news_per_keyword(&client, keyword.searchterm.unwrap()).await;
            match result_news {
                Ok(news) => {
                    for news in news {
                        create_news(&connection, news, user_id).await;
                    }
                    // let string = pretty_print_news(news);
                    // log::info!("{}", string);
                }
                Err(err) => log::error!("{:?}", err),
            }
        }
    }
    println!("I run every 1 seconds");
}
