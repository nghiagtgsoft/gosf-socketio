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
  setTimeout(() => {
    socket.send("\"HELLO GOLANG!!\"")
  }, 1000);
    socket.on('disconnect', (socket) => {
      console.log('a user disconnected');
    });
  
    socket.on('message', (message) => {
      console.log('a user says: ' + message);
      socket.emit('message', "hello YOU!")
    });
  
    console.log('a user connected');
  });
  
  httpServer.listen(8001);