use actix_cors::Cors;
use actix_web::{App, Error, HttpResponse, HttpServer, Responder, get, post, web};
use serde::{Deserialize, Serialize};
use todo::create_todo_service_client::CreateTodoServiceClient;
mod auth {
    tonic::include_proto!("auth.v1");
}
mod todo {
    tonic::include_proto!("todo.v1");
}
#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}
#[derive(Debug, Serialize, Deserialize)]
struct CreateTodoResponse {
    id: String,
    title: String,
    description: String,
    is_done: bool,
}
#[derive(Debug, Serialize, Deserialize)]
struct CreateTodoRequest {
    title: String,
    description: Option<String>,
    is_done: bool,
}
#[post("/todo")]
async fn create_todo(body: web::Json<CreateTodoRequest>) -> Result<HttpResponse, Error> {
    let request = body.into_inner();

    let grpc_request = CreateTodoRequest {
        title: request.title,
        description: request.description,
        is_done: request.is_done,
    };
    // Temp
    Ok(HttpResponse::Ok().json(obj))
}
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        let cors = Cors::default()
            .allow_any_origin()
            .send_wildcard()
            .allow_any_method()
            .allow_any_header();
        App::new().service(hello).wrap(cors)
    })
    .workers(4)
    .bind(("127.0.0.1", 8000))?
    .run()
    .await
}
