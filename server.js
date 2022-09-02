import { createServer } from "http";
import { Server } from "socket.io";

const httpServer = createServer();
const io = new Server(httpServer, {
  // options
});

io.on("connect_error", (err) => {
    console.log(`connect_error due to ${err.message}`);
  });



io.on('connection', (socket) => {
    socket.on('disconnect', (socket) => {
      console.log('a user disconnected');
    });
  
    socket.on('message', (message) => {
      console.log('a user says: ' + message);
    });
  
    console.log('a user connected');
  });
  
  httpServer.listen(8001);