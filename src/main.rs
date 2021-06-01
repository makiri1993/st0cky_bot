use std::env;

use sqlx::{PgPool, Pool, Postgres};

use crate::bot::run;

mod bot;
mod console;
mod database;
mod models;

#[tokio::main]
async fn main() {
    log::set_logger(&console::CONSOLE_LOGGER).unwrap();
    log::set_max_level(log::LevelFilter::Debug);

    log::info!("Initialized logging...");

    let connection = init_db().await;
    run(connection).await;
}

async fn init_db() -> Pool<Postgres> {
    log::info!("Initialize db...");

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    PgPool::connect(&database_url)
        .await
        .unwrap_or_else(|_| panic!("Error connecting to {}", database_url))
    // Arc::new(Mutex::new(
    //     PgConnection::establish(&database_url)
    //         .expect(&format!("Error connecting to {}", database_url)),
    // ))
}
