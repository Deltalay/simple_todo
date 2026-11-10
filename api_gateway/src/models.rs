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
#[derive(Debug, Serialize, Deserialize)]
pub struct RegisterRequest {
    pub username: String,
    pub email: String,
    pub password: String,
}
#[derive(Debug, Serialize, Deserialize)]
pub struct RegisterResponse {
    pub status: i32,
    pub message: Option<String>,
}
