use crate::models::keyword::handle_keyword_creation;
use crate::models::user::handle_user_creation;
use sqlx::{Pool, Postgres};
use std::error::Error;
use teloxide::{prelude::*, utils::command::BotCommand};
use tokio_stream::wrappers::UnboundedReceiverStream;

pub async fn run(connection: Pool<Postgres>) {
    log::info!("Starting simple_commands_bot...");

    let bot = Bot::from_env().auto_send();
    let bot_name: String = "st0cky".to_string();

    type DispatcherHandler<T> = DispatcherHandlerRx<AutoSend<Bot>, T>;

    Dispatcher::new(bot)
        .messages_handler(move |rx: DispatcherHandler<Message>| {
            UnboundedReceiverStream::new(rx)
                .commands(bot_name)
                .for_each_concurrent(None, move |(context, command)| {
                    let connection_clone = connection.clone();

                    async move {
                        answer(context, command, connection_clone).await.unwrap();
                    }
                })
        })
        .dispatch()
        .await;
}

#[derive(BotCommand, PartialEq, Debug)]
#[command(rename = "lowercase", description = "These commands are supported:")]
enum Command {
    #[command(description = "display this text.")]
    Help,
    #[command(description = "add a search keyword")]
    AddKeyword(String),
}

fn get_telegram_details(cx: &UpdateWithCx<AutoSend<Bot>, Message>) -> (i64, Option<String>) {
    let x = cx.update.from().expect("No user data available");

    (x.to_owned().id, x.to_owned().username)
}

async fn answer(
    cx: UpdateWithCx<AutoSend<Bot>, Message>,
    command: Command,
    connection: Pool<Postgres>,
) -> Result<(), Box<dyn Error + Send + Sync>> {
    let telegram_details = get_telegram_details(&cx);
    log::info!("Received a command from user {:?}", telegram_details.0);
    let new_user = handle_user_creation(&connection, telegram_details).await;
    match command {
        Command::Help => cx.answer(Command::descriptions()).send().await?,
        Command::AddKeyword(searchterm) => {
            handle_keyword_creation(&connection, &searchterm, new_user.id).await?;

            cx.answer(format!(
                "I added \"{}\" to your search keywords.",
                searchterm
            ))
            .await?
        }
    };

    Ok(())
}
