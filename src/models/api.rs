use chrono::{DateTime, Local, NaiveDateTime, Utc};
use serde::{de, Deserialize, Deserializer};
use sqlx::types::BigDecimal;
use std::fmt::Display;
use std::str::FromStr;

#[derive(Deserialize, Debug, Clone)]
pub struct ContextualNewsApiResponse {
    pub value: Vec<ContextualNews>,
}

#[derive(Deserialize, Debug, Clone)]
pub struct ContextualNews {
    #[serde(deserialize_with = "from_str")]
    pub id: BigDecimal,
    pub title: String,
    pub url: String,
    pub description: String,
    pub body: String,
    pub snippet: String,
    // provider: [,
    //   {
    //     _type: Organization,,
    //     name: Transfermarkt,
    //   }
    // ],
    #[serde(rename(deserialize = "datePublished"))]
    pub date_published: NaiveDateTime,
}

fn from_str<'de, T, D>(deserializer: D) -> Result<T, D::Error>
where
    T: FromStr,
    T::Err: Display,
    D: Deserializer<'de>,
{
    let s = String::deserialize(deserializer)?;
    T::from_str(&s).map_err(de::Error::custom)
}
