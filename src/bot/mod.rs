use crate::models::api::BingNewsApiResponse;
use crate::models::keyword::handle_keyword_creation;
use crate::models::user::handle_user_creation;
use sqlx::{Pool, Postgres};
use std::env;
use std::error::Error;
use teloxide::payloads::SendMessageSetters;
use teloxide::{prelude::*, utils::command::BotCommand};
use tokio_stream::wrappers::UnboundedReceiverStream;

pub async fn run(connection: Pool<Postgres>) {
    log::info!("Starting simple_commands_bot...");

    let bot = Bot::from_env().auto_send();
    let bot_name: String = "st0cky".to_string();
    let client = reqwest::Client::new();

    type DispatcherHandler<T> = DispatcherHandlerRx<AutoSend<Bot>, T>;

    Dispatcher::new(bot)
        .messages_handler(move |rx: DispatcherHandler<Message>| {
            UnboundedReceiverStream::new(rx)
                .commands(bot_name)
                .for_each_concurrent(None, move |(context, command)| {
                    let connection_clone = connection.clone();
                    let client_clone = client.clone();
                    async move {
                        answer(context, command, connection_clone, client_clone)
                            .await
                            .unwrap();
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
    #[command(description = "get news for a specific searchterm")]
    GetNews(String),
}

fn get_telegram_details(cx: &UpdateWithCx<AutoSend<Bot>, Message>) -> (i64, Option<String>) {
    let x = cx.update.from().expect("No user data available");

    (x.to_owned().id, x.to_owned().username)
}

const BING_URL: &str = "https://api.cognitive.microsoft.com/bing/v7.0/news/search";
const BING_HEADER_KEY: &str = "Ocp-Apim-Subscription-Key";

async fn answer(
    cx: UpdateWithCx<AutoSend<Bot>, Message>,
    command: Command,
    connection: Pool<Postgres>,
    client: reqwest::Client,
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
        Command::GetNews(searchterm) => {
            let response = client
                .get(BING_URL)
                .header(
                    BING_HEADER_KEY,
                    env::var("BING_API_TOKEN").expect("BING_API_TOKEN not set"),
                )
                .query(&[("q", &searchterm)])
                .send()
                .await?
                .json::<BingNewsApiResponse>()
                .await?;

            let pretty_printed_news = response.value.iter().fold("".to_string(), |curr, next| {
                format!(
                    "{}<b><u>🔔 {}</u></b>\n\n{}\n\n🤓 <a href=\"{}\">Link</a>\n\n\n\n",
                    curr, next.name, next.description, next.url
                )
            });
            cx.answer(pretty_printed_news)
                .parse_mode(teloxide::types::ParseMode::Html)
                .await?
        }
    };

    Ok(())
}
