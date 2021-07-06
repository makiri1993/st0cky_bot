use crate::models::api::{ContextualNews, ContextualNewsApiResponse};
use crate::models::keyword::handle_keyword_creation;
use crate::models::user::handle_user_creation;
use sqlx::{Pool, Postgres};
use std::env;
use std::error::Error;
use teloxide::{prelude::*, types::BotCommand as BotStruct, utils::command::BotCommand};
use tokio_stream::wrappers::UnboundedReceiverStream;

#[derive(BotCommand, PartialEq, Debug)]
#[command(rename = "lowercase", description = "These commands are supported:")]
enum TelegramCommand {
    #[command(description = "display this text.")]
    Help,
    #[command(description = "add a search keyword")]
    AddKeyword(String),
    #[command(description = "get news for a specific searchterm")]
    GetNews(String),
}

pub async fn run(connection: Pool<Postgres>, client: reqwest::Client) {
    log::info!("Starting simple_commands_bot...");

    let bot = Bot::from_env().auto_send();
    let bot_name: String = "st0cky".to_string();

    type DispatcherHandler<T> = DispatcherHandlerRx<AutoSend<Bot>, T>;
    set_custom_commands(&bot).await;

    Dispatcher::new(bot)
        .messages_handler(move |rx: DispatcherHandler<Message>| {
            UnboundedReceiverStream::new(rx)
                .commands(bot_name)
                .for_each_concurrent(None, move |(context, command)| {
                    let connection_clone = connection.clone();
                    let client_clone = client.clone();
                    async move {
                        answer(context, command, connection_clone, &client_clone)
                            .await
                            .unwrap();
                    }
                })
        })
        .dispatch()
        .await;
}

async fn set_custom_commands(bot: &AutoSend<Bot>) {
    let set_commands_result: Result<u32, teloxide::RequestError> = bot
        .set_my_commands(vec![
            BotStruct::new("/help", "Get a list of all commands and how to use them."),
            BotStruct::new(
                "/getnews",
                "Add a searchterm after the command and get your news",
            ),
            BotStruct::new(
                "/addkeyword",
                "Add a searchterm so st0cky knows what to look for.",
            ),
        ])
        .await;

    match set_commands_result {
        Ok(_) => {
            log::info!("Custom commands were set in Telegram.");
        }
        Err(_) => {
            log::error!("Custom commands were set in Telegram.");
        }
    };
}

fn get_telegram_details(cx: &UpdateWithCx<AutoSend<Bot>, Message>) -> (i64, Option<String>) {
    let x = cx.update.from().expect("No user data available");

    (x.to_owned().id, x.to_owned().username)
}

const CONTEXTUAL_SEARCH_URL: &str =
    "https://contextualwebsearch-websearch-v1.p.rapidapi.com/api/search/NewsSearchAPI";
const RAPID_KEY_HEADER_KEY: &str = "x-rapidapi-key";

async fn answer(
    cx: UpdateWithCx<AutoSend<Bot>, Message>,
    command: TelegramCommand,
    connection: Pool<Postgres>,
    client: &reqwest::Client,
) -> Result<(), Box<dyn Error + Send + Sync>> {
    let telegram_details = get_telegram_details(&cx);
    log::info!("Received a command from user {:?}", telegram_details.0);
    let new_user = handle_user_creation(&connection, telegram_details).await;
    match command {
        TelegramCommand::Help => cx.answer(TelegramCommand::descriptions()).send().await?,
        TelegramCommand::AddKeyword(searchterm) => {
            if searchterm.is_empty() {
                cx.answer("Sorry you didn't provide any keyword.").await?
            } else {
                handle_keyword_creation(&connection, &searchterm, new_user.id).await?;

                cx.answer(format!(
                    "I added \"{}\" to your search keywords.",
                    searchterm
                ))
                .await?
            }
        }
        TelegramCommand::GetNews(searchterm) => {
            let pretty_printed_news =
                fetch_news_per_keyword_and_pretty_print(client, searchterm).await?;
            cx.answer(pretty_printed_news)
                .parse_mode(teloxide::types::ParseMode::Html)
                .await?
        }
    };

    Ok(())
}

async fn fetch_news_per_keyword_and_pretty_print(
    client: &reqwest::Client,
    searchterm: String,
) -> Result<String, Box<dyn Error + Send + Sync>> {
    let news = fetch_news_per_keyword(client, searchterm).await?;
    let pretty_printed_news = pretty_print_news(news);

    Ok(pretty_printed_news)
}

pub async fn fetch_news_per_keyword(
    client: &reqwest::Client,
    searchterm: String,
) -> Result<Vec<ContextualNews>, reqwest::Error> {
    let request = client
        .get(CONTEXTUAL_SEARCH_URL)
        .header(
            RAPID_KEY_HEADER_KEY,
            env::var("RAPID_KEY_TOKEN").expect("RAPID_KEY_TOKEN not set"),
        )
        .query(&[
            ("q", searchterm.as_str()),
            ("pageNumber", "1"),
            ("pageSize", "10"),
            ("autoCorrect", "true"),
        ]);
    log::info!("{:?}", request);
    let response = request
        .send()
        .await?
        .json::<ContextualNewsApiResponse>()
        .await?;
    log::info!("Fetched {} news", response.value.len());
    Ok(response.value)
}

pub fn pretty_print_news(news: Vec<ContextualNews>) -> String {
    news.iter().fold("".to_string(), |curr, next| {
        format!(
            "{}<b><u>🔔 {}</u></b>\n\n{}\n\n🤓 <a href=\"{}\">Link</a>\n\n\n\n",
            curr, next.title, next.description, next.url
        )
    })
}
