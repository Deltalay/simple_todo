use actix_cors::Cors;
use actix_web::{App, Error, HttpResponse, HttpServer, Responder, get, post, web};

use auth::auth_service_client::AuthServiceClient;
use todo::create_todo_service_client::CreateTodoServiceClient;
mod auth {
    tonic::include_proto!("auth.v1");
}
mod todo {
    tonic::include_proto!("todo.v1");
}
mod models;
#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}

#[post("/todo")]
async fn create_todo(body: web::Json<models::CreateTodoRequest>) -> Result<HttpResponse, Error> {
    let request = body.into_inner();

    let grpc_request = todo::CreateTodoRequest {
        title: request.title,
        description: request.description.unwrap_or_default(),
        is_done: request.is_done,
        id: None,
    };
    let mut client = CreateTodoServiceClient::connect("http://127.0.0.1:50051")
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;
    let request = tonic::Request::new(grpc_request);
    let response_result = client
        .create_todo(request)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    let response = response_result.into_inner();
    let todo = response
        .todo
        .ok_or_else(|| actix_web::error::ErrorInternalServerError("Todo not found"))?;
    let serial_res = models::CreateTodoResponse {
        id: todo.id,
        title: todo.title,
        description: todo.description,
        is_done: todo.is_done,
    };
    Ok(HttpResponse::Ok().json(serial_res))
}
#[post("/auth/register")]
async fn register_auth(body: web::Json<models::RegisterRequest>) -> Result<HttpResponse, Error> {
    let request = body.into_inner();
    let grpc_request = auth::RegisterRequest {
        email: request.email,
        password: request.password,
        username: request.username,
    };
    let mut client = AuthServiceClient::connect("http://127.0.0.1:50052")
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;
    let request = tonic::Request::new(grpc_request);
    let response_result = client
        .register(request)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;
    let response = response_result.into_inner();
    let res = models::RegisterResponse {
        status: response.status,
        message: response.message,
    };
    Ok(HttpResponse::Ok().json(res))
}
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        let cors = Cors::default()
            .allow_any_origin()
            .send_wildcard()
            .allow_any_method()
            .allow_any_header();
        App::new()
            .service(hello)
            .service(create_todo)
            .service(register_auth)
            .wrap(cors)
    })
    .workers(4)
    .bind(("127.0.0.1", 8000))?
    .run()
    .await
}
