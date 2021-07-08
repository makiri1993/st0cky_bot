#[derive(Clone)]
pub struct News {
    pub id: i32,
    pub searchterm: Option<String>,
    pub user_id: i64,
}

impl News {
    pub fn new(id: i32, searchterm: Option<String>, user_id: i64) -> Self {
        Self {
            id,
            searchterm,
            user_id,
        }
    }
}
