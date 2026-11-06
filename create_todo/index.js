const protoLoader = require('@grpc/proto-loader');
const grpcLibrary = require('@grpc/grpc-js');
const pakcageDef = protoLoader.loadSync("../shared/proto/todo.proto")
const packageObject = grpcLibrary.loadPackageDefinition(pakcageDef);
const server = new grpcLibrary.Server();
function createTodo(request)
{
    
}
server.addService(packageObject.CreateTodoService.service, {
    createTodo: (call, callback) => {
        callback(null, createTodo(call.request))
    }
})