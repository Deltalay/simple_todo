require("dotenv").config();
import { loadSync } from "@grpc/proto-loader";
import { loadPackageDefinition, Server } from "@grpc/grpc-js";
import { Pool } from "pg";
const connectionString = process.env.DATABASE_URL;
if (!connectionString) {
    throw new Error("Please set DATABASE_URL in .env");
}
const pool = await Pool({
    onConnect: () => {
        console.log("Connected to database");
    },
    max: 10,
    connectionString,
});

const pakcageDef = loadSync("../shared/proto/todo.proto");
const packageObject = loadPackageDefinition(pakcageDef);
const server = new Server();
function createTodo(request) {}
server.addService(packageObject.CreateTodoService.service, {
    createTodo: (call, callback) => {
        callback(null, createTodo(call.request));
    },
});
