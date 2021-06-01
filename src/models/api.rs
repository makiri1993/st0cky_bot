use chrono::{DateTime, Utc};
use serde::Deserialize;

#[derive(Deserialize, Debug, Clone)]
pub struct BingNewsApiResponse {
    pub value: Vec<BingNews>,
}

#[derive(Deserialize, Debug, Clone)]
pub struct BingNews {
    pub name: String,
    pub url: String,
    pub description: String,
    // provider: [,
    //   {
    //     _type: Organization,,
    //     name: Transfermarkt,
    //   }
    // ],
    #[serde(rename(deserialize = "datePublished"))]
    date_published: DateTime<Utc>,
    category: Option<String>,
}
