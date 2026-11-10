use serde::{Deserialize, Serialize};
#[derive(Debug, Serialize, Deserialize)]
pub struct CreateTodoResponse {
    pub id: String,
    pub title: String,
    pub description: String,
    pub is_done: bool,
}
#[derive(Debug, Serialize, Deserialize)]
pub struct CreateTodoRequest {
    pub title: String,
    pub description: Option<String>,
    pub is_done: bool,
}
