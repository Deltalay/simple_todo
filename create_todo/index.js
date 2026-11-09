import dotenv from "dotenv";

dotenv.config();
import { loadSync } from "@grpc/proto-loader";
import {
    loadPackageDefinition,
    Server,
    ServerCredentials,
} from "@grpc/grpc-js";
import { Pool } from "pg";
const connectionString = process.env.DATABASE_URL;
if (!connectionString) {
    throw new Error("Please set DATABASE_URL in .env");
}
const pool = new Pool({
    max: 10,
    connectionString,
});
pool.on("connect", () => {
    console.log("Connected to database");
});
const pakcageDef = loadSync("../shared_proto/todo.proto");
const packageObject = loadPackageDefinition(pakcageDef);
const server = new Server();
function createTodo(request) {
    let title = request.title;
    let description = request.description;
    let is_done = request.is_done;
    if (!title) title = "";
    if (!description) description = "";
    if (!is_done) is_done = false;
    console.log(title, description, is_done);
    return {
        todo: {
            id: "123",
            title: request.title,
            description: request.description || "",
            is_done: request.is_done,
        },
    };
}
server.addService(packageObject.todo.v1.CreateTodoService.service, {
    createTodo: (call, callback) => {
        callback(null, createTodo(call.request));
    },
});
server.bindAsync(
    "0.0.0.0:50051",
    ServerCredentials.createInsecure(),
    (error, port) => {
        if (error) {
            console.error("Failed to bind:", error);
            return;
        }

        console.log(`gRPC server listening on port ${port}`);
    },
);
